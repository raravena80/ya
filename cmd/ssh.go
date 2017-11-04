// Copyright © 2017 Ricardo Aravena <raravena@branch.io>
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
	"github.com/raravena80/ya/common"
	"github.com/raravena80/ya/ops"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var command string

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Run command acrosss multiple servers",
	Long: `Run a command across multiple servers,
using SSH.`,
	Run: func(cmd *cobra.Command, args []string) {
		var options []func(*common.Options)
		options = append(options,
			common.SetMachines(viper.GetStringSlice("ya.machines")))
		options = append(options,
			common.SetUser(viper.GetString("ya.user")))
		options = append(options,
			common.SetPort(viper.GetInt("ya.port")))
		options = append(options,
			common.SetCmd(viper.GetString("ya.command")))
		options = append(options,
			common.SetKey(viper.GetString("ya.key")))
		options = append(options,
			common.SetUseAgent(viper.GetBool("ya.useagent")))
		options = append(options,
			common.SetTimeout(viper.GetString("ya.timeout")))
		ops.Run(options...)
	},
}

func init() {
	RootCmd.AddCommand(sshCmd)

	// Local flags
	sshCmd.Flags().StringVarP(&command, "command", "c", "", "Command to run")
	viper.BindPFlag("ya.ssh.command", sshCmd.Flags().Lookup("command"))
}
