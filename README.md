# KeiranScript

KeiranScript is a simple programming language that can be compiled to assembly.
Although we plan to support direct compilation to cross-platform binaries, we currently only support asm output.

## Installation

Just clone the repo and build the project using [Go](https://golang.org). We will include pre-built binaries soon.

## Usage

```bash
$ ks example.ks
```
This will produce a file ending in .asm in the current directory.
This can be compiled to machine code manually, and then run on your system.

There is currently no documentation for the language, but reading the `example.ks` file included in the `examples` directory should give you a good idea. You can also read the source code if you wish.