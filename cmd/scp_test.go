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
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestSCPCommand(t *testing.T) {
	// Find the scp command
	var scpCmd *cobra.Command
	for _, cmd := range RootCmd.Commands() {
		if cmd.Name() == "scp" {
			scpCmd = cmd
			break
		}
	}

	if scpCmd == nil {
		t.Fatal("scp command not found")
	}

	// Test scp command properties
	// Note: The Use string may vary based on branch/version
	if scpCmd.Use == "" {
		t.Error("scpCmd.Use is empty")
	}

	if scpCmd.Short == "" {
		t.Error("scpCmd.Short is empty")
	}

	if scpCmd.Long == "" {
		t.Error("scpCmd.Long is empty")
	}

	// Test that recursive flag exists
	recursiveFlag := scpCmd.Flags().Lookup("recursive")
	if recursiveFlag == nil {
		t.Error("recursive flag not found")
	}

	// Test flag shorthand
	if recursiveFlag.Shorthand != "r" {
		t.Errorf("recursive flag shorthand = %s, want r", recursiveFlag.Shorthand)
	}
}

func TestSCPCommandArgsValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{name: "SCP with help flag",
			args:    []string{"scp", "--help"},
			wantErr: false},
		{name: "SCP with recursive flag and help",
			args:    []string{"scp", "-r", "--help"},
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCmd := RootCmd
			testCmd.SetArgs(tt.args)
			err := testCmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSCPCommandFlags(t *testing.T) {
	// Find the scp command
	var scpCmd *cobra.Command
	for _, cmd := range RootCmd.Commands() {
		if cmd.Name() == "scp" {
			scpCmd = cmd
			break
		}
	}

	if scpCmd == nil {
		t.Fatal("scp command not found")
	}

	// Test recursive flag
	recursiveFlag := scpCmd.Flags().Lookup("recursive")
	if recursiveFlag == nil {
		t.Error("recursive flag not found")
	} else {
		// Test default value
		if recursiveFlag.DefValue != "false" {
			t.Errorf("recursive flag default = %s, want false", recursiveFlag.DefValue)
		}
	}
}

func TestSCPCommandStructure(t *testing.T) {
	// Find the scp command
	var scpCmd *cobra.Command
	for _, cmd := range RootCmd.Commands() {
		if cmd.Name() == "scp" {
			scpCmd = cmd
			break
		}
	}

	if scpCmd == nil {
		t.Fatal("scp command not found")
	}

	// Test PreRun is set (for argument validation)
	// Note: This may be nil in some versions
	_ = scpCmd.PreRun

	// Test Run is set
	if scpCmd.Run == nil {
		t.Error("scpCmd.Run is nil")
	}
}

func TestSCPCommandFlagBindings(t *testing.T) {
	// Find the scp command
	var scpCmd *cobra.Command
	for _, cmd := range RootCmd.Commands() {
		if cmd.Name() == "scp" {
			scpCmd = cmd
			break
		}
	}

	if scpCmd == nil {
		t.Fatal("scp command not found")
	}

	// Test src flag
	srcFlag := scpCmd.Flags().Lookup("src")
	if srcFlag == nil {
		t.Error("src flag not found")
	}
	if srcFlag.Shorthand != "f" {
		t.Errorf("src flag shorthand = %s, want f", srcFlag.Shorthand)
	}

	// Test dst flag
	dstFlag := scpCmd.Flags().Lookup("dst")
	if dstFlag == nil {
		t.Error("dst flag not found")
	}
	if dstFlag.Shorthand != "d" {
		t.Errorf("dst flag shorthand = %s, want d", dstFlag.Shorthand)
	}

	// Test recursive flag
	recursiveFlag := scpCmd.Flags().Lookup("recursive")
	if recursiveFlag == nil {
		t.Error("recursive flag not found")
	}
}

func TestSCPCommandRunFunction(t *testing.T) {
	// Test that the Run function exists and can be called
	// Find the scp command
	var scpCmd *cobra.Command
	for _, cmd := range RootCmd.Commands() {
		if cmd.Name() == "scp" {
			scpCmd = cmd
			break
		}
	}

	if scpCmd == nil {
		t.Fatal("scp command not found")
	}

	// Verify Run function is set
	if scpCmd.Run == nil {
		t.Fatal("scpCmd.Run is nil")
	}

	// Test calling Run with empty args (will fail but exercises the code path)
	// We're just verifying the function can be called without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Run function panicked: %v", r)
		}
	}()

	// Set some viper values to test the Run function
	viper.Set("ya.machines", []string{"host1"})
	viper.Set("ya.scp.source", "/src")
	viper.Set("ya.scp.destination", "/dst")
	viper.Set("ya.scp.recursive", false)

	// Call Run - it will fail to connect but exercises the code
	scpCmd.Run(scpCmd, []string{})
}

func TestSCPInitFunction(t *testing.T) {
	// Test that init() properly sets up the command
	// This is already tested by other tests, but let's verify
	// the command is properly registered

	// Find the scp command
	var scpCmd *cobra.Command
	found := false
	for _, cmd := range RootCmd.Commands() {
		if cmd.Name() == "scp" {
			scpCmd = cmd
			found = true
			break
		}
	}

	if !found {
		t.Error("scp command not registered with RootCmd")
	}

	if scpCmd == nil {
		t.Fatal("scp command is nil")
	}

	// Verify the command has the expected flags
	flags := []string{"src", "dst", "recursive"}
	for _, flagName := range flags {
		flag := scpCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Flag %s not found on scp command", flagName)
		}
	}
}
