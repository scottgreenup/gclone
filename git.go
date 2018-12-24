package main

import (
	"regexp"

	"github.com/pkg/errors"
)

type URLType string

const (
	URLTypeSSH  URLType = "ssh"
	URLTypeGit  URLType = "git"
	URLTypeHTTP URLType = "http"
	URLTypeFTP  URLType = "ftp"
	URLTypeSCP  URLType = "scp"
)

type GitURL struct {

	// The hostname segment of the URL
	Hostname string

	// The path of the URL, without the .git in the final segment
	Path string

	// The port used, automatically filled in if missing.
	Port string

	// The username if present, automatically filled in if possible.
	Username string

	// The type of URL
	Type URLType
}

func NewGitURL(url string) (*GitURL, error) {
	expressions := map[URLType]*regexp.Regexp{
		URLTypeSSH:   regexp.MustCompile(`^ssh://(.+@)?(.+)(:[0-9]+)?(/.+)+.git(/)?$`),
		URLTypeGit:   regexp.MustCompile(`^git://(.+)(:[0-9]+)?(/.+)+.git(/|#.+)?$`),
		URLTypeHTTP:  regexp.MustCompile(`^http(s)?://(.+)(:[0-9]+)?(/.+)+.git(/)?$`),
		URLTypeFTP:   regexp.MustCompile(`^ftp(s)?://(.+)(:[0-9]+)?(/.+)+.git(/)?$`),
		URLTypeSCP:   regexp.MustCompile(`^(\w+@)(.+):(.+)(/.+)*.git(/)?$`),
	}

	for k, expr := range expressions {
		if expr.MatchString(url) {
			return &GitURL{
				Type: k,
			}, nil
		}
	}

	return nil, errors.New("unable to determine type of URL")
}
