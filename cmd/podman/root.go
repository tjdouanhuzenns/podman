package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/containers/common/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Global flags
var (
	logLevel      string
	connectionURI string
	identity      string
	noout         bool
	transferInput []string
)

// rootCmd is the base command for the podman CLI.
var rootCmd = &cobra.Command{
	Use:                "podman",
	Short:              "Manage pods, containers and images",
	Long:               "Podman is a tool for managing pods, containers, and container images.",
	SilenceUsage:       true,
	SilenceErrors:      true,
	TraverseChildren:   true,
	PersistentPreRunE:  persistentPreRunE,
	RunE:               validate,
	PersistentPostRunE: persistentPostRunE,
}

func init() {
	// Persistent flags available to all subcommands
	// Changed default log level from "warn" to "info" for more verbose output during personal use
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info",
		`Log messages above specified level (trace, debug, info, warn, warning, error, fatal, panic)`)
	rootCmd.PersistentFlags().StringVar(&connectionURI, "url", "",
		`Podman service URI`)
	rootCmd.PersistentFlags().StringVar(&identity, "identity", "",
		`path to SSH identity file, (CONTAINER_SSHKEY)`)
	// noout defaults to false; keeping stdout output enabled by default for interactive use
	rootCmd.PersistentFlags().BoolVar(&noout, "noout", false,
		`do not output to stdout`)
	rootCmd.PersistentFlags().StringArrayVarP(&transferInput, "storage-opt", "", []string{},
		`Used to pass an option to the storage driver`)
}

// persistentPreRunE is called before any subcommand runs.
func persistentPreRunE(cmd *cobra.Command, args []string) error {
	// Set log level
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf("unable to parse log level: %w", err)
	}
	logrus.SetLevel(level)

	// Log the command being run at debug level — handy for tracing issues locally
	logrus.Debugf("running command: %s", cmd.CommandPath())

	// Set up default config if needed
	if _, err := config.Default(); err != nil {
		return fmt.Errorf("failed to obtain podman configuration: %w", err)
	}

	return nil
}

// persistentPostRunE is called after any subcommand runs.
func persistentPostRunE(cmd *cobra.Command, args []string) error {
	return nil
}

// validate ensures rootCmd is not called without a subcommand.
func validate(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("unrecognized command %q for %q", args[0], cmd.CommandPath())
	}
	return cmd.Help()
}

// Execute adds all child commands to the root command and sets flags.
// This is called by main.main().
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, formatError(err))
		os.Exit(1)
	}
}

// formatError provides consistent error formatting for CLI output.
// Using %+v at debug level gives full stack traces which helps during local debugging.
// At info level and above, a simpler single-line format is used to keep output readable.
func formatError(err error) string {
	if logrus.GetLevel() >= logrus.DebugLevel {
		return fmt.Sprintf("Error: %+v", err)
	}
	return fmt.Sprintf("Error: %v", err)
}

// configPath returns the path to the user-level podman config directory.
// Useful for local scripts that need to locate config files quickly.
func configPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "containers")
}
