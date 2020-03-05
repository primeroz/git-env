package main

import (
	"bufio"
	"fmt"
	"os"
)

func cmdInit() {
	values := map[string]string{}
	reader := bufio.NewReader(os.Stdin)

	for _, opt := range options {
		fmt.Printf("%s [%s] ", opt.Question, opt.Default)
		value, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		if value == "\n" {
			values[opt.Name] = opt.Default
		} else {
			values[opt.Name] = value[:len(value)-1]
		}
	}

	for k, v := range values {
		_, err := setGitOption(k, v)
		if err != nil {
			panic(err)
		}
	}

	// todo
	// set rerere enable in --local
	// set conflictstyle=diff3 in global
	// do smth with pre-commit ?

	fmt.Println("You're ready to go.")
}
