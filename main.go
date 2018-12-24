package main

import (
	"flag"
	"fmt"
	"regexp"
)


type GitURL struct {

	// The hostname segment of the URL
	hostname string

	// The path of the URL, without the .git in the final segment
	path string

	// The port used, automatically filled in if missing.
	port string

	// The username if present, automatically filled in if possible.
	username string
}

func NewGitURL(url string) *GitURL {
	expressions := map[string]*regexp.Regexp{
		"ssh":   regexp.MustCompile(`^ssh://([A-Za-z0-9]+@)?([A-Za-z0-9]+)(\.[A-Za-z0-9]+)+(:[0-9]+)?(/[A-Za-z0-9+-\.]+)*/[A-Za-z0-9]+.git(/)?$`),
		"git":   regexp.MustCompile(`^git://([A-Za-z0-9]+)(\.[A-Za-z0-9]+)+(:[0-9]+)?(/[A-Za-z0-9+-\.]+)*/[A-Za-z0-9]+.git(/)?$`),
		"https": regexp.MustCompile(`^http(s)?://([A-Za-z0-9]+)(\.[A-Za-z0-9]+)+(:[0-9]+)?(/[A-Za-z0-9+-\.]+)*/[A-Za-z0-9]+.git(/)?$`),
		"ftps":  regexp.MustCompile(`^ftp(s)?://([A-Za-z0-9]+)(\.[A-Za-z0-9]+)+(:[0-9]+)?(/[A-Za-z0-9+-\.]+)*/[A-Za-z0-9]+.git(/)?$`),
		"scp":   regexp.MustCompile(`^([A-Za-z0-9]+@)?([A-Za-z0-9]+)(\.[A-Za-z0-9]+)+:([A-Za-z0-9+-\.]+(/[A-Za-z0-9+-\.]+)*/)?[A-Za-z0-9]+.git(/)?$`),
	}

	for name, expr := range expressions {
		s := expr.FindString(url)
		fmt.Println(name, s)
	}

	return nil
}

func main() {
	var repository string
	flag.StringVar(&repository, "repository", "", "The URL to the git repository")
	flag.Parse()

	fmt.Println(repository)
	NewGitURL(repository)
}
