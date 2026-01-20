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
	"github.com/raravena80/ya/common"
	"github.com/raravena80/ya/ops"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// scpCmd represents the scp command
var scpCmd = &cobra.Command{
	Use:   "scp [options] <source> <destination>",
	Short: "Copy files to multiple servers",
	Long: `Copy files to multiple servers.
You can specify the source and destination files,
the source files are local and the destination files
are in the remote servers.

Arguments:
  <source>       Source file or directory (local)
  <destination>  Destination file or directory (remote)

Example:
  ya scp -m host1,host2 -u user /path/to/file /remote/path`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		options := BuildCommonOptions()
		options = append(options,
			common.SetSource(args[0]))
		options = append(options,
			common.SetDestination(args[1]))
		options = append(options,
			common.SetIsRecursive(viper.GetBool("ya.scp.recursive")))
		options = append(options,
			common.SetOp("scp"))
		ops.SSHSession(options...)
	},
}

func init() {
	// Add scpCmd to cobra
	RootCmd.AddCommand(scpCmd)
	scpCmd.Flags().BoolP("recursive", "r", false, "Set recursive copy")
	viper.BindPFlag("ya.scp.recursive", scpCmd.Flags().Lookup("recursive"))
}
