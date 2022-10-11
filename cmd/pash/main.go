// The pash package provides a CLI tool for generating password hashes.
// The command reads one password per line on stdin and outputs a hash of that password suitable for use in storage.
// When passwords are entered into a terminal echo is disabled.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/vanti-dev/bsp-ew/internal/util/pass"
	"golang.org/x/term"
	"os"
	"syscall"
)

var (
	check string
)

func init() {
	flag.StringVar(&check, "check", "", "Check a password hash instead of generating them")
}

func main() {
	flag.Parse()

	if len(check) > 0 {
		checkHash([]byte(check))
		return
	}

	if term.IsTerminal(syscall.Stdin) {
		readFromStdin()
	} else {
		readFromScanner(bufio.NewScanner(os.Stdin))
	}
}

func checkHash(hash []byte) {
	input, err := readOnePass()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Err reading password: %v\n", err)
		os.Exit(1)
	}

	err = pass.Compare(hash, input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintf(os.Stderr, "Read hash: '%v'\n", string(hash))
		os.Exit(2)
	}

	fmt.Println("Success")
}

func readOnePass() ([]byte, error) {
	if term.IsTerminal(syscall.Stdin) {
		fmt.Print("Password: ")
		input, err := term.ReadPassword(syscall.Stdin)
		fmt.Println()
		return input, err
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return nil, scanner.Err()
		}
		return scanner.Bytes(), nil
	}
}

func readFromStdin() {
	for {
		fmt.Print("Password: ")
		input, err := term.ReadPassword(syscall.Stdin)
		fmt.Println()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Err reading password: %v\n", err)
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
