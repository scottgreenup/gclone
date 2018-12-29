package git

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertURLType(t *testing.T, expected URLType, actual *GitURL, url string) {
	require.NotNil(t, actual, fmt.Sprintf("Expected %s but received %+v from %q", expected, actual, url))
	assert.True(t, expected == actual.Type, fmt.Sprintf("Expected %s but received %s from %q", expected, actual.Type, url))
}

func TestGitURL(t *testing.T) {
	t.Run("git", func(t *testing.T) {

		testcases := []string{
			"git://github.com/scottgreenup/gclone.git",
			"git://host.xz/path/to/repo.git/",
			"git://github.com/ember-cli/ember-cli.git#ff786f9f",
			"git://github.com/ember-cli/ember-cli.git#gh-pages",
			"git://github.com/ember-cli/ember-cli.git#master",
			"git://github.com/ember-cli/ember-cli.git#Quick-Fix",
			"git://github.com/ember-cli/ember-cli.git#quick_fix",
			"git://github.com/ember-cli/ember-cli.git#v0.1.0",
		}

		for _, tc := range testcases {
			gu, err := NewURL(tc)
			require.NoError(t, err, tc)
			assertURLType(t, URLTypeGit, gu, tc)
		}
	})

	t.Run("http", func(t *testing.T) {

		testcases := []string{
			"https://github.com/kubernetes/kubernetes.git",
			"http://192.168.101.127/user/project.git",
			"http://github.com/user/project.git",
			"http://host.xz/path/to/repo.git/",
			"https://192.168.101.127/user/project.git",
			"https://github.com/user/project.git",
			"https://host.xz/path/to/repo.git/",
			"https://username::;*%$:@github.com/username/repository.git",
			"https://username:$fooABC@:@github.com/username/repository.git",
			"https://username:password@github.com/username/repository.git",
		}

		for _, tc := range testcases {
			gu, err := NewURL(tc)
			require.NoError(t, err, tc)
			assertURLType(t, URLTypeHTTP, gu, tc)
		}
	})

	t.Run("scp", func(t *testing.T) {

		testcases := []string{
			"git@github.com:kubernetes/kubernetes.git",
			"git@gitlab.com:facebook/react.git",
			"git@192.168.101.127:user/project.git",
			"git@github.com:user/project.git",
			"git@github.com:user/some-project.git",
			"git@github.com:user/some_project.git",
		}

		for _, tc := range testcases {
			gu, err := NewURL(tc)
			require.NoError(t, err)
			assertURLType(t, URLTypeSCP, gu, tc)
		}
	})

	t.Run("ssh", func(t *testing.T) {

		testcases := []string{
			"ssh://host.xz/path/to/repo.git/",
			"ssh://host.xz/path/to/repo.git/",
			"ssh://host.xz/~/path/to/repo.git",
			"ssh://host.xz/~user/path/to/repo.git/",
			"ssh://host.xz:7777/path/to/repo.git/",
			"ssh://user@host.xz/path/to/repo.git/",
			"ssh://user@host.xz/path/to/repo.git/",
			"ssh://user@host.xz/~/path/to/repo.git",
			"ssh://user@host.xz/~user/path/to/repo.git/",
			"ssh://user@host.xz:7777/path/to/repo.git/",
		}

		for _, tc := range testcases {
			gu, err := NewURL(tc)
			require.NoError(t, err, tc)
			assertURLType(t, URLTypeSSH, gu, tc)
		}
	})
}

