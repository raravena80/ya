// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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

var (
	src string
	dst string
)

// scpCmd represents the scp command
var scpCmd = &cobra.Command{
	Use:   "scp",
	Short: "Copy files to multiple servers",
	Long: `Copy files to multiple servers.
You can specify the source and destination files,
the source files are local and the destination files
are in the remote servers.`,
	Run: func(cmd *cobra.Command, args []string) {
		var options []func(*common.Options)
		options = append(options,
			common.SetMachines(viper.GetStringSlice("ya.machines")))
		options = append(options,
			common.SetUser(viper.GetString("ya.user")))
		options = append(options,
			common.SetPort(viper.GetInt("ya.port")))
		options = append(options,
			common.SetKey(viper.GetString("ya.key")))
		options = append(options,
			common.SetUseAgent(viper.GetBool("ya.useagent")))
		options = append(options,
			common.SetTimeout(viper.GetInt("ya.timeout")))
		options = append(options,
			common.SetSource(viper.GetString("ya.source")))
		options = append(options,
			common.SetDestination(viper.GetString("ya.destination")))
		options = append(options,
			common.SetOp("scp"))
		ops.SSHSession(options...)
	},
}

func init() {
	// Add scpCmd to cobra
	RootCmd.AddCommand(scpCmd)
	scpCmd.Flags().StringVarP(&src, "src", "f", "", "Source file or directory")
	viper.BindPFlag("ya.scp.src", sshCmd.Flags().Lookup("source"))
	scpCmd.Flags().StringVarP(&dst, "dst", "d", "", "Destination file or directory")
	viper.BindPFlag("ya.scp.dst", sshCmd.Flags().Lookup("destination"))
}
