package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace-bin/pkg/actions/net"
	"github.com/spf13/cobra"
)

var address_addCmd = &cobra.Command{
	Use:   "add",
	Short: "add new protocol address",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	carapace.Gen(address_addCmd).Standalone()

	addressCmd.AddCommand(address_addCmd)

	carapace.Gen(address_addCmd).PositionalCompletion(
		net.ActionIpv4Addresses(), // TODO ip6
		carapace.ActionValues("dev"),
		net.ActionDevices(net.IncludedDevices{Wifi: true, Ethernet: true}),
	)
}
