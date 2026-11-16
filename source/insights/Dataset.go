package insights

import "antispam/structs"
import "antispam/types"
import "antispam/utils/fuzzy"
import "encoding/json"
import "errors"
import "fmt"
import "os"
import "path/filepath"
import "sort"
import "strings"

type JsonFile struct {
	RelPath          string
	AbsPath          string
	Original         []byte
	Entries          []*Entry
	Dirty            bool
	TrailingNewlines int
}

type HostsFile struct {
	RelPath          string
	AbsPath          string
	Original         []byte
	Lines            []string
	Dirty            bool
	TrailingNewlines int
}

type Entry struct {
	Spammer         structs.Spammer
	File            *JsonFile
	Index           int
	Unblockable     bool
	AliasesPresent  bool
	NetworksPresent bool
}

type Dataset struct {
	Root             string
	Entries          []*Entry
	Unblockable      []*Entry
	Files            []*JsonFile
	UnblockableFiles []*JsonFile
	HostsFiles       []*HostsFile

	Spammers            structs.SpammerMap
	UnblockableSpammers structs.SpammerMap
	Hosts               structs.HostMap
}

type FileDiff struct {
	Path string
	Diff string
}

func NewDataset(root string) *Dataset {

	var dataset Dataset

	dataset.Root = root
	dataset.Entries = make([]*Entry, 0)
	dataset.Unblockable = make([]*Entry, 0)
	dataset.Files = make([]*JsonFile, 0)
	dataset.UnblockableFiles = make([]*JsonFile, 0)
	dataset.HostsFiles = make([]*HostsFile, 0)

	dataset.Spammers = structs.NewSpammerMap()
	dataset.UnblockableSpammers = structs.NewSpammerMap()
	dataset.Hosts = structs.NewHostMap()

	return &dataset

}

func Load(root string) (*Dataset, error) {

	absolute, err := filepath.Abs(root)

	if err != nil {
		return nil, err
	}

	dataset := NewDataset(absolute)

	err1 := dataset.loadJsonFolder("phishers", false)

	if err1 != nil {
		return nil, err1
	}

	err2 := dataset.loadJsonFolder("spammers", false)

	if err2 != nil {
		return nil, err2
	}

	err3 := dataset.loadJsonFolder("spammers/unblockable", true)

	if err3 != nil {
		return nil, err3
	}

	err4 := dataset.loadHostsFolder("hosts")

	if err4 != nil {
		return nil, err4
	}

	dataset.Rebuild()

	return dataset, nil

}

func (dataset *Dataset) loadJsonFolder(folder string, unblockable bool) error {

	pattern := filepath.Join(dataset.Root, folder, "*.json")
	files, err := filepath.Glob(pattern)

	if err != nil {
		return err
	}

	sort.Strings(files)

	for _, path := range files {

		buffer, err := os.ReadFile(path)

		if err != nil {
			return err
		}

		var instances []structs.Spammer
		var raw_instances []map[string]json.RawMessage

		err = json.Unmarshal(buffer, &instances)

		if err != nil {
			return errors.New("Cannot parse \"" + path + "\": " + err.Error())
		}

		json.Unmarshal(buffer, &raw_instances)

		name, err := filepath.Rel(dataset.Root, path)

		if err != nil {
			name = path
		}

		file := &JsonFile{
			RelPath:          filepath.ToSlash(filepath.Join("source/insights", name)),
			AbsPath:          path,
			Original:         buffer,
			Entries:          make([]*Entry, 0),
			Dirty:            false,
			TrailingNewlines: countTrailingNewlines(buffer),
		}

		for index, instance := range instances {

			aliases_present := false
			networks_present := false

			if index < len(raw_instances) {

				_, aliases_present = raw_instances[index]["aliases"]
				_, networks_present = raw_instances[index]["networks"]

			}

			entry := &Entry{
				Spammer:         instance,
				File:            file,
				Index:           index,
				Unblockable:     unblockable,
				AliasesPresent:  aliases_present,
				NetworksPresent: networks_present,
			}

			file.Entries = append(file.Entries, entry)

			if unblockable == true {
				dataset.Unblockable = append(dataset.Unblockable, entry)
			} else {
				dataset.Entries = append(dataset.Entries, entry)
			}

		}

		if unblockable == true {
			dataset.UnblockableFiles = append(dataset.UnblockableFiles, file)
		} else {
			dataset.Files = append(dataset.Files, file)
		}

	}

	return nil

}

