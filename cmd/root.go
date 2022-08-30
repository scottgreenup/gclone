package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/scottgreenup/gclone/pkg/parse"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RootCommandOutput struct {
	TargetDirectory string `json:"targetDirectory"`
}

func RootCommandRun(result *parse.TransformResult) (*RootCommandOutput, error) {

	// Set up the `git clone` command.
	gitArgs := result.GitArgs
	_, _ = fmt.Fprintf(os.Stderr, "git clone %s\n", strings.Join(gitArgs, " "))
	gitArgs = append([]string{"clone", "--progress"}, gitArgs...)
	command := exec.Command("git", gitArgs...)

	// Set up getting the output from `git clone`.
	commandReader, err := command.StderrPipe()
	if err != nil {
		return nil, err
	}

	// Run a goroutine to forward the command output.
	go func() {
		p := make([]byte, 1)
		for {
			_, err := commandReader.Read(p)

			// TODO handle non io.EOF errors.
			if err != nil {
				break
			}
			_, _ = fmt.Fprintf(os.Stderr, "%s", string(p))
		}
	}()

	// Run the command.
	if err := command.Start(); err != nil {
		return nil, err
	}

	// Wait for it to finish.
	if err := command.Wait(); err != nil {
		return nil, err
	}

	if err := os.Chdir(result.TargetDirectoryPath); err != nil {
		return nil, err
	}

	return &RootCommandOutput{
		TargetDirectory: result.TargetDirectoryPath,
	}, nil
}

var rootCmd = &cobra.Command{
	Use:   "gclone",
	Short: "An improved git cloning experience",
	Long: `A drop in replacement for git clone that allows for configuring the automatic organising of repositories that
are cloned.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		transformConfig := parse.DefaultTransformConfig()
		if s := viper.GetString("DefaultDirectory"); s != "" {
			transformConfig.DefaultDirectory = s
		}

		// TODO: How does the user discover this?
		transformConfig.FailOnExisting = viper.GetBool("FailOnExisting")

		result, err := parse.Transform(args, transformConfig)

		if err != nil {
			if err == parse.TransformErrorBadUsage {
				if err := cmd.Help(); err != nil {
					_, _ = fmt.Fprintln(os.Stderr, err)
				}
				os.Exit(1)
			}

			return err
		}

		// The directory may already exist.
		if fi, err := os.Stat(result.TargetDirectoryPath); err == nil {

			if fi.IsDir() && transformConfig.FailOnExisting {
				return errors.Errorf("%s already exists, can not clone to it", result.TargetDirectoryPath)
			}

			if fi.IsDir() {
				resp := &RootCommandOutput{
					TargetDirectory: result.TargetDirectoryPath,
				}
				if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
					return err
				}
				return nil
			}

			return err
		}

		resp, err := RootCommandRun(result)

		if err != nil {
			return err
		}

		return json.NewEncoder(os.Stdout).Encode(resp)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.AddConfigPath("$HOME/.config/gclone")
	viper.AddConfigPath("/etc/gclone")
	viper.SetConfigName("config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using config file: ", viper.ConfigFileUsed())
	}
}
