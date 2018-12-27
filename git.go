package main

import (
	"github.com/pkg/errors"
	"regexp"
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

var expressions map[URLType]struct{
	*regexp.Regexp
	parser
}

type parser func(url string) (*GitURL)

func init() {

	username := `(\w+)`
	domain := `((\w+)(\.\w+)*)`
	port := `(:[0-9]+)?`
	path := `((.+/)*(.+)).git(/)?`

	expressions = map[URLType]struct{
		*regexp.Regexp
		parser
	}{
		URLTypeSSH: {
			regexp.MustCompile(`^ssh://(` + username + `@)?` + domain + port + `/` + path + `$`),
			parseSSH,
		},
		URLTypeGit: {
			regexp.MustCompile(`^git://` + domain + port + `/` + path + `(#.+)?$`),
			parseGit,
		},
		URLTypeHTTP: {
			regexp.MustCompile(`^http(s)?://(` + username + `:.+@)?` + domain +  port + `/` + path + `$`),
			parseHTTP,
		},
		/* Unsupported until we have better test data
		URLTypeFTP: {
			regexp.MustCompile(`^ftp(s)?://` + domain + port + `/` + path + `$`),
			parseFTP,
		},
		*/
		URLTypeSCP: {
			regexp.MustCompile(`^(` + username + `@)` + domain + `:` + path + `$`),
			parseSCP,
		},
	}

}

func NewGitURL(url string) (*GitURL, error) {
	for _, expr := range expressions {
		if expr.Regexp.MatchString(url) {
			return expr.parser(url), nil
		}
	}

	return nil, errors.New("unable to determine type of URL")
}

func parseSSH(url string) (*GitURL) {
	r := expressions[URLTypeSSH].Regexp
	matches := r.FindStringSubmatch(url)
	if matches == nil {
		return nil
	}
	if len(matches) < 6{
		return nil
	}
	return &GitURL{
		Username: matches[2],
		Hostname: matches[3],
		Path: matches[7],
		Type: URLTypeSSH,
	}
}

func parseGit(url string) (*GitURL) {
	r := expressions[URLTypeGit].Regexp
	matches := r.FindStringSubmatch(url)
	if matches == nil {
		return nil
	}
	if len(matches) < 6{
		return nil
	}
	return &GitURL{
		Username: matches[2],
		Hostname: matches[3],
		Path: matches[7],
		Type: URLTypeGit,
	}
}

func parseHTTP(url string) (*GitURL) {
	r := expressions[URLTypeHTTP].Regexp
	matches := r.FindStringSubmatch(url)
	if matches == nil {
		return nil
	}
	if len(matches) < 6{
		return nil
	}
	return &GitURL{
		Username: matches[3],
		Hostname: matches[4],
		Path: matches[8],
		Type: URLTypeHTTP,
	}
}

func parseFTP(url string) (*GitURL) {
	return nil
}

func parseSCP(url string) (*GitURL) {
	return nil
}
