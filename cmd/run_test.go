package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// TODO
// write the --config file from a var in the test
// the test is failing because we are finding the config file in the search path
// we need to fix this...
//
// helper need to be called and defer file cleanup

func TestPrecedence(t *testing.T) {
	// set up config file if needed
	configFH, err := config()
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	if configFH != nil {
		defer configFH.Close()
		defer os.RemoveAll(configFH.Name())
	}

	tests := []struct {
		name          string
		root          *cobra.Command
		command       *cobra.Command
		args          []string // the command line flag and args
		expectedValue string
		runE          func(*cobra.Command, []string) error
		environment   map[string]string
	}{
		{
			name:          "When we have no config and call the run command",
			root:          rootCmd,
			command:       runCmd,
			args:          []string{"run"},
			expectedValue: "--precedence => default",
			runE: func(c *cobra.Command, args []string) error {
				out := c.OutOrStdout()
				fmt.Fprintf(out, "--precedence => %s", viper.GetString("precedence"))
				return nil
			},
		},
		{
			name:          "When we have a config and call the run command",
			root:          rootCmd,
			command:       runCmd,
			args:          []string{"run", "--config", configFH.Name()},
			expectedValue: "--precedence => config",
			runE: func(c *cobra.Command, args []string) error {
				out := c.OutOrStdout()
				fmt.Fprintf(out, "--precedence => %s", viper.GetString("precedence"))
				// spew.Config.Dump(viper.AllSettings()) //Debug
				return nil
			},
		},
		{
			name:    "When we have a config and call the run command and the setting is in env",
			root:    rootCmd,
			command: runCmd,
			//args:          []string{"run", "--config", "/Users/kpatton/.cobra-precedence/config.yaml"},
			args:          []string{"run"},
			expectedValue: "--precedence => environment",
			runE: func(c *cobra.Command, args []string) error {
				out := c.OutOrStdout()
				fmt.Fprintf(out, "--precedence => %s", viper.GetString("precedence"))
				return nil
			},
			environment: map[string]string{"precedence": "environment"},
		},
		{
			name:    "When we have a config, call the run command, the setting is in env and we set the flag",
			root:    rootCmd,
			command: runCmd,
			//args:          []string{"run", "--config", "/Users/kpatton/.cobra-precedence/config.yaml"},
			args:          []string{"run", "--precedence", "flag"},
			expectedValue: "--precedence => flag",
			runE: func(c *cobra.Command, args []string) error {
				out := c.OutOrStdout()
				fmt.Fprintf(out, "--precedence => %s", viper.GetString("precedence"))
				return nil
			},
			environment: map[string]string{"precedence": "environment"},
		},
		{
			name:    "When we have a config access the value using \".\" notation",
			root:    rootCmd,
			command: runCmd,
			//args:          []string{"run", "--config", "/Users/kpatton/.cobra-precedence/config.yaml"},
			args:          []string{"run", "--precedence", "flag"},
			expectedValue: "nested value: baz",
			runE: func(c *cobra.Command, args []string) error {
				out := c.OutOrStdout()
				fmt.Fprintf(out, "nested value: %s", viper.GetString("foo.bar"))
				return nil
			},
			environment: map[string]string{"precedence": "environment"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// set the command output to our buffer so we can capute the output
			output := &bytes.Buffer{}
			tt.command.SetOut(output)

			// pass the arguments to the rootCmd
			tt.root.SetArgs(tt.args)

			// replace the RunE on "command" with our test version
			tt.command.RunE = tt.runE

			// export any defined environment vars
			exportEnv(tt.environment)

			// execute the rootCmd
			err := tt.root.Execute()
			assert.NoError(t, err)

			got, _ := ioutil.ReadAll(output)
			assert.Equal(t, tt.expectedValue, string(got))
		})
	}
}

func exportEnv(vars map[string]string) {
	for k, v := range vars {
		k = strings.Join([]string{cfgPrefix, k}, "_")
		k := strings.ToUpper(k)
		os.Setenv(k, v)
	}
}

// config create a yaml file with configuration settings to load
func config() (*os.File, error) {
	settings := []byte(`
---
precedence: config
foo:
  bar: baz
`)

	cfg, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, fmt.Errorf("error creating config: %w", err)
	}

	newName := cfg.Name() + ".yaml"

	err = os.Rename(cfg.Name(), newName)
	if err != nil {
		return nil, fmt.Errorf("error renaming config: %w", err)
	}
	cfg.Close()

	fh, err := os.Open(newName)
	if err != nil {
		return nil, fmt.Errorf("error opeing renamed config: %w", err)
	}

	err = ioutil.WriteFile(fh.Name(), settings, 0700)
	if err != nil {
		return nil, fmt.Errorf("error writing renamed config: %w", err)
	}

	return fh, err
}
