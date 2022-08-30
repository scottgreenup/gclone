package parse

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/scottgreenup/gclone/pkg/git"
)

type TransformConfig struct {
	DefaultDirectory string
	FailOnExisting   bool
}

func DefaultTransformConfig() TransformConfig {
	return TransformConfig{
		DefaultDirectory: "~/code",
		FailOnExisting:   false,
	}
}

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	TransformErrorBadUsage Error = "BadUsage"
)

type TransformResult struct {
	GitArgs             []string
	TargetDirectoryPath string
}

func Transform(args []string, config TransformConfig) (*TransformResult, error) {
	gitArgs := make([]string, 0, len(args))
	ourArgs := make([]string, 0, 2)

	processingGitArgs := false
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			processingGitArgs = true
		}

		if processingGitArgs {
			gitArgs = append(gitArgs, arg)
		} else {
			ourArgs = append(ourArgs, arg)
		}
	}

	// TODO: Do this using cobra?
	// We're expecting 1 or 2 arguments.
	//  ourArgs[0] should be repo
	//  ourArgs[1] should be path
	if len(ourArgs) == 0 || len(ourArgs) > 2 {
		return nil, TransformErrorBadUsage
	}

	switch len(ourArgs) {
	case 1:
		repo := ourArgs[0]
		gu, err := git.NewURL(repo)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
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
		return &TransformResult{
			GitArgs:             result,
			TargetDirectoryPath: expandedPath,
		}, nil

	case 2:
		repo := ourArgs[0]
		path := ourArgs[1]
		expandedPath, expandErr := homedir.Expand(path)
		if expandErr != nil {
			return nil, expandErr
		}
		result := gitArgs[:]
		result = append(result, repo, expandedPath)
		return &TransformResult{
			GitArgs:             result,
			TargetDirectoryPath: expandedPath,
		}, nil

	default:
		return nil, TransformErrorBadUsage
	}
}
