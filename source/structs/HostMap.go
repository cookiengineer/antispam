package structs

import "strings"

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

	_, ok1 := hostmap.Hosts[value]

	if ok1 == true {
		result = true
	}

	if result == false {

		if strings.Contains(value, ".") {

			tmp := strings.Split(value, ".")

			if len(tmp) > 2 {
				tmp = tmp[len(tmp)-2:]
			}

			domain := strings.Join(tmp, ".")
			_, ok2 := hostmap.Hosts[domain]

			if ok2 == true {
				result = true
			}

		}

	}

	return result

}
