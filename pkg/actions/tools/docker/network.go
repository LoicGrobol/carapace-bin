package docker

import (
	"strings"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace-bin/pkg/styles"
)

// ActionNetworks completes network names
//   bridge (bridge/local)
//   host (host/local)
func ActionNetworks() carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		return carapace.ActionExecCommand("docker", "network", "ls", "--format", "{{.Name}}\n{{.Driver}}/{{.Scope}}")(func(output []byte) carapace.Action {
			vals := strings.Split(string(output), "\n")
			return carapace.ActionValuesDescribed(vals[:len(vals)-1]...).Style(styles.Docker.Network)
		})
	})
}
