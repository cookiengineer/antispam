package actions

import "antispam/structs"

func Classify(file string) *structs.Spammer {

	var result *structs.Spammer = nil

	// TODO: Parse E-Mail
	// TODO: Lookup Domains in HostMap
	// TODO: Lookup IPs in SpammerMap

	return result

}
