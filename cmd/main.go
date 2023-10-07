package main

import (
	"os"

	split "github.com/shimabukuromeg/splitclone"
)

func main() {

	cli := &split.CLI{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}

	cli.Run(os.Args[1:])
}