func (dataset *Dataset) loadHostsFolder(folder string) error {

	pattern := filepath.Join(dataset.Root, folder, "*.hosts")
	files, err := filepath.Glob(pattern)

	if err != nil {
		return err
	}

	sort.Strings(files)

	for _, path := range files {

		buffer, err := os.ReadFile(path)

		if err != nil {
			return err
		}

		name, err := filepath.Rel(dataset.Root, path)

		if err != nil {
			name = path
		}

		lines := make([]string, 0)
		raw_lines := strings.Split(strings.TrimRight(string(buffer), "\n"), "\n")

		for _, line := range raw_lines {

			line = strings.TrimSpace(line)

			if line != "" {
				lines = append(lines, line)
			}

		}

		dataset.HostsFiles = append(dataset.HostsFiles, &HostsFile{
			RelPath:          filepath.ToSlash(filepath.Join("source/insights", name)),
			AbsPath:          path,
			Original:         buffer,
			Lines:            lines,
			Dirty:            false,
			TrailingNewlines: countTrailingNewlines(buffer),
		})

	}

	return nil

}

func (dataset *Dataset) Rebuild() {

	dataset.Spammers = structs.NewSpammerMap()
	dataset.UnblockableSpammers = structs.NewSpammerMap()
	dataset.Hosts = structs.NewHostMap()

	for _, entry := range dataset.Entries {
		dataset.Spammers.AddSpammer(entry.Spammer)
	}

	dataset.Spammers.FlushAliases()

	for _, entry := range dataset.Unblockable {
		dataset.UnblockableSpammers.AddSpammer(entry.Spammer)
	}

	dataset.UnblockableSpammers.FlushAliases()

	for _, file := range dataset.HostsFiles {

		for _, line := range file.Lines {
			dataset.Hosts.AddDomain(line)
		}

	}

}

func (dataset *Dataset) Reload() error {

	fresh, err := Load(dataset.Root)

	if err != nil {
		return err
	}

	dataset.Entries = fresh.Entries
	dataset.Unblockable = fresh.Unblockable
	dataset.Files = fresh.Files
	dataset.UnblockableFiles = fresh.UnblockableFiles
	dataset.HostsFiles = fresh.HostsFiles
	dataset.Rebuild()

	return nil

}

func (dataset *Dataset) AllEntries() []*Entry {

	result := make([]*Entry, 0)
	result = append(result, dataset.Entries...)
	result = append(result, dataset.Unblockable...)

	return result

}

func (entry *Entry) Domain() string {
	return entry.Spammer.Domain
}

func (entry *Entry) FilePath() string {

	if entry.File != nil {
		return entry.File.RelPath
	}

	return ""

}

func (entry *Entry) HasNetwork(value string) bool {

	prefixes := types.ParseNetworks(value)

	if len(prefixes) == 0 {
		return false
	}

	for _, network := range entry.Spammer.Networks {

		if types.NetworkContainsPrefix(network, value) == true {
			return true
		}

	}

	return false

}

func (entry *Entry) AddNetwork(value string) []string {

	result := make([]string, 0)
	prefixes := types.ParseNetworks(value)

	for _, prefix := range prefixes {

		network := prefix.Masked().String()

		if entry.HasNetwork(network) == true {
			continue
		}

		entry.Spammer.Networks = append(entry.Spammer.Networks, network)
		result = append(result, network)

	}

	if len(result) > 0 && entry.File != nil {
		entry.File.Dirty = true
	}

	return result

}

func (entry *Entry) AddAlias(value string) bool {

	value = strings.TrimSpace(strings.ToLower(value))

	if value == "" || value == entry.Spammer.Domain {
		return false
	}

	for _, alias := range entry.Spammer.Aliases {

		if alias == value {
			return false
		}

	}

	entry.Spammer.Aliases = append(entry.Spammer.Aliases, value)

	if entry.File != nil {
		entry.File.Dirty = true
	}

	return true

}

func (dataset *Dataset) AddNetwork(entry *Entry, value string) ([]string, error) {

	if entry == nil {
		return nil, errors.New("no target entry selected")
	}

	if entry.Unblockable == true {
		return nil, errors.New("cannot add networks to unblockable entries")
	}

	result := entry.AddNetwork(value)

	if len(result) > 0 {
		dataset.Rebuild()
	}

	return result, nil

}

