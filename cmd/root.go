// Copyright © 2017 Ricardo Aravena <raravena80@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// exitFunc is a function that exits the program.
// In production, this is os.Exit, but can be replaced for testing.
var exitFunc = os.Exit

// printfFunc is a function that prints to stdout.
// In production, this is fmt.Printf, but can be replaced for testing.
var printfFunc = fmt.Printf

// printlnFunc is a function that prints to stdout with a newline.
// In production, this is fmt.Println, but can be replaced for testing.
var printlnFunc = fmt.Println

var (
	cfgFile       string
	user          string
	key           string
	port          int
	timeout       int
	connectTimeout int
	commandTimeout int
	agentsock     string
	machines      []string
	Version       string
	Gitcommit     string
	outputFormat  string
	dryRun        bool
	hostPatterns  []string
	hostExcludes  []string
	showProgress  bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ya",
	Short: "Ya runs commands or copies files across multiple servers",
	Long: `Ya runs commands or copies files or directories,
across multiple servers, using SSH or SCP`,
	Version: fmt.Sprintf("%v\ncommit %v", Version, Gitcommit),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		printlnFunc(err)
		exitFunc(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	curUser := os.Getenv("LOGNAME")
	home := os.Getenv("HOME")
	if home == "" {
		var err error
		home, err = homedir.Dir()
		if err != nil {
			printfFunc("Error: could not determine home directory: %v\n", err)
			exitFunc(1)
		}
	}
	sshKey := home + "/.ssh/id_rsa"

	// Persistent flags
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ya.yaml)")
	RootCmd.PersistentFlags().StringSliceVarP(&machines, "machines", "m", []string{}, "Hosts to run command on")
	viper.BindPFlag("ya.machines", RootCmd.PersistentFlags().Lookup("machines"))
	RootCmd.PersistentFlags().IntVarP(&port, "port", "p", 22, "Ssh port to connect to")
	viper.BindPFlag("ya.port", RootCmd.PersistentFlags().Lookup("port"))
	RootCmd.PersistentFlags().StringVarP(&user, "user", "u", curUser, "User to run the command as")
	viper.BindPFlag("ya.user", RootCmd.PersistentFlags().Lookup("user"))
	RootCmd.PersistentFlags().StringVarP(&key, "key", "k", sshKey, "Ssh key to use for authentication, full path")
	viper.BindPFlag("ya.key", RootCmd.PersistentFlags().Lookup("key"))
	RootCmd.PersistentFlags().BoolP("useagent", "a", false, "Use agent for authentication")
	viper.BindPFlag("ya.useagent", RootCmd.PersistentFlags().Lookup("useagent"))
	RootCmd.PersistentFlags().IntVarP(&timeout, "timeout", "t", 5, "Timeout for connection")
	viper.BindPFlag("ya.timeout", RootCmd.PersistentFlags().Lookup("timeout"))
	RootCmd.PersistentFlags().StringVarP(&agentsock, "agentsock", "s", os.Getenv("SSH_AUTH_SOCK"), "SSH agent socket file. If using SSH agent")
	viper.BindPFlag("ya.agentsock", RootCmd.PersistentFlags().Lookup("agentsock"))
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Set verbose output")
	viper.BindPFlag("ya.verbose", RootCmd.PersistentFlags().Lookup("verbose"))
	RootCmd.PersistentFlags().IntVar(&connectTimeout, "connect-timeout", 0, "Connection timeout override in seconds")
	viper.BindPFlag("ya.connect-timeout", RootCmd.PersistentFlags().Lookup("connect-timeout"))
	RootCmd.PersistentFlags().IntVar(&commandTimeout, "command-timeout", 0, "Command execution timeout override in seconds")
	viper.BindPFlag("ya.command-timeout", RootCmd.PersistentFlags().Lookup("command-timeout"))
	RootCmd.PersistentFlags().StringVarP(&outputFormat, "output-format", "o", "text", "Output format: text, json, yaml, table")
	viper.BindPFlag("ya.output-format", RootCmd.PersistentFlags().Lookup("output-format"))
	RootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "n", false, "Preview operations without executing")
	viper.BindPFlag("ya.dry-run", RootCmd.PersistentFlags().Lookup("dry-run"))
	RootCmd.PersistentFlags().StringSliceVarP(&hostPatterns, "host", "H", []string{}, "Host patterns to match")
	viper.BindPFlag("ya.host-patterns", RootCmd.PersistentFlags().Lookup("host"))
	RootCmd.PersistentFlags().StringSliceVar(&hostExcludes, "host-exclude", []string{}, "Host patterns to exclude")
	viper.BindPFlag("ya.host-excludes", RootCmd.PersistentFlags().Lookup("host-exclude"))
	RootCmd.PersistentFlags().BoolVarP(&showProgress, "progress", "P", false, "Show progress indicators for file transfers")
	viper.BindPFlag("ya.show-progress", RootCmd.PersistentFlags().Lookup("progress"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			printlnFunc(err)
			exitFunc(1)
		}

		// Search config in home directory with name ".ya" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ya")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		printlnFunc("Using config file:", viper.ConfigFileUsed())
	}
}
