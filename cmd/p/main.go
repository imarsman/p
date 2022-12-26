package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/imarsman/p/cmd/internal/args"
)

func wait() {
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	text := strings.TrimSpace(input.Text())
	if strings.HasPrefix(text, `!`) {
		var buffStdOut bytes.Buffer
		var buffStdErr bytes.Buffer
		cmd := exec.Command("bash", "-c", strings.TrimSpace(text[1:]))
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
	if strings.TrimSpace(text) != "" {
		fmt.Println(text)
	}
	return
}

func main() {
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
		defer file.Close()

		// loop until all lines processed
		for {
			count := 1
			scanner := bufio.NewScanner(file)

			// optionally, resize scanner's capacity for lines over 64K, see next example
			for scanner.Scan() {
				if count == args.Args.Number {
					wait()
					count = 1
				}
				fmt.Println(scanner.Text())
				count++
			}
			break
		}
	}
}