func (dataset *Dataset) AddAlias(entry *Entry, value string) (bool, error) {

	if entry == nil {
		return false, errors.New("no target entry selected")
	}

	if entry.Unblockable == true {
		return false, errors.New("cannot add aliases to unblockable entries")
	}

	result := entry.AddAlias(value)

	if result == true {
		dataset.Rebuild()
	}

	return result, nil

}

func (dataset *Dataset) AddHost(file *HostsFile, domain string) bool {

	if file == nil {
		return false
	}

	domain = strings.TrimSpace(strings.ToLower(domain))

	if domain == "" {
		return false
	}

	for _, line := range file.Lines {

		if line == domain {
			return false
		}

	}

	file.Lines = append(file.Lines, domain)
	file.Dirty = true
	dataset.Rebuild()

	return true

}

func (dataset *Dataset) FindEntry(domain string) *Entry {

	for _, entry := range dataset.Entries {

		if entry.Spammer.Domain == domain {
			return entry
		}

	}

	for _, entry := range dataset.Unblockable {

		if entry.Spammer.Domain == domain {
			return entry
		}

	}

	return nil

}

func (dataset *Dataset) FindHostsFile(domain string) *HostsFile {

	for _, file := range dataset.HostsFiles {

		for _, line := range file.Lines {

			if hostMatches(domain, line) == true {
				return file
			}

		}

	}

	return nil

}

func hostMatches(domain string, host string) bool {

	if domain == host {
		return true
	}

	if strings.HasSuffix(domain, "."+host) == true {
		return true
	}

	return false

}

func (dataset *Dataset) MatchValue(value string) *structs.Match {

	if types.IsIPv4(value) == true || types.IsIPv6(value) == true {

		if dataset.Hosts.Hosts[value] == true {
			return dataset.hostsValueMatch(value)
		}

	}

	if types.IsIPv4(value) == true {

		tmp := dataset.Spammers.SearchIPv4(value)

		if tmp != nil {
			return dataset.entryMatch(tmp, false)
		}

		tmp = dataset.UnblockableSpammers.SearchIPv4(value)

		if tmp != nil {
			return dataset.entryMatch(tmp, true)
		}

	} else if types.IsIPv6(value) == true {

		tmp := dataset.Spammers.SearchIPv6(value)

		if tmp != nil {
			return dataset.entryMatch(tmp, false)
		}

		tmp = dataset.UnblockableSpammers.SearchIPv6(value)

		if tmp != nil {
			return dataset.entryMatch(tmp, true)
		}

	} else if types.IsDomain(value) == true {

		if dataset.Hosts.SearchDomain(value) == true {

			host := dataset.hostsMatch(value)
			file := dataset.FindHostsFile(value)

			if file != nil {
				return &structs.Match{
					Domain: host,
					File:   file.RelPath,
					Kind:   "hosts",
				}
			}

			return &structs.Match{
				Domain: host,
				Kind:   "hosts",
			}

		}

		tmp := dataset.Spammers.SearchDomain(value)

		if tmp != nil {
			return dataset.entryMatch(tmp, false)
		}

		tmp = dataset.UnblockableSpammers.SearchDomain(value)

		if tmp != nil {
			return dataset.entryMatch(tmp, true)
		}

	}

	return nil

}

func (dataset *Dataset) hostsValueMatch(value string) *structs.Match {

	file := dataset.FindHostsFile(value)
	match := &structs.Match{
		Domain: value,
		Kind:   "hosts",
	}

	if file != nil {
		match.File = file.RelPath
	}

	return match

}

func (dataset *Dataset) hostsMatch(value string) string {

	if dataset.Hosts.Hosts[value] == true {
		return value
	}

	tmp := strings.Split(value, ".")

	if len(tmp) > 2 {

		parent := strings.Join(tmp[len(tmp)-2:], ".")

		if dataset.Hosts.Hosts[parent] == true {
			return parent
		}

	}

	return value

}

func (dataset *Dataset) entryMatch(spammer *structs.Spammer, unblockable bool) *structs.Match {

	entry := dataset.FindEntry(spammer.Domain)

	kind := "spammer"

	if unblockable == true {
		kind = "unblockable"
	}

	if entry != nil {

		return &structs.Match{
			Domain:      entry.Spammer.Domain,
			File:        entry.FilePath(),
			Kind:        kind,
			Unblockable: unblockable,
		}

	}

	return &structs.Match{
		Domain:      spammer.Domain,
		Kind:        kind,
		Unblockable: unblockable,
	}

}

