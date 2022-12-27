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

	"github.com/imarsman/p/cmd/internal/args"
)

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
			fmt.Println(buffStdOut.String())
		}
		return
	}
	fmt.Println(text)

	return
}

func output(scanner *bufio.Scanner) {
	count := 1
	for scanner.Scan() {
		if count == args.Args.Number+1 {
			wait()
			count = 1
		}
		fmt.Println(scanner.Text())
		count++
	}
	if scanner.Err() != nil {
		fmt.Printf("%v\n", scanner.Err())
		os.Exit(1)
	}
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
