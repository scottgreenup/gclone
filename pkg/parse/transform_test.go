package parse

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func doTransform(s string, config TransformConfig) ([]string, error) {
	args := strings.Split(s, " ")
	return Transform(args, config)
}

func TestTransform(t *testing.T) {
	t.Run("help", func(t *testing.T) {
		gitArgs, err := doTransform("--help", DefaultTransformConfig())
		require.Error(t, err)
		assert.Empty(t, gitArgs)
	})

	t.Run("repository", func(t *testing.T) {
		config := DefaultTransformConfig()
		gitArgs, err := doTransform("https://github.com/kubernetes/kubernetes.git", config)
		require.NoError(t, err)
		assert.True(t, strings.HasSuffix(gitArgs[len(gitArgs)-1], "github.com/kubernetes/kubernetes"))
	})

	t.Run("both", func(t *testing.T) {
		config := DefaultTransformConfig()
		gitArgs, err := doTransform("https://github.com/kubernetes/kubernetes.git someDir", config)
		require.NoError(t, err)
		assert.Equal(t, []string{
			"https://github.com/kubernetes/kubernetes.git",
			"someDir",
		}, gitArgs)
	})

	t.Run("both with flags", func(t *testing.T) {
		config := DefaultTransformConfig()
		gitArgs, err := doTransform("https://github.com/kubernetes/kubernetes.git someDir -- -q --something --no-tags", config)
		require.NoError(t, err)
		assert.Equal(t, []string{
			"--depth",
			"1",
			"-q",
			"--no-tags",
			"https://github.com/kubernetes/kubernetes.git",
			"someDir",
		}, gitArgs)
	})

	t.Run("flags before", func(t *testing.T) {
		config := DefaultTransformConfig()
		_, err := doTransform("--something https://github.com/kubernetes/kubernetes.git someDir -- -q --something --no-tags", config)
		require.Error(t, err)
	})
}
}
