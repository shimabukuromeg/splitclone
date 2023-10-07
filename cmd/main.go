package main

import (
	"os"

	"github.com/shimabukuromeg/splitclone/app"
)

func main() {

	cli := &app.CLI{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}

	cli.Run(os.Args[1:])
}
