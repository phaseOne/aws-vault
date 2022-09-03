package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func TerminalPrompt(message string) (string, error) {
	return withTerminal(func(in, out *os.File) (string, error) {
		fmt.Fprint(out, message)

		reader := bufio.NewReader(in)
		text, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		return strings.TrimSpace(text), nil
	})
}

func TerminalSecretPrompt(message string) (string, error) {
	return withTerminal(func(in, out *os.File) (string, error) {
		fmt.Fprint(out, message)

		text, err := term.ReadPassword(int(in.Fd()))
		if err != nil {
			return "", err
		}

		fmt.Println()

		return strings.TrimSpace(string(text)), nil
	}
}

func TerminalMfaPrompt(mfaSerial string) (string, error) {
	return TerminalPrompt(mfaPromptMessage(mfaSerial))
}

func init() {
	Methods["terminal"] = TerminalMfaPrompt
}

// withTerminal runs f with the terminal input and output files, if available.
// withTerminal does not open a non-terminal stdin, so the caller does not need
// to check stdinInUse.
// Copyright 2021 The age Authors
// This method is governed by a BSD-style license found at
// https://github.com/FiloSottile/age/blob/main/LICENSE
func withTerminal(f func(in, out *os.File) (string, error)) (string, error) {
	if tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0); err == nil {
		defer tty.Close()
		return f(tty, tty)
	} else
		return f(os.Stdin, os.Stderr)
	}
}
