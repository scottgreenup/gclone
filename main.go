package main

import (
	"flag"
	"fmt"
)

func main() {
	var repository string
	flag.StringVar(&repository, "repository", "", "The URL to the git repository")
	flag.Parse()

	fmt.Println(repository)
	NewGitURL(repository)
}
