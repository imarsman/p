package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"unicode"

	"github.com/imarsman/p/cmd/internal/args"
)

func centerString(str string, pad string, width int) string {
	padWidth := int(float64(width-len(str)) / 2)
	if str != "" {
		return strings.Repeat(pad, padWidth) + " " + str + " " + strings.Repeat(pad, width-(padWidth+len(str))-2)
	}
	return strings.Repeat(pad, padWidth) + strings.Repeat(pad, width-(padWidth+len(str)))
}

func wait() {
	tty, err := os.Open("/dev/tty")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var r = bufio.NewReader(tty)
	text, err := r.ReadString('\n')
	if len(text) == 0 && err != nil {
		if err == io.EOF {
			return
		}
	}
	if err != nil {
		fmt.Println(err)
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	if text == `q` {
		os.Exit(0)
	}

	if text == `!` {
		return
	}

	if strings.HasPrefix(text, `!`) {
		var buffStdOut bytes.Buffer
		var buffStdErr bytes.Buffer

		text = strings.TrimSpace(text[1:])
		cmd := exec.Command("bash", "-c", text)
		cmd.Stdout = &buffStdOut
		cmd.Stderr = &buffStdErr
		err := cmd.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			if args.Args.Pretty {
				fmt.Println(centerString("command", `-`, 80))
			}
			fmt.Println(strings.TrimSpace(buffStdOut.String()))
			if args.Args.Pretty {
				fmt.Println(centerString("", `-`, 80))
			}
		}
		return
	}
	fmt.Println(text)

	return
}

func output(scanner *bufio.Scanner) {
	count := 1
	total := 1
	for scanner.Scan() {
		if count == args.Args.Lines+1 {
			wait()
			count = 1
		}
		if args.Args.Number {
			fmt.Printf("%8d %s\n", total, supressIf(scanner.Text()))
		} else {
			fmt.Println(supressIf(scanner.Text()))
		}
		count++
		total++
	}
	if scanner.Err() != nil {
		fmt.Printf("%v\n", scanner.Err())
		os.Exit(1)
	}
}

func supressIf(input string) string {
	if args.Args.Supress == false {
		return input
	}
	var sb strings.Builder
	for _, r := range input {
		if unicode.IsPunct(r) {
			sb.WriteRune(r)
			continue
		} else if unicode.IsPrint(r) {
			sb.WriteRune(r)
			continue
		}
		// Don't add if not caught above
		// sb.WriteRune(-1)
	}

	return sb.String()
}

func main() {
	var fileList []*os.File
	var stdin string

	// Currently I can't find a way to read stdin and then read user input in
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		var err error
		buf := new(strings.Builder)
		_, err = io.Copy(buf, os.Stdin)
		if err != nil {
			fmt.Println(err)
		}
		stdin = string(buf.String())
	}

	for _, path := range args.Args.Files {
		if _, err := os.Stat(path); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			} else {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}
		}
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		fileList = append(fileList, file)
		defer file.Close()
	}

	// write stdout
	if len(stdin) > 0 {
		scanner := bufio.NewScanner(strings.NewReader(stdin))
		output(scanner)
	}

	// output files
	for _, file := range fileList {
		scanner := bufio.NewScanner(file)
		output(scanner)
	}
}
