package utils

import (
	"strings"

	"github.com/hashicorp/nomad/client/allocdir"
	"github.com/hashicorp/nomad/client/config"
	cstructs "github.com/hashicorp/nomad/client/structs"
	"github.com/hashicorp/nomad/client/taskenv"
	"github.com/hashicorp/nomad/nomad/structs"
)

// SetEnvvars sets path and host env vars depending on the FS isolation used.
func SetEnvvars(envBuilder *taskenv.Builder, fsi cstructs.FSIsolation, taskDir *allocdir.TaskDir, conf *config.Config) {
	// Set driver-specific environment variables
	switch fsi {
	case cstructs.FSIsolationNone:
		// Use host paths
		envBuilder.SetAllocDir(taskDir.SharedAllocDir)
		envBuilder.SetTaskLocalDir(taskDir.LocalDir)
		envBuilder.SetSecretsDir(taskDir.SecretsDir)
	default:
		// filesystem isolation; use container paths
		envBuilder.SetAllocDir(allocdir.SharedAllocContainerPath)
		envBuilder.SetTaskLocalDir(allocdir.TaskLocalContainerPath)
		envBuilder.SetSecretsDir(allocdir.TaskSecretsContainerPath)
	}

	// Set the host environment variables for non-image based drivers
	if fsi != cstructs.FSIsolationImage {
		filter := strings.Split(conf.ReadDefault("env.blacklist", config.DefaultEnvBlacklist), ",")
		envBuilder.SetHostEnvvars(filter)
	}
}

// CgroupsMounted returns true if the cgroups are mounted on a system otherwise
// returns false
func CgroupsMounted(node *structs.Node) bool {
	_, ok := node.Attributes["unique.cgroup.mountpoint"]
	return ok
}
