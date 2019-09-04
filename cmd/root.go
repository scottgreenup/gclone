package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
	Run: func(cmd *cobra.Command, args []string) {
		transformConfig := parse.DefaultTransformConfig()
		if s := viper.GetString("DefaultDirectory"); s != "" {
			transformConfig.DefaultDirectory = s
		}

		result, err := parse.Transform(args, transformConfig)

		if err != nil {
			if err == parse.TransformErrorBadUsage {
				if err := cmd.Help(); err != nil {
					_, _ = fmt.Fprintln(os.Stderr, err)
				}
				os.Exit(1)
			}

			fmt.Println(err)
			os.Exit(1)
		}

		resp, err := RootCommandRun(result)

		if err != nil {
			fmt.Println("%s", err.Error())
			os.Exit(1)
		}

		if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
			fmt.Println("%s", err.Error())
		}
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
		fmt.Println("Using config file: ", viper.ConfigFileUsed())
	}
}
