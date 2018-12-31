package parse

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/scottgreenup/gclone/pkg/git"
)

type TransformConfig struct {
	DefaultDirectory string
}

func DefaultTransformConfig() TransformConfig {
	return TransformConfig{
		DefaultDirectory: "~/code",
	}
}

type Error string

func (e Error) Error() string {
	return string(e)

}

const (
	TransformErrorBadUsage Error = "BadUsage"
)

func Transform(args []string, config TransformConfig) ([]string, error) {
	gitArgs := make([]string, 0, len(args))
	doubleDash := false

	ourArgs := make([]string, 0, 2)

	for _, arg := range args {
		if arg == "--" {
			doubleDash = true
			continue
		} else if arg == "--help" || arg == "-h" {
			return nil, TransformErrorBadUsage
		} else if doubleDash {
			gitArgs = append(gitArgs, arg)
		} else {
			ourArgs = append(ourArgs, arg)
		}
	}

	if len(ourArgs) == 0 || len(ourArgs) > 2 {
		return nil, TransformErrorBadUsage
	}

	switch len(ourArgs) {
	case 1:
		// if args == 1 then we treat it as a repo
		//      parse repo
		//      mkdir -p the auto-generated dir
		//      git clone with repo and AGD
		repo := ourArgs[0]
		gu, err := git.NewURL(repo)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// TODO normalise, spaces to -, etc...
		path := filepath.Join(config.DefaultDirectory, gu.Hostname, gu.Path)
		expandedPath, expandErr := homedir.Expand(path)
		if expandErr != nil {
			return nil, expandErr
		}
		result := gitArgs[:]
		result = append(result, repo, expandedPath)
		return result, nil

	case 2:
		// if args == 2 then repo and directory
		//      mkdir -p the directory
		//      call git clone with repo and directory
		//      if it fails, kill the directory
		repo := ourArgs[0]
		path := ourArgs[1]
		expandedPath, expandErr := homedir.Expand(path)
		if expandErr != nil {
			return nil, expandErr
		}
		result := gitArgs[:]
		result = append(result, repo, expandedPath)
		return result, nil

	default:
		return nil, TransformErrorBadUsage
	}
}
