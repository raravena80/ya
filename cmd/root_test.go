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
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestInitConfig(t *testing.T) {
	// Create a temporary directory for test config
	tmpDir, err := os.MkdirTemp("", "ya-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test config file
	configContent := `
ya:
  user: testuser
  port: 2222
`
	configFile := filepath.Join(tmpDir, ".ya.yaml")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Set cfgFile to point to our test config
	originalCfgFile := cfgFile
	cfgFile = configFile

	// Reset viper to ensure clean state
	viper.Reset()

	// Call initConfig
	initConfig()

	// Verify config was read
	if viper.GetString("ya.user") != "testuser" {
		t.Errorf("Expected user 'testuser', got '%s'", viper.GetString("ya.user"))
	}
	if viper.GetInt("ya.port") != 2222 {
		t.Errorf("Expected port 2222, got %d", viper.GetInt("ya.port"))
	}

	// Restore original cfgFile
	cfgFile = originalCfgFile
}

func TestInitConfigNoFile(t *testing.T) {
	// Set cfgFile to non-existent file
	originalCfgFile := cfgFile
	nonExistentFile := "/tmp/non-existent-ya-config-12345.yaml"
	cfgFile = nonExistentFile

	// Reset viper to ensure clean state
	viper.Reset()

	// Call initConfig - should not fail, just no config loaded
	initConfig()

	// Restore original cfgFile
	cfgFile = originalCfgFile
}

func TestInitConfigDefaultLocation(t *testing.T) {
	// Test default config file location (in home directory)
	originalCfgFile := cfgFile
	cfgFile = ""

	// Reset viper to ensure clean state
	viper.Reset()

	// Call initConfig - should use default location
	initConfig()

	// Restore original cfgFile
	cfgFile = originalCfgFile
}

func TestExecuteWithHelp(t *testing.T) {
	// Set up RootCmd with help
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Test that we can at least create the command and check its properties
	if RootCmd == nil {
		t.Error("RootCmd should not be nil")
	}

	// Verify RootCmd has expected properties
	if RootCmd.Use != "ya" {
		t.Errorf("Expected Use 'ya', got '%s'", RootCmd.Use)
	}

	// Note: Execute() calls os.Exit(1) on error, so we can't test it directly
	// The actual execution path is tested via integration tests
}

