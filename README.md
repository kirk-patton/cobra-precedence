# cobra-precedence
Sample Cobra Project demonstrating cobra & viper precedence

This sample Cobra project demonstrates how to load values for multiple locations.

* flag default value
* config file
* environment variable
* explicit call to --flag 

The project uses Viper as the source of truth for the values collected.  It was written to better understand Cobra.
The overall idea is that cobra/viper logic is separate from the core logic of whatever the client app is supposed to do.  That logic should be contained in a separate file/library with its own unit tests.  Cobra/viper just collects the settings to pass along.  This test program does not do anything useful other that demonstrate precedence.

The tests start with no flags or config file and add them one test at a time.

### Notes

When testing, flags/arguments are passed to the "rootCmd". The same holds true for calls to rootCmd.Execute.  This is the same flow as a cli application would use, the difference being that we supply arguments using [cmd.SetArgs([]string)](https://pkg.go.dev/github.com/spf13/cobra#Command.SetArgs) which allows to control the argument during testing.

When setting up a cobra "command", we use [RunE](https://git.vzbuilders.com/kpatton/cobra-precedence/blob/main/cmd/run.go#L49-L51). It allows flexibility to return and check errors from the command.

The [output](https://git.vzbuilders.com/kpatton/cobra-precedence/blob/main/cmd/run_test.go#L117-L118) of the RunE command is set to return an iowriter so that the output can be captured in a buffer for [examination](https://git.vzbuilders.com/kpatton/cobra-precedence/blob/main/cmd/run_test.go#L133) but the test program.

The "command" we have specified [configures](https://git.vzbuilders.com/kpatton/cobra-precedence/blob/main/cmd/run.go#L49-L51) viper.

### Reference

* [How to test CLI commands made with Go and Cobra](https://gianarb.it/blog/golang-mockmania-cli-command-with-cobra)