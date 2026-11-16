package actions

import "antispam/insights"
import "antispam/structs"
import "antispam/types"
import "errors"
import "os"
import "sort"

func Analyze(file string, dataset *insights.Dataset) (*structs.Analysis, error) {

	var result *structs.Analysis = nil

	if dataset == nil {
		return nil, errors.New("no dataset loaded")
	}

	stat, err0 := os.Stat(file)

	if err0 != nil || stat.IsDir() {
		return nil, errors.New("cannot read file \"" + file + "\"")
	}

	buffer, err1 := os.ReadFile(file)

	if err1 != nil || types.IsEmail(buffer) == false {
		return nil, errors.New("cannot parse file \"" + file + "\"")
	}

	email := types.ParseEmail(buffer)

	if email == nil {
		return nil, errors.New("cannot parse file \"" + file + "\"")
	}

	verdict, match, reasons := dataset.Classify(email)

	analysis := &structs.Analysis{
		File:       file,
		Email:      email,
		Verdict:    verdict,
		Match:      match,
		Reasons:    reasons,
		Indicators: make([]*structs.Indicator, 0),
	}

	seen := make(map[string]bool)

	append_indicator := func(kind string, value string, authoritative bool) {

		key := kind + ":" + value

		if seen[key] == true {
			return
		}

		seen[key] = true

		indicator := &structs.Indicator{
			Kind:             kind,
			Value:            value,
			Private:          false,
			Authoritative:    authoritative,
			Candidates:       make([]structs.FuzzyCandidate, 0),
			PTR:              make([]string, 0),
			ForwardConfirmed: make([]string, 0),
			Routes:           make([]string, 0),
		}

		if kind == "ipv4" || kind == "ipv6" {
			indicator.Private = types.IsPrivateIP(value)
		}

		known := dataset.MatchValue(value)

		if known != nil {

			if known.Unblockable == true {
				indicator.Unblockable = known
			} else {
				indicator.Known = known
			}

		}

		if kind == "domain" {
			indicator.Candidates = dataset.FuzzyCandidates(value)
		}

		analysis.Indicators = append(analysis.Indicators, indicator)

	}

	for _, value := range email.SenderIPs {

		kind := "ipv4"

		if types.IsIPv6(value) == true {
			kind = "ipv6"
		}

		append_indicator(kind, value, true)

	}

	for _, value := range email.IPv4s {
		append_indicator("ipv4", value, false)
	}

	for _, value := range email.IPv6s {
		append_indicator("ipv6", value, false)
	}

	for _, value := range email.Domains {
		append_indicator("domain", value, false)
	}

	if len(email.SenderIPs) == 0 {

		for _, indicator := range analysis.Indicators {

			if (indicator.Kind == "ipv4" || indicator.Kind == "ipv6") && indicator.Private == false {
				indicator.Authoritative = true
				break
			}

		}

	}

	sort.SliceStable(analysis.Indicators, func(a int, b int) bool {

		left := analysis.Indicators[a]
		right := analysis.Indicators[b]

		if left.Authoritative != right.Authoritative {
			return left.Authoritative == true
		}

		left_ip := left.Kind == "ipv4" || left.Kind == "ipv6"
		right_ip := right.Kind == "ipv4" || right.Kind == "ipv6"

		if left_ip != right_ip {
			return left_ip == true
		}

		if left.Private != right.Private {
			return left.Private == false
		}

		return false

	})

	result = analysis

	return result, nil

}
