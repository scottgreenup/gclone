package main

import (
	"github.com/pkg/errors"
	"regexp"
	"strings"
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

var expressions map[URLType]*regexp.Regexp

func init() {
	expressions = map[URLType]*regexp.Regexp{
		URLTypeSSH:   regexp.MustCompile(`^ssh://(.+@)?([\w\.]+)(:[0-9]+)?(/.+)+.git(/)?$`),
		URLTypeGit:   regexp.MustCompile(`^git://([\w\.]+)(:[0-9]+)?(/.+)+.git(/|#.+)?$`),
		URLTypeHTTP:  regexp.MustCompile(`^http(s)?://([\w]+:.+@)?([\w\.]+)(:[0-9]+)?(/.+)+.git(/)?$`),
		URLTypeFTP:   regexp.MustCompile(`^ftp(s)?://([\w\.]+)(:[0-9]+)?(/.+)+.git(/)?$`),
		URLTypeSCP:   regexp.MustCompile(`^(\w+@)([\w\.]+):(.+)(/.+)*.git(/)?$`),
	}

}

func NewGitURL(url string) (*GitURL, error) {
	for k, expr := range expressions {
		if expr.MatchString(url) {
			return &GitURL{
				Type: k,
			}, nil
		}
	}

	return nil, errors.New("unable to determine type of URL")
}

func parseSSH(url string) (*GitURL) {
	r := expressions[URLTypeSSH]
	if !r.MatchString(url) {
		return &GitURL{}
	}
	matches := r.FindStringSubmatch(url)
	return &GitURL{
		Username: strings.TrimSuffix(matches[1], "@"),
		Hostname: matches[2],
		Path: matches[4],
	}
}
