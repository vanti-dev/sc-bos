package transport

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/cenkalti/backoff/v4"
	"golang.org/x/crypto/ssh"
)

type Ssh struct {
	*Connection
	conf SshConfig

	// control access to client & session
	lock    sync.RWMutex
	client  *ssh.Client
	session *ssh.Session
	stdin   io.WriteCloser
	stdout  io.Reader
}

// NewSsh creates a new Ssh transport with the given config
func NewSsh(conf SshConfig) *Ssh {
	conf.defaults()
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 0 // never give up
	s := &Ssh{
		conf: conf,
	}
	conn := NewConnection(conf.ConnectionConfig, bo, s.dial, readWriteCloserFuncs(s.read, s.write, s.close))
	s.Connection = conn
	return s
}

func (s *Ssh) read(p []byte) (n int, err error) {
	s.lock.RLock()
	stdout := s.stdout
	s.lock.RUnlock()

	n, err = stdout.Read(p)
	if err != nil {
		s.state.update(Disconnected)
		if err == io.EOF {
			err = nil
		}
	}
	return
}

func (s *Ssh) write(p []byte) (n int, err error) {
	s.lock.RLock()
	stdin := s.stdin
	s.lock.RUnlock()
	return stdin.Write(p)
}

func (s *Ssh) close() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.client != nil {
		err := (*s.client).Close()
		s.client = nil
		return err
	}
	return nil
}

func (s *Ssh) dial() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.client != nil {
		_ = s.session.Close()
		_ = (*s.client).Close()

		s.session = nil
		s.stdin = nil
		s.stdout = nil
		s.client = nil
	}

	var hostKeyCallback func(hostname string, remote net.Addr, key ssh.PublicKey) error
	if s.conf.IgnoreHostKey {
		hostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", s.conf.Ip, s.conf.Port), &ssh.ClientConfig{
		User:            s.conf.Username,
		HostKeyCallback: hostKeyCallback,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.conf.Password),
			ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
				// Just send the password back for all questions
				answers = make([]string, len(questions))
				for i := range answers {
					answers[i] = s.conf.Password
				}

				return answers, nil
			}),
		},
		Timeout: s.conf.Timeout.Duration,
	})
	if err != nil {
		return err
	}

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	err = session.RequestPty("xterm", 40, 80, modes)
	if err != nil {
		return err
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	err = session.Shell()
	if err != nil {
		return err
	}

	s.client = client
	s.session = session
	s.stdin = stdin
	s.stdout = stdout

	return nil
}
