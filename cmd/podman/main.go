package main

import (
	"os"

	_ "github.com/containers/podman/v5/cmd/podman/completion"
	_ "github.com/containers/podman/v5/cmd/podman/containers"
	_ "github.com/containers/podman/v5/cmd/podman/generate"
	_ "github.com/containers/podman/v5/cmd/podman/healthcheck"
	_ "github.com/containers/podman/v5/cmd/podman/images"
	_ "github.com/containers/podman/v5/cmd/podman/machine"
	_ "github.com/containers/podman/v5/cmd/podman/manifest"
	_ "github.com/containers/podman/v5/cmd/podman/networks"
	_ "github.com/containers/podman/v5/cmd/podman/play"
	_ "github.com/containers/podman/v5/cmd/podman/pods"
	"github.com/containers/podman/v5/cmd/podman/registry"
	_ "github.com/containers/podman/v5/cmd/podman/secrets"
	_ "github.com/containers/podman/v5/cmd/podman/system"
	_ "github.com/containers/podman/v5/cmd/podman/volumes"
	"github.com/containers/podman/v5/pkg/rootless"
	"github.com/sirupsen/logrus"
)

func main() {
	// reexec is used to re-run the binary in certain contexts (e.g. rootless)
	// must be called before anything else in main
	if reexecDone := rootless.TryReexecAsRootless(); reexecDone {
		return
	}

	app := registry.PodmanConfig()

	if err := app.Execute(); err != nil {
		if registry.GetExitCode() == 0 {
			registry.SetExitCode(registry.ExecErrorCodeGeneric)
		}
		// Log at Error level so the message is visible without debug flags.
		logrus.Errorf("%v", err)
	}

	exitCode := registry.GetExitCode()

	// Always log the exit code at debug level for easier troubleshooting.
	logrus.Debugf("Podman exiting with code %d", exitCode)

	os.Exit(exitCode)
}
