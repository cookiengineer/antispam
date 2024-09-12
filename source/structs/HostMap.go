package structs

type HostMap struct {
	Hosts map[string]bool `json:"hosts"`
}

func NewHostMap() HostMap {

	var hostmap HostMap

	hostmap.Hosts = make(map[string]bool)

	return hostmap

}

func (hostmap *HostMap) AddDomain(value string) {
	hostmap.Hosts[value] = true
}

func (hostmap *HostMap) SearchDomain(value string) bool {

	var result bool = false

	_, ok := hostmap.Hosts[value]

	if ok == true {
		result = true
	}

	return result

}
