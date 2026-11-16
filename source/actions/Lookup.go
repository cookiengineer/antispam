package actions

import "antispam/protocols/dns"
import "antispam/protocols/whois"
import "antispam/structs"
import "antispam/types"
import "sync"
import "time"

type lookupCacheEntry struct {
	Value     structs.IPLookup
	ExpiresAt time.Time
}

type LookupCache struct {
	mutex   sync.Mutex
	entries map[string]lookupCacheEntry
}

func NewLookupCache() *LookupCache {

	return &LookupCache{
		entries: make(map[string]lookupCacheEntry),
	}

}

func (cache *LookupCache) Get(ip string) (structs.IPLookup, bool) {

	if cache == nil {
		return structs.IPLookup{}, false
	}

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	entry, ok := cache.entries[ip]

	if ok == false {
		return structs.IPLookup{}, false
	}

	if time.Now().After(entry.ExpiresAt) == true {
		delete(cache.entries, ip)
		return structs.IPLookup{}, false
	}

	value := entry.Value
	value.Cached = true

	return value, true

}

func (cache *LookupCache) Set(ip string, value structs.IPLookup) {

	if cache == nil {
		return
	}

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.entries[ip] = lookupCacheEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

}

func LookupIP(ip string, cache *LookupCache) structs.IPLookup {

	ip = types.CanonicalIP(ip)

	if cached, ok := cache.Get(ip); ok == true {
		return cached
	}

	result := structs.IPLookup{
		IP:               ip,
		PTR:              make([]string, 0),
		ForwardConfirmed: make([]string, 0),
		Routes:           make([]string, 0),
		Errors:           make([]string, 0),
		RetrievedAt:      time.Now(),
	}

	if types.IsPrivateIP(ip) == true {
		result.Errors = append(result.Errors, "IP address is not publicly routable")
		return result
	}

	if types.IsIPv4(ip) == false && types.IsIPv6(ip) == false {
		result.Errors = append(result.Errors, "Not an IPv4 or IPv6 address")
		return result
	}

	resolve_ptr(ip, &result)

	if types.IsIPv4(ip) == true {

		info := whois.ResolveIPv4(ip, nil)

		result.AddressBlock = info.AddressBlock
		result.ASN = info.Origin
		result.Organization = info.OrganizationName
		result.NetworkName = info.NetworkName
		result.Registry = info.Source
		result.Country = info.Country

	} else if types.IsIPv6(ip) == true {

		info := whois.ResolveIPv6(ip, nil)

		result.AddressBlock = info.AddressBlock
		result.ASN = info.Origin
		result.Organization = info.OrganizationName
		result.NetworkName = info.NetworkName
		result.Registry = info.Source
		result.Country = info.Country

	}

	if result.ASN == "" {
		result.ASN = whois.ResolveASN(ip)
	}

	if result.ASN != "" {

		routes := whois.ResolveRoutes(result.ASN, nil)

		for _, route := range routes {

			for _, network := range types.NormalizeNetworkList(route.AddressBlock) {
				appendLookupRoute(&result, network)
			}

		}

	}

	for _, network := range types.NormalizeNetworkList(result.AddressBlock) {
		appendLookupRoute(&result, network)
	}

	if len(result.PTR) == 0 && len(result.Routes) == 0 && result.ASN == "" && result.Organization == "" {
		result.Errors = append(result.Errors, "No WHOIS or DNS data found")
	}

	if len(result.Errors) == 0 || len(result.Routes) > 0 || result.ASN != "" || result.Organization != "" {
		cache.Set(ip, result)
	}

	return result

}

func resolve_ptr(ip string, result *structs.IPLookup) {

	response, err := dns.Resolve(ip)

	if err != nil {
		result.Errors = append(result.Errors, "Reverse DNS failed: "+err.Error())
		return
	}

	if response.Type != "response" {
		result.Errors = append(result.Errors, "Reverse DNS returned no response")
		return
	}

	for _, answer := range response.Answers {

		if answer.Type == dns.TypePTR {

			ptr := answer.ToDomain()

			if ptr != "" {
				appendUnique(&result.PTR, ptr)
			}

		}

	}

	for _, ptr := range result.PTR {

		forward, err := dns.Resolve(ptr)

		if err != nil || forward.Type != "response" {
			continue
		}

		for _, answer := range forward.Answers {

			if answer.Type == dns.TypeA {

				value := types.CanonicalIP(answer.ToIPv4())

				if value == ip {
					appendUnique(&result.ForwardConfirmed, ptr)
				}

			} else if answer.Type == dns.TypeAAAA {

				value := types.CanonicalIP(answer.ToIPv6())

				if value == ip {
					appendUnique(&result.ForwardConfirmed, ptr)
				}

			}

		}

	}

}

func appendLookupRoute(result *structs.IPLookup, network string) {

	if network == "" {
		return
	}

	for _, item := range result.Routes {

		if item == network {
			return
		}

	}

	result.Routes = append(result.Routes, network)

}

func appendUnique(values *[]string, value string) {

	if value == "" {
		return
	}

	for _, item := range *values {

		if item == value {
			return
		}

	}

	*values = append(*values, value)

}