func (dataset *Dataset) Classify(email *types.Email) (structs.Verdict, *structs.Match, []string) {

	var verdict structs.Verdict = structs.VerdictUnknown
	var match *structs.Match = nil
	reasons := make([]string, 0)

	if email == nil {
		return verdict, match, reasons
	}

	// Specific host entries take precedence over provider networks.
	for _, domain := range email.Domains {

		if dataset.Hosts.SearchDomain(domain) == true {

			host := dataset.hostsMatch(domain)
			file := dataset.FindHostsFile(domain)
			match = &structs.Match{Domain: host, Kind: "hosts"}

			if file != nil {
				match.File = file.RelPath
			}

			reasons = append(reasons, domain)
			verdict = structs.VerdictSpam

			return verdict, match, reasons

		}

	}

	ips := make([]string, 0)
	ips = append(ips, email.SenderIPs...)
	ips = append(ips, email.IPv4s...)
	ips = append(ips, email.IPv6s...)

	for _, value := range ips {

		if dataset.Hosts.Hosts[value] == true {

			match = dataset.hostsValueMatch(value)
			reasons = append(reasons, value)
			verdict = structs.VerdictSpam

			return verdict, match, reasons

		}

	}

	for _, domain := range email.Domains {

		tmp := dataset.Spammers.SearchDomain(domain)

		if tmp != nil {
			match = dataset.entryMatch(tmp, false)
			reasons = append(reasons, domain)
			verdict = structs.VerdictSpam
			return verdict, match, reasons
		}

	}

	for _, value := range email.SenderIPs {

		if types.IsIPv4(value) == true {

			tmp := dataset.Spammers.SearchIPv4(value)

			if tmp != nil {
				match = dataset.entryMatch(tmp, false)
				reasons = append(reasons, value)
				verdict = structs.VerdictSpam
				return verdict, match, reasons
			}

		} else if types.IsIPv6(value) == true {

			tmp := dataset.Spammers.SearchIPv6(value)

			if tmp != nil {
				match = dataset.entryMatch(tmp, false)
				reasons = append(reasons, value)
				verdict = structs.VerdictSpam
				return verdict, match, reasons
			}

		}

	}

	for _, ipv4 := range email.IPv4s {

		tmp := dataset.Spammers.SearchIPv4(ipv4)

		if tmp != nil {
			match = dataset.entryMatch(tmp, false)
			reasons = append(reasons, ipv4)
			verdict = structs.VerdictSpam
			return verdict, match, reasons
		}

	}

	for _, ipv6 := range email.IPv6s {

		tmp := dataset.Spammers.SearchIPv6(ipv6)

		if tmp != nil {
			match = dataset.entryMatch(tmp, false)
			reasons = append(reasons, ipv6)
			verdict = structs.VerdictSpam
			return verdict, match, reasons
		}

	}

	for _, domain := range email.Domains {

		tmp := dataset.UnblockableSpammers.SearchDomain(domain)

		if tmp != nil {
			match = dataset.entryMatch(tmp, true)
			reasons = append(reasons, domain)
			verdict = structs.VerdictUnblockable
			return verdict, match, reasons
		}

	}

	for _, value := range email.SenderIPs {

		if types.IsIPv4(value) == true {

			tmp := dataset.UnblockableSpammers.SearchIPv4(value)

			if tmp != nil {
				match = dataset.entryMatch(tmp, true)
				reasons = append(reasons, value)
				verdict = structs.VerdictUnblockable
				return verdict, match, reasons
			}

		} else if types.IsIPv6(value) == true {

			tmp := dataset.UnblockableSpammers.SearchIPv6(value)

			if tmp != nil {
				match = dataset.entryMatch(tmp, true)
				reasons = append(reasons, value)
				verdict = structs.VerdictUnblockable
				return verdict, match, reasons
			}

		}

	}

	for _, ipv4 := range email.IPv4s {

		tmp := dataset.UnblockableSpammers.SearchIPv4(ipv4)

		if tmp != nil {
			match = dataset.entryMatch(tmp, true)
			reasons = append(reasons, ipv4)
			verdict = structs.VerdictUnblockable
			return verdict, match, reasons
		}

	}

	for _, ipv6 := range email.IPv6s {

		tmp := dataset.UnblockableSpammers.SearchIPv6(ipv6)

		if tmp != nil {
			match = dataset.entryMatch(tmp, true)
			reasons = append(reasons, ipv6)
			verdict = structs.VerdictUnblockable
			return verdict, match, reasons
		}

	}

	return verdict, match, reasons

}

