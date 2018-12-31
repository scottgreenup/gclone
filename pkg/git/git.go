package git

import (
	"regexp"

	"github.com/pkg/errors"
)

type URLType string

const (
	URLTypeSSH  URLType = "ssh"
	URLTypeGit  URLType = "git"
	URLTypeHTTP URLType = "http"
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
	username := `(?P<username>\w+)`
	domain := `(?P<hostname>(\w+)(\.\w+)*)`
	port := `(?P<port>:[0-9]+)?`
	path := `(?P<path>(.+/)*(.+)).git(/)?`

	expressions = map[URLType]*regexp.Regexp{
		URLTypeSSH:  regexp.MustCompile(`^ssh://(` + username + `@)?` + domain + port + `/` + path + `$`),
		URLTypeGit:  regexp.MustCompile(`^git://` + domain + port + `/` + path + `(#.+)?$`),
		URLTypeHTTP: regexp.MustCompile(`^http(s)?://(` + username + `:.+@)?` + domain + port + `/` + path + `$`),
		URLTypeSCP:  regexp.MustCompile(`^(` + username + `@)` + domain + `:` + path + `$`),
	}
}

func NewURL(url string) (*GitURL, error) {
	for t, expr := range expressions {
		if expr.MatchString(url) {
			return parse(url, t), nil
		}
	}

	return nil, errors.New("unable to determine type of URL")
}

func find(url string, r *regexp.Regexp) map[string]string {
	matches := r.FindStringSubmatch(url)
	if matches == nil {
		return nil
	}
	result := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = matches[i]
		}
	}
	return result
}

func parse(url string, t URLType) *GitURL {
	matches := find(url, expressions[t])
	if matches == nil {
		return nil
	}
	return &GitURL{
		Username: matches["username"],
		Hostname: matches["hostname"],
		Path:     matches["path"],
		Port:     matches["port"],
		Type:     t,
	}
}
