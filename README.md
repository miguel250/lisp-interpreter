# Lisp Interpreter

#### Build and run
```bash
go build
cat test.lisp | ./lisp-interpreter
```

#### REPL
```bash
./lisp-interpreter -r
```

#### Built-in functions
* `(setq x 4)`: defined symbol
* `(print x)`: Print variable to stdout
* `(list (1 "hello" 1.3))`: Create a list
* `(first (list (1 "hello" 1.3)))`: Return first value of a list

#### Todo
* Seprate parsing, evaluation and built-ins into their own directories.
* Add support for `if` operators.