func (dataset *Dataset) FuzzyCandidates(query string) []structs.FuzzyCandidate {

	result := make([]structs.FuzzyCandidate, 0)
	candidates := make([]fuzzy.Candidate, 0)
	lookup := make(map[string]*Entry)

	for _, entry := range dataset.AllEntries() {

		parts := make([]string, 0)
		parts = append(parts, entry.Spammer.Domain)
		parts = append(parts, entry.Spammer.Aliases...)

		if entry.File != nil {
			base := filepath.Base(entry.File.RelPath)
			parts = append(parts, strings.TrimSuffix(base, ".json"))
		}

		key := entry.Spammer.Domain + "|" + entry.FilePath()
		lookup[key] = entry

		candidates = append(candidates, fuzzy.NewCandidate(key, strings.Join(parts, " ")))

	}

	matches := fuzzy.MatchAll(query, candidates, 0.55)

	for _, item := range matches {

		entry := lookup[item.Key]

		if entry == nil {
			continue
		}

		result = append(result, structs.FuzzyCandidate{
			Domain:      entry.Spammer.Domain,
			File:        entry.FilePath(),
			Score:       item.Score,
			Reason:      item.Reason,
			Unblockable: entry.Unblockable,
		})

		if len(result) >= 8 {
			break
		}

	}

	return result

}

func (dataset *Dataset) RouteCandidates(entry *Entry, ip string, routes []string, address_block string) []structs.RouteCandidate {

	result := make([]structs.RouteCandidate, 0)
	seen := make(map[string]bool)

	append_route := func(value string, source string) {

		normalized := types.NormalizeNetworkList(value)

		for _, network := range normalized {

			if seen[network] == true {
				continue
			}

			seen[network] = true

			missing := true

			if entry != nil && entry.HasNetwork(network) == true {
				missing = false
			}

			result = append(result, structs.RouteCandidate{
				Prefix:       network,
				Source:       source,
				ContainsIP:   types.NetworkContains(network, ip),
				Missing:      missing,
				AddressCount: types.NetworkAddressCount(network),
			})

		}

	}

	// The registry allocation is the most relevant route for the observed
	// IP. The observed prefix is always available as a conservative fallback.
	append_route(address_block, "registry")
	append_route(types.ObservedPrefix(ip), "observed")

	for _, route := range routes {
		append_route(route, "route")
	}

	sort.Slice(result, func(a int, b int) bool {

		if result[a].ContainsIP != result[b].ContainsIP {
			return result[a].ContainsIP == true
		}

		if result[a].Missing != result[b].Missing {
			return result[a].Missing == true
		}

		source_a := routeSourceRank(result[a].Source)
		source_b := routeSourceRank(result[b].Source)

		if source_a != source_b {
			return source_a < source_b
		}

		bits_a := networkBits(result[a].Prefix)
		bits_b := networkBits(result[b].Prefix)

		if bits_a != bits_b {
			return bits_a > bits_b
		}

		return result[a].Prefix < result[b].Prefix

	})

	return result

}

func routeSourceRank(value string) int {

	if value == "registry" {
		return 0
	} else if value == "route" {
		return 1
	} else if value == "observed" {
		return 2
	}

	return 3

}

func networkBits(value string) int {

	prefixes := types.ParseNetworks(value)

	if len(prefixes) > 0 {
		return prefixes[0].Bits()
	}

	return 0

}

func (dataset *Dataset) PendingFiles() []string {

	result := make([]string, 0)

	for _, file := range dataset.Files {

		if file.Dirty == true {
			result = append(result, file.RelPath)
		}

	}

	for _, file := range dataset.UnblockableFiles {

		if file.Dirty == true {
			result = append(result, file.RelPath)
		}

	}

	for _, file := range dataset.HostsFiles {

		if file.Dirty == true {
			result = append(result, file.RelPath)
		}

	}

	return result

}

func (dataset *Dataset) Diff() []FileDiff {

	result := make([]FileDiff, 0)

	for _, file := range dataset.Files {

		if file.Dirty == true {

			next := file.Render()
			diff := UnifiedDiff(file.RelPath, file.Original, next)

			if diff != "" {
				result = append(result, FileDiff{Path: file.RelPath, Diff: diff})
			}

		}

	}

	for _, file := range dataset.UnblockableFiles {

		if file.Dirty == true {

			next := file.Render()
			diff := UnifiedDiff(file.RelPath, file.Original, next)

			if diff != "" {
				result = append(result, FileDiff{Path: file.RelPath, Diff: diff})
			}

		}

	}

	for _, file := range dataset.HostsFiles {

		if file.Dirty == true {

			next := file.Render()
			diff := UnifiedDiff(file.RelPath, file.Original, next)

			if diff != "" {
				result = append(result, FileDiff{Path: file.RelPath, Diff: diff})
			}

		}

	}

	return result

}

