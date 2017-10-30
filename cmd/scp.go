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
	"fmt"

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
		fmt.Println("scp not implemented yet")
	},
}

func init() {
	// Add scpCmd to cobra
	RootCmd.AddCommand(scpCmd)
	scpCmd.Flags().StringVarP(&src, "src", "s", "", "Source file or directory")
	viper.BindPFlag("ya.scp.src", sshCmd.Flags().Lookup("source"))
	scpCmd.Flags().StringVarP(&dst, "dst", "d", "", "Destination file or directory")
	viper.BindPFlag("ya.scp.dst", sshCmd.Flags().Lookup("destination"))
}
