package transport

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/textproto"
	"os"
	"testing"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func TestConnection(t *testing.T) {
	tcp := NewTcp(TcpConfig{
		Ip:   "localhost",
		Port: 6860,
	}, zap.L())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	reader := bufio.NewReader(os.Stdin)
	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error {
		return tcp.Connect(ctx)
	})
	grp.Go(func() error {
		reader := bufio.NewReader(tcp)
		tp := textproto.NewReader(reader)
		// responses are delimited by \r\n\r\n, and the response itself may well be multiple lines separated by \r\n.
		for {
			log.Printf("read line")
			line, err := tp.ReadLine()
			if err != nil {
				log.Printf("got error: %s", err)
			}
			log.Printf("got line: %s", line)
		}
	})
	grp.Go(func() error {
		for {
			text, _ := reader.ReadString('\n')
			log.Printf("write: %s", text)
			n, err := tcp.Write([]byte(text))
			if err != nil {
				log.Printf("write error: %s", err)
			} else {
				log.Printf("written %d", n)
			}
		}
	})
	grp.Go(func() error {
		prev := Idle
		for {
			state, change := tcp.WaitForStateChange(ctx, prev)
			if !change {
				return fmt.Errorf("no change")
			}
			prev = state
			if state == Connected {
				log.Printf("connected")
			} else if state == Disconnected {
				log.Printf("disconnected")
			}
		}
	})
	err := grp.Wait()
	if err != nil {
		t.Error(err)
	}
}