func TestLegalParsing(t *testing.T) {

	t.Run("ssh", func(t *testing.T) {

		testcases := []struct{
			input string
			output GitURL
		}{
			{ "ssh://host.xz/path/to/repo.git/", GitURL{
				Hostname: "host.xz",
				Path: "path/to/repo",
			}},

			// TODO what does the tilde mean?
			{ "ssh://host.xz/~/path/to/repo.git", GitURL{
				Hostname: "host.xz",
				Path: "~/path/to/repo",
			}},

			// TODO what does the tilde mean?
			{ "ssh://host.xz/~user/path/to/repo.git/", GitURL{
				Hostname: "host.xz",
				Path: "~user/path/to/repo",
			}},
			{ "ssh://host.xz:1/path/to/repo.git/", GitURL{
				Hostname: "host.xz",
				Path: "path/to/repo",
			}},
			{ "ssh://user@host.xz/path/to/repo.git/", GitURL{
				Hostname: "host.xz",
				Path:     "path/to/repo",
				Username: "user",
			}},
			{ "ssh://user@host.xz/path/to/repo.git/", GitURL{
				Hostname: "host.xz",
				Path:     "path/to/repo",
				Username: "user",
			}},
			{ "ssh://user@host.xz/~/path/to/repo.git", GitURL{
				Hostname: "host.xz",
				Path:     "~/path/to/repo",
				Username: "user",
			}},
			{ "ssh://user@host.xz/~user/path/to/repo.git/", GitURL{
				Hostname: "host.xz",
				Path:     "~user/path/to/repo",
				Username: "user",
			}},
			{ "ssh://user@host.xz:1/path/to/repo.git/", GitURL{
				Hostname: "host.xz",
				Path:     "path/to/repo",
				Username: "user",
			}},
			{ "ssh://user@some.crazy.domain.net.au:1/path/to/repo.git/", GitURL{
				Hostname: "some.crazy.domain.net.au",
				Path:     "path/to/repo",
				Username: "user",
			}},
			{ "ssh://user@some.crazy.domain.net.au:1/a.git", GitURL{
				Hostname: "some.crazy.domain.net.au",
				Path:     "a",
				Username: "user",
			}},
		}

		for _, tc := range testcases {
			gu := parse(tc.input, URLTypeSSH)
			require.NotNil(t, gu, tc.input)

			assert.Equal(t, tc.output.Username, gu.Username)
			assert.Equal(t, tc.output.Hostname, gu.Hostname)
			assert.Equal(t, tc.output.Path, gu.Path)
			assert.Equal(t, URLTypeSSH, gu.Type)
		}
	})

	t.Run("http", func(t *testing.T) {

		testcases := []struct{
			input string
			output GitURL
		}{
			{"https://github.com/kubernetes/kubernetes.git", GitURL{
				Hostname: "github.com",
				Path: "kubernetes/kubernetes",
				Type: URLTypeHTTP,
			}},
			{"http://192.168.101.127/user/project.git", GitURL{
				Hostname: "192.168.101.127",
				Path: "user/project",
				Type: URLTypeHTTP,
			}},
			{"http://github.com/user/project.git", GitURL{
				Hostname: "github.com",
				Path: "user/project",
				Type: URLTypeHTTP,
			}},
			{"http://host.xz/path/to/repo.git/", GitURL{
				Hostname: "host.xz",
				Path: "path/to/repo",
				Type: URLTypeHTTP,
			}},
			{"https://192.168.101.127/user/project.git", GitURL{
				Hostname: "192.168.101.127",
				Path: "user/project",
				Type: URLTypeHTTP,
			}},
			{"https://github.com/user/project.git", GitURL{
				Hostname: "github.com",
				Path: "user/project",
				Type: URLTypeHTTP,
			}},
			{"https://host.xz/path/to/repo.git/", GitURL{
				Hostname: "host.xz",
				Path: "path/to/repo",
				Type: URLTypeHTTP,
			}},
			{"https://username::;*%$:@github.com/username/repository.git", GitURL{
				Hostname: "github.com",
				Path: "username/repository",
				Type: URLTypeHTTP,
				Username: "username",
			}},
			{"https://username:$fooABC@:@github.com/username/repository.git", GitURL{
				Hostname: "github.com",
				Path: "username/repository",
				Type: URLTypeHTTP,
				Username: "username",
			}},
			{"https://username:password@github.com/username/repository.git", GitURL{
				Hostname: "github.com",
				Path: "username/repository",
				Type: URLTypeHTTP,
				Username: "username",
			}},
		}

		for _, tc := range testcases {
			gu := parse(tc.input, URLTypeHTTP)
			require.NotNil(t, gu, tc.input)

			assert.Equal(t, tc.output.Username, gu.Username, tc.input)
			assert.Equal(t, tc.output.Hostname, gu.Hostname, tc.input)
			assert.Equal(t, tc.output.Path, gu.Path, tc.input)
			assert.Equal(t, URLTypeHTTP, gu.Type, tc.input)
		}
	})

	t.Run("scp", func(t *testing.T) {

		testcases := []struct{
			input string
			output GitURL
		}{
			{"git@github.com:kubernetes/kubernetes.git", GitURL{
				Hostname: "github.com",
				Username: "git",
				Path: "kubernetes/kubernetes",
			}},
			{"git@gitlab.com:facebook/react.git", GitURL{
				Hostname: "gitlab.com",
				Username: "git",
				Path: "facebook/react",
			}},
			{"git@192.168.101.127:user/project.git", GitURL{
				Hostname: "192.168.101.127",
				Username: "git",
				Path: "user/project",
			}},
			{"git@github.com:user/project.git", GitURL{
				Hostname: "github.com",
				Username: "git",
				Path: "user/project",
			}},
			{"git@github.com:user/some-project.git", GitURL{
				Hostname: "github.com",
				Username: "git",
				Path: "user/some-project",
			}},
			{"git@github.com:user/some_project.git", GitURL{
				Hostname: "github.com",
				Username: "git",
				Path: "user/some_project",
			}},
		}

		for _, tc := range testcases {
			gu := parse(tc.input, URLTypeSCP)
			require.NotNil(t, gu, tc.input)

			assert.Equal(t, tc.output.Username, gu.Username, tc.input)
			assert.Equal(t, tc.output.Hostname, gu.Hostname, tc.input)
			assert.Equal(t, tc.output.Path, gu.Path, tc.input)
		}
	})
}

func TestIllegalParsing(t *testing.T) {

	t.Run("ssh", func(t *testing.T) {

		testcases := []string{
			"ssh://host./path/to/repo.git/",
			"ssh:://host.xz/~user/path/to/repo.git/",
			"ssh://host.xz:port/path/to/repo.git/",
			"ssh://@host.xz/path/to/repo.git/",
			"ssh://user@user@host.xz/path/to/repo.git/",
			"ssh://user@.xz/~/path/to/repo.git",
			"ssh://user@./~user/path/to/repo.git/",
			"ssh://@user@host.xz:1/path/to/repo.git/",
			"ssh://user@some.crazy.domain.net.au@:1/path/to/repo.git/",
			"ssh:/user@some.crazy.domain.net.au:1/a.git",
		}

		for _, tc := range testcases {
			gu := parse(tc, URLTypeSSH)
			require.Nil(t, gu, tc)
		}
	})
}
