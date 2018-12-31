package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/scottgreenup/gclone/pkg/parse"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

		gitArgs, err := parse.Transform(args, transformConfig)

		if err != nil {
			if err == parse.TransformErrorBadUsage {
				cmd.Help()
				os.Exit(1)
			}

			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("git clone %s\n", strings.Join(gitArgs, " "))

		gitArgs = append([]string{"clone", "--progress"}, gitArgs...)
		command := exec.Command("git", gitArgs...)

		commandReader, err := command.StderrPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		go func() {
			p := make([]byte, 1)
			for {
				_, err := commandReader.Read(p)

				if err != nil {
					if err == io.EOF {
						break
					} else {
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}
				}

				fmt.Printf("%s", string(p))
			}
		}()

		if err := command.Start(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := command.Wait(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
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
