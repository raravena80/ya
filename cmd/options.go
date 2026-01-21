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
// from the viper configuration. It returns a slice of option functions that can be
// applied to the Options struct.
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

	// Optional timeout overrides
	if ct := viper.GetInt("ya.connect-timeout"); ct > 0 {
		options = append(options, common.SetConnectTimeout(ct))
	}
	if ct := viper.GetInt("ya.command-timeout"); ct > 0 {
		options = append(options, common.SetCommandTimeout(ct))
	}

	// Output format
	if fmt := viper.GetString("ya.output-format"); fmt != "" {
		options = append(options, common.SetOutputFormat(fmt))
	}

	// Dry-run mode
	if viper.GetBool("ya.dry-run") {
		options = append(options, common.SetDryRun(true))
	}

	// Host patterns
	if patterns := viper.GetStringSlice("ya.host-patterns"); len(patterns) > 0 {
		options = append(options, common.SetHostPatterns(patterns))
	}

	// Host excludes
	if excludes := viper.GetStringSlice("ya.host-excludes"); len(excludes) > 0 {
		options = append(options, common.SetHostExcludes(excludes))
	}

	// Progress indicators
	if viper.GetBool("ya.show-progress") {
		options = append(options, common.SetShowProgress(true))
	}

	return options
}
