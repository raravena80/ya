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

var (
	cfgFile   string
	user      string
	key       string
	port      int
	timeout   int
	agentsock string
	machines  []string
	Version   string
	Gitcommit string
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
		fmt.Println(err)
		//go:nocovline // os.Exit hard to test in unit tests
		os.Exit(1)
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
			fmt.Printf("Error: could not determine home directory: %v\n", err)
			//go:nocovline // os.Exit hard to test in unit tests
			os.Exit(1)
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
			fmt.Println(err)
			//go:nocovline // os.Exit hard to test in unit tests
			os.Exit(1)
		}

		// Search config in home directory with name ".ya" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ya")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
