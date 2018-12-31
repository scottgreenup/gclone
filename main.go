package main

import (
	"fmt"
	"os"

	"github.com/scottgreenup/gclone/cmd"
)

func usage() {
	fmt.Println("usage: gclone [FLAGS] repository [DIR] [-- [GCLONE FLAGS]]")
	os.Exit(1)
}

func main() {
	cmd.Execute()
}
