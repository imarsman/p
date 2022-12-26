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
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)

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
	if strings.TrimSpace(text) != "" {
		fmt.Println(text)
	}
	return
}

// var input *os.File

func main() {
	var fileList []*os.File
	var stdin string

	// stat, _ := os.Stdin.Stat()
	// if (stat.Mode() & os.ModeCharDevice) == 0 {
	// 	// scanner := bufio.NewScanner(os.Stdin)

	// 	// for scanner.Scan() {
	// 	// 	sb.WriteString(fmt.Sprintf("%s\n", scanner.Text()))
	// 	// }
	// 	contents, err := io.ReadAll(os.Stdin)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	stdin = string(contents)
	// 	os.Stdin.Close()

	// 	input, _, err = os.Pipe()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	os.Stdin = input
	// }

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

	if len(stdin) > 0 {
		scanner := bufio.NewScanner(strings.NewReader(stdin))
		count := 1
		for scanner.Scan() {
			if count == args.Args.Number {
				fmt.Println("waiting")
				wait()
				count = 1
			}
			fmt.Println(scanner.Text(), count)
			count++
		}
	}

	for _, file := range fileList {
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
