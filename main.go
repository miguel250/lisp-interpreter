package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	replPtr := flag.Bool("r", false, "REPL mode")
	flag.Parse()

	scope := newScope(nil)
	b := newBuiltins()

	for k, v := range b.fn {
		scope.set(k, v)
	}

	repl := *replPtr

	if repl {
		for {
			fmt.Print(">> ")
			input(scope, true)
		}
	} else {
		input(scope, false)
	}
}

func input(scope *scope, repl bool) {

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		experssions, err := parse(scanner.Text())
		if err != nil {
			fmt.Println(err)
			break
		}

		failed := false
		for _, e := range experssions {
			e, err = eval(e, scope)
			if err != nil {
				fmt.Println(err)
				failed = true
				break
			}

			if repl && e != nil {
				fmt.Println(e)
			}

		}

		if failed {
			break
		}

		if repl {
			fmt.Print(">> ")
		}
	}
}