func (dataset *Dataset) Save() ([]string, error) {

	changed := make([]string, 0)

	for _, file := range dataset.Files {

		if file.Dirty == true {

			next := file.Render()
			err := writeAtomic(file.AbsPath, next)

			if err != nil {
				return changed, err
			}

			file.Original = next
			file.Dirty = false
			changed = append(changed, file.RelPath)

		}

	}

	for _, file := range dataset.UnblockableFiles {

		if file.Dirty == true {

			next := file.Render()
			err := writeAtomic(file.AbsPath, next)

			if err != nil {
				return changed, err
			}

			file.Original = next
			file.Dirty = false
			changed = append(changed, file.RelPath)

		}

	}

	for _, file := range dataset.HostsFiles {

		if file.Dirty == true {

			next := file.Render()
			err := writeAtomic(file.AbsPath, next)

			if err != nil {
				return changed, err
			}

			file.Original = next
			file.Dirty = false
			changed = append(changed, file.RelPath)

		}

	}

	dataset.Rebuild()

	return changed, nil

}

func writeAtomic(path string, buffer []byte) error {

	tmp := path + ".tmp"
	mode := os.FileMode(0644)

	if stat, err := os.Stat(path); err == nil {
		mode = stat.Mode().Perm()
	}

	err := os.WriteFile(tmp, buffer, mode)

	if err != nil {
		return err
	}

	err = os.Rename(tmp, path)

	if err != nil {
		os.Remove(tmp)
		return err
	}

	return nil

}

func (file *JsonFile) Render() []byte {

	instances := make([]structs.Spammer, 0)

	for _, entry := range file.Entries {
		instances = append(instances, entry.Spammer)
	}

	return marshalSpammers(instances, file.Entries, file.TrailingNewlines)

}

func (file *HostsFile) Render() []byte {

	trailing := file.TrailingNewlines

	if trailing < 1 {
		trailing = 0
	}

	return []byte(strings.Join(file.Lines, "\n") + strings.Repeat("\n", trailing))

}

func countTrailingNewlines(buffer []byte) int {

	result := 0

	for index := len(buffer) - 1; index >= 0; index-- {

		if buffer[index] == '\n' {
			result++
		} else {
			break
		}

	}

	return result

}

func renderJsonArray(values []string, present bool) string {

	if len(values) == 0 {

		if present == true {
			return "[\n\t]"
		}

		return ""

	}

	lines := make([]string, 0)

	for _, value := range values {
		lines = append(lines, "\t\t"+quoteJson(value))
	}

	return "[\n" + strings.Join(lines, ",\n") + "\n\t]"

}

func marshalSpammers(instances []structs.Spammer, entries []*Entry, trailing int) []byte {

	builder := strings.Builder{}
	builder.WriteString("[")

	for index, instance := range instances {

		if index > 0 {
			builder.WriteString(", ")
		}

		entry := &Entry{}

		if index < len(entries) {
			entry = entries[index]
		}

		fields := make([]string, 0)

		fields = append(fields, "\t\"domain\": "+quoteJson(instance.Domain))

		aliases := renderJsonArray(instance.Aliases, entry.AliasesPresent)

		if aliases != "" {
			fields = append(fields, "\t\"aliases\": "+aliases)
		}

		networks := renderJsonArray(instance.Networks, entry.NetworksPresent)

		if networks != "" {
			fields = append(fields, "\t\"networks\": "+networks)
		}

		builder.WriteString("{\n")
		builder.WriteString(strings.Join(fields, ",\n"))
		builder.WriteString("\n}")

	}

	builder.WriteString("]")

	if trailing < 1 {
		trailing = 1
	}

	builder.WriteString(strings.Repeat("\n", trailing))

	return []byte(builder.String())

}

func quoteJson(value string) string {

	buffer, err := json.Marshal(value)

	if err == nil {
		return string(buffer)
	}

	return "\"\""

}

func (dataset *Dataset) String() string {
	return fmt.Sprintf("Dataset{root=%s entries=%d unblockable=%d hosts=%d}", dataset.Root, len(dataset.Entries), len(dataset.Unblockable), len(dataset.HostsFiles))
}
