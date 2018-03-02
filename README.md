# Lisp Interpreter [![Build Status](https://travis-ci.org/miguel250/lisp-interpreter.svg?branch=master)](https://travis-ci.org/miguel250/lisp-interpreter) [![Coverage Status](https://coveralls.io/repos/github/miguel250/lisp-interpreter/badge.svg?branch=master)](https://coveralls.io/github/miguel250/lisp-interpreter?branch=master)

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
