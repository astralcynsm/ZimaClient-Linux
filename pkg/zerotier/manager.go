package zerotier

import (
	"rnctl/pkg/utils"
)

func GetInfo() (string, error) {
	return utils.RunCommand("zerotier-cli", "info", "-j")
}

func JoinNetwork(nwid string) (string, error) {
	return utils.RunCommand("zerotier-cli", "join", nwid)
}

func LeaveNetwork(nwid string) (string, error) {
	return utils.RunCommand("zerotier-cli", "leave", nwid)
}

func ListNetworks() (string, error) {
	return utils.RunCommand("zerotier-cli", "listnetworks", "-j")
}
