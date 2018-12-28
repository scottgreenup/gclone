package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	var repository string
	flag.StringVar(&repository, "repository", "", "The URL to the git repository")
	flag.Parse()

	fmt.Println(repository)
	gu, err := NewGitURL(repository)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// TODO normalise, spaces to -, etc...
	path := filepath.Join("$HOME", "code", gu.Hostname, gu.Path)
	fmt.Printf("mkdir -p %s && git clone %s %s\n", path, repository, path)
}
