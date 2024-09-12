package actions

import "antispam/structs"

func Mark(file string) *structs.Spammer {

	var result *structs.Spammer = nil

	// TODO: Parse E-Mail
	// TODO: Lookup Domains in HostMap
	// TODO: Lookup IPs in SpammerMap
	// TODO: If entry does not exist, result=&structs.NewSpammer()
	// TODO: Add missing Domains
	// TODO: Add missing IPs

	return result

}
