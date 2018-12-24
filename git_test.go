package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertURLType(t *testing.T, expected URLType, actual URLType, url string) {
	assert.True(t, expected == actual, fmt.Sprintf("Expected %s but received %s from %q", expected, actual, url))
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
			gu, err := NewGitURL(tc)
			require.NoError(t, err, tc)
			assertURLType(t, URLTypeGit, gu.Type, tc)
		}
	})

	t.Run("https", func(t *testing.T) {

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
			gu, err := NewGitURL(tc)
			require.NoError(t, err)
			assertURLType(t, URLTypeHTTP, gu.Type, tc)
		}
	})

	t.Run("scp", func(t *testing.T) {

		testcases := []string{
			"git@github.com:kubernetes/kubernetes.git",
			"git@gitlab.com:facebook/react.git",
			"git@192.168.101.127:user/project.git",
			"git@github.com:user/project.git",
			"git@github.com:user/some-project.git",
			"git@github.com:user/some-project.git",
			"git@github.com:user/some_project.git",
			"git@github.com:user/some_project.git",
		}

		for _, tc := range testcases {
			gu, err := NewGitURL(tc)
			require.NoError(t, err)
			assertURLType(t, URLTypeSCP, gu.Type, tc)
		}
	})

	t.Run("ssh", func(t *testing.T) {

		testcases := []string{
			"ssh://host.xz/path/to/repo.git/",
			"ssh://host.xz/path/to/repo.git/",
			"ssh://host.xz/~/path/to/repo.git",
			"ssh://host.xz/~user/path/to/repo.git/",
			"ssh://host.xz:port/path/to/repo.git/",
			"ssh://user@host.xz/path/to/repo.git/",
			"ssh://user@host.xz/path/to/repo.git/",
			"ssh://user@host.xz/~/path/to/repo.git",
			"ssh://user@host.xz/~user/path/to/repo.git/",
			"ssh://user@host.xz:port/path/to/repo.git/",
		}

		for _, tc := range testcases {
			gu, err := NewGitURL(tc)
			require.NoError(t, err)
			assertURLType(t, URLTypeSSH, gu.Type, tc)
		}
	})
}
