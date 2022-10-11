// The pash package provides a CLI tool for generating password hashes.
// The command reads one password per line on stdin and outputs a hash of that password suitable for use in storage.
// When passwords are entered into a terminal echo is disabled.
package main

import (
	"bufio"
	"fmt"
	"github.com/vanti-dev/bsp-ew/internal/util/pass"
	"golang.org/x/term"
	"os"
	"syscall"
)

func main() {
	if term.IsTerminal(int(syscall.Stdin)) {
		readFromStdin()
	} else {
		readFromScanner(bufio.NewScanner(os.Stdin))
	}
}

func readFromStdin() {
	for {
		fmt.Print("Password: ")
		input, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Err reading password: %v", err)
			os.Exit(1)
		}
		printPassHash(input)
	}
}

func readFromScanner(s *bufio.Scanner) {
	for s.Scan() {
		input := s.Bytes()
		if len(input) == 0 {
			continue // ignore empty lines
		}
		printPassHash(input)
	}
}

func printPassHash(input []byte) {
	hash, err := pass.Hash(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Hash error: %v", err)
		return
	}
	printHash(hash)
}

func printHash(hash []byte) {
	fmt.Printf("%v\n", string(hash))
}
