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
	"github.com/spf13/viper"
)

// BuildCommonOptions constructs the common options shared between SSH and SCP commands
func BuildCommonOptions() []func(*common.Options) {
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
	return options
}
