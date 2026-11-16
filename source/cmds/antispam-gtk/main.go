package main

import "antispam/actions"
import "antispam/bindings/gtk"
import "antispam/insights"
import "antispam/structs"
import "fmt"
import "os"
import "path/filepath"
import "sort"
import "strings"
import "sync/atomic"

type Helper struct {
	app *gtk.Application

	insightsRoot string
	mailsRoot    string

	dataset      *insights.Dataset
	datasetError string
	cache        *actions.LookupCache

	analyses          []*structs.Analysis
	visible           []int
	selected          int
	selectedIndicator int

	targets            []*insights.Entry
	targetEntry        *insights.Entry
	hosts              []*insights.HostsFile
	domainValues       []string
	routeCandidates    []structs.RouteCandidate
	currentUnblockable bool

	running     atomic.Bool
	cancel      atomic.Bool
	loadCounter atomic.Int64
	loading     bool
	updating    bool
	filterValue int

	window *gtk.Window

	folderLabel *gtk.Label
	statusLabel *gtk.Label
	progress    *gtk.ProgressBar

	emailList  *gtk.ListBox
	summary    *gtk.Label
	indicators *gtk.ListBox
	message    *gtk.TextView

	targetBox  *gtk.Box
	targetDrop *gtk.DropDown

	prefixEntry *gtk.Entry

	routeBox  *gtk.Box
	routeDrop *gtk.DropDown

	hostsBox  *gtk.Box
	hostsDrop *gtk.DropDown

	domainBox  *gtk.Box
	domainDrop *gtk.DropDown

	pendingLabel  *gtk.Label
	analyzeButton *gtk.Button
}

func main() {

	arguments := parseArgs()

	root := findRepositoryRoot()

	if root == "" {

		cwd, err := os.Getwd()

		if err == nil {
			root = cwd
		}

	}

	insights_root := filepath.Join(root, "source", "insights")
	mails_root := filepath.Join(root, "mails")

	if arguments.insights != "" {
		insights_root = arguments.insights
	}

	if arguments.mails != "" {
		mails_root = arguments.mails
	}

	helper := &Helper{
		app:               gtk.NewApplication("network.tholian.AntispamHelper"),
		insightsRoot:      insights_root,
		mailsRoot:         mails_root,
		cache:             actions.NewLookupCache(),
		selected:          -1,
		selectedIndicator: -1,
		filterValue:       -1,
	}

	helper.loadDataset()

	helper.app.OnActivate(func() {
		helper.activate()
	})

	os.Exit(helper.app.Run())

}

type arguments struct {
	insights string
	mails    string
}

func parseArgs() arguments {

	result := arguments{}

	for _, value := range os.Args[1:] {

		if strings.HasPrefix(value, "--insights=") {
			result.insights = strings.TrimSpace(value[11:])
		} else if strings.HasPrefix(value, "--mails=") {
			result.mails = strings.TrimSpace(value[8:])
		}

	}

	return result

}

func findRepositoryRoot() string {

	starts := make([]string, 0)

	cwd, err := os.Getwd()

	if err == nil {
		starts = append(starts, cwd)
	}

	executable, err := os.Executable()

	if err == nil {
		starts = append(starts, filepath.Dir(executable))
	}

	for _, start := range starts {

		current := start

		for {

			stat, err := os.Stat(filepath.Join(current, "source", "insights"))

			if err == nil && stat.IsDir() == true {
				return current
			}

			parent := filepath.Dir(current)

			if parent == current {
				break
			}

			current = parent

		}

	}

	return ""

}

func (helper *Helper) activate() {

	helper.window = gtk.NewWindow(helper.app)
	helper.window.SetTitle("Antispam Helper")
	helper.window.SetDefaultSize(1280, 860)

	helper.buildUI()
	helper.window.Present()

	if helper.datasetError != "" {
		helper.status(helper.datasetError)
	}

	helper.loadEmails()

}

func (helper *Helper) buildUI() {

	root := gtk.NewBox(gtk.OrientationVertical, 0)
	root.SetMarginStart(6)
	root.SetMarginEnd(6)
	root.SetMarginTop(6)
	root.SetMarginBottom(6)

	toolbar := gtk.NewBox(gtk.OrientationHorizontal, 8)

	folder_button := gtk.NewButton("Choose mails folder")

	folder_button.OnClick(func() {

		gtk.ShowFolderDialog(helper.window.AsPtr(), func(path string) {

			if path != "" {
				helper.mailsRoot = path
				helper.loadEmails()
			}

		})

	})

	helper.folderLabel = gtk.NewLabel(validUTF8(helper.mailsRoot))
	helper.folderLabel.SetSelectable(true)
	helper.folderLabel.SetHExpand(true)
	helper.folderLabel.SetHAlign(gtk.AlignStart)

	reload_button := gtk.NewButton("Reload")

	reload_button.OnClick(func() {
		helper.reloadAll()
	})

	filter := gtk.NewDropDown([]string{"All", "Unknown", "Spam", "Unblockable", "Has match"})

	filter.OnChanged(func() {
		helper.filterValue = int(filter.GetSelected())
		helper.refreshEmailList()
	})

	toolbar.Append(folder_button.AsPtr())
	toolbar.Append(helper.folderLabel.AsPtr())
	toolbar.Append(filter.AsPtr())
	toolbar.Append(reload_button.AsPtr())

	content := gtk.NewBox(gtk.OrientationHorizontal, 6)
	content.SetVExpand(true)

	helper.emailList = gtk.NewListBox()

	helper.emailList.OnRowActivatedIndex(func(index int) {
		helper.selectVisible(index)
	})

	left_scroll := gtk.NewScrolledWindow()
	left_scroll.SetPolicy(gtk.PolicyAutomatic, gtk.PolicyAutomatic)
	left_scroll.SetVExpand(true)
	left_scroll.SetSizeRequest(360, -1)
	left_scroll.SetChild(helper.emailList.AsPtr())

	right := gtk.NewBox(gtk.OrientationVertical, 8)
	right.SetHExpand(true)
	right.SetVExpand(true)

	helper.summary = gtk.NewLabel("Select an email")
	helper.summary.SetWrap(true)
	helper.summary.SetXAlign(0)
	helper.summary.SetSelectable(true)

	helper.indicators = gtk.NewListBox()

	helper.indicators.OnRowActivatedIndex(func(index int) {
		helper.selectIndicator(index)
	})

	indicator_scroll := gtk.NewScrolledWindow()
	indicator_scroll.SetPolicy(gtk.PolicyAutomatic, gtk.PolicyAutomatic)
	indicator_scroll.SetSizeRequest(-1, 220)
	indicator_scroll.SetChild(helper.indicators.AsPtr())

	learning := gtk.NewGrid()
	learning.SetColumnSpacing(6)
	learning.SetRowSpacing(6)
	learning.SetHExpand(true)

	helper.targetBox = gtk.NewBox(gtk.OrientationHorizontal, 0)
	helper.targetBox.SetHExpand(true)

	helper.prefixEntry = gtk.NewEntry()
	helper.prefixEntry.SetHExpand(true)
	helper.prefixEntry.SetPlaceholder("CIDR prefix, e.g. 203.0.113.0/24")

	prefix_row := gtk.NewBox(gtk.OrientationHorizontal, 8)
	prefix_row.SetHExpand(true)

	add_prefix := gtk.NewButton("Add prefix")

	add_prefix.OnClick(func() {
		helper.addPrefix()
	})

	prefix_row.Append(helper.prefixEntry.AsPtr())
	prefix_row.Append(add_prefix.AsPtr())

	helper.routeBox = gtk.NewBox(gtk.OrientationHorizontal, 0)
	helper.routeBox.SetHExpand(true)

	helper.hostsBox = gtk.NewBox(gtk.OrientationHorizontal, 0)
	helper.domainBox = gtk.NewBox(gtk.OrientationHorizontal, 0)

	hosts_row := gtk.NewBox(gtk.OrientationHorizontal, 8)
	hosts_row.SetHExpand(true)

	add_host := gtk.NewButton("Add host")

	add_host.OnClick(func() {
		helper.addHost()
	})

	hosts_row.Append(helper.hostsBox.AsPtr())
	hosts_row.Append(helper.domainBox.AsPtr())
	hosts_row.Append(add_host.AsPtr())

	helper.pendingLabel = gtk.NewLabel("No pending changes")
	helper.pendingLabel.SetHExpand(true)
	helper.pendingLabel.SetXAlign(0)

	review_button := gtk.NewButton("Review & Save")
	review_button.OnClick(func() {
		helper.showDiff()
	})

	pending_row := gtk.NewBox(gtk.OrientationHorizontal, 8)
	pending_row.Append(helper.pendingLabel.AsPtr())
	pending_row.Append(review_button.AsPtr())

	label_target := gtk.NewLabel("Target:")
	label_target.SetXAlign(0)
	label_prefix := gtk.NewLabel("Prefix:")
	label_prefix.SetXAlign(0)
	label_route := gtk.NewLabel("WHOIS route:")
	label_route.SetXAlign(0)
	label_hosts := gtk.NewLabel("Hosts fallback:")
	label_hosts.SetXAlign(0)
	label_pending := gtk.NewLabel("Dataset:")
	label_pending.SetXAlign(0)

	learning.Attach(label_target.AsPtr(), 0, 0, 1, 1)
	learning.Attach(helper.targetBox.AsPtr(), 1, 0, 1, 1)
	learning.Attach(label_prefix.AsPtr(), 0, 1, 1, 1)
	learning.Attach(prefix_row.AsPtr(), 1, 1, 1, 1)
	learning.Attach(label_route.AsPtr(), 0, 2, 1, 1)
	learning.Attach(helper.routeBox.AsPtr(), 1, 2, 1, 1)
	learning.Attach(label_hosts.AsPtr(), 0, 3, 1, 1)
	learning.Attach(hosts_row.AsPtr(), 1, 3, 1, 1)
	learning.Attach(label_pending.AsPtr(), 0, 4, 1, 1)
	learning.Attach(pending_row.AsPtr(), 1, 4, 1, 1)

	helper.message = gtk.NewTextView()
	helper.message.SetEditable(false)
	helper.message.SetMonospace(true)
	helper.message.SetCursorVisible(false)
	helper.message.SetWrapModeWordChar()

	message_scroll := gtk.NewScrolledWindow()
	message_scroll.SetPolicy(gtk.PolicyAutomatic, gtk.PolicyAutomatic)
	message_scroll.SetVExpand(true)
	message_scroll.SetChild(helper.message.AsPtr())

	right.Append(helper.summary.AsPtr())
	right.Append(indicator_scroll.AsPtr())
	right.Append(learning.AsPtr())
	right.Append(message_scroll.AsPtr())

	content.Append(left_scroll.AsPtr())
	content.Append(right.AsPtr())

	action_bar := gtk.NewBox(gtk.OrientationHorizontal, 8)

	helper.analyzeButton = gtk.NewButton("Analyze IPs")

	helper.analyzeButton.OnClick(func() {
		helper.startLookup()
	})

	spam_button := gtk.NewButton("Archive as Spam")

	spam_button.OnClick(func() {
		helper.confirmArchive(true)
	})

	not_spam_button := gtk.NewButton("Archive as Not Spam")

	not_spam_button.OnClick(func() {
		helper.confirmArchive(false)
	})

	action_bar.Append(helper.analyzeButton.AsPtr())
	action_bar.Append(spam_button.AsPtr())
	action_bar.Append(not_spam_button.AsPtr())

	status := gtk.NewBox(gtk.OrientationHorizontal, 8)

	helper.progress = gtk.NewProgressBar()
	helper.progress.SetShowText(false)
	helper.progress.SetSizeRequest(220, -1)

	helper.statusLabel = gtk.NewLabel("Ready")
	helper.statusLabel.SetXAlign(0)
	helper.statusLabel.SetHExpand(true)
	helper.statusLabel.SetHAlign(gtk.AlignStart)
	helper.statusLabel.SetSelectable(true)

	status.Append(helper.progress.AsPtr())
	status.Append(helper.statusLabel.AsPtr())

	root.Append(toolbar.AsPtr())
	root.Append(content.AsPtr())
	root.Append(action_bar.AsPtr())
	root.Append(status.AsPtr())

	helper.window.SetChild(root.AsPtr())

	helper.rebuildTargets()
	helper.rebuildHosts()
	helper.rebuildDomains(nil)
	helper.updatePending()

}

func (helper *Helper) loadDataset() {

	dataset, err := insights.Load(helper.insightsRoot)

	if err != nil {
		helper.datasetError = "Cannot load insights: " + err.Error()
		helper.status(helper.datasetError)
		return
	}

	helper.dataset = dataset
	helper.datasetError = ""

}

func (helper *Helper) loadEmails() {

	if helper.window == nil {
		return
	}

	if helper.dataset == nil {
		helper.status("No insights dataset loaded")
		return
	}

	counter := helper.loadCounter.Add(1)
	mails_root := helper.mailsRoot
	dataset := helper.dataset

	if helper.running.Load() == true {
		helper.cancel.Store(true)
	}

	helper.loading = true
	helper.progress.SetFraction(0)
	helper.status("Loading emails from " + mails_root)

	go func() {

		entries, err := os.ReadDir(mails_root)

		if err != nil {

			gtk.RunOnMain(func() {

				if helper.loadCounter.Load() != counter {
					return
				}

				helper.loading = false
				helper.status("Cannot read mails folder: " + err.Error())

			})

			return

		}

		files := make([]string, 0)

		for _, entry := range entries {

			if entry.IsDir() == true {
				continue
			}

			name := entry.Name()

			if strings.HasPrefix(name, ".") == true {
				continue
			}

			files = append(files, filepath.Join(mails_root, name))

		}

		sort.Strings(files)

		analyses := make([]*structs.Analysis, 0)
		failures := 0

		for index, file := range files {

			if helper.loadCounter.Load() != counter {
				return
			}

			analysis, err := actions.Analyze(file, dataset)

			if err != nil {
				failures++
				continue
			}

			analyses = append(analyses, analysis)

			if index%5 == 0 {

				fraction := float64(index+1) / float64(len(files))
				message := fmt.Sprintf("Loading emails %d/%d...", index+1, len(files))

				gtk.RunOnMain(func() {

					if helper.loadCounter.Load() != counter {
						return
					}

					helper.progress.SetFraction(fraction)
					helper.status(message)

				})

			}

		}

		sort.SliceStable(analyses, func(a int, b int) bool {

			left := analyses[a]
			right := analyses[b]

			if left.Email.Date.Equal(right.Email.Date) == false {
				return left.Email.Date.After(right.Email.Date)
			}

			return left.File > right.File

		})

		gtk.RunOnMain(func() {

			if helper.loadCounter.Load() != counter {
				return
			}

			helper.loading = false
			helper.progress.SetFraction(1)
			helper.analyses = analyses
			helper.selected = -1
			helper.selectedIndicator = -1

			helper.refreshEmailList()

			if helper.selected == -1 && len(helper.visible) > 0 {
				helper.selected = helper.visible[0]
			}

			helper.updateDetail()
			helper.folderLabel.SetText(validUTF8(mails_root))
			helper.status(fmt.Sprintf("Loaded %d emails (%d unparsable)", len(helper.analyses), failures))

		})

	}()

}

func (helper *Helper) reloadAll() {

	if helper.loading == true {
		helper.status("Restarting email load...")
		helper.loadEmails()
		return
	}

	if helper.dataset != nil && len(helper.dataset.PendingFiles()) > 0 {

		helper.showConfirm(
			"Discard dataset changes?",
			"There are unsaved dataset changes. Reloading will discard them.",
			"Discard & Reload",
			func() {
				helper.reloadNow()
			},
		)

		return

	}

	helper.reloadNow()

}

func (helper *Helper) reloadNow() {

	helper.loadDataset()

	if helper.dataset == nil {
		return
	}

	helper.rebuildTargets()
	helper.rebuildHosts()
	helper.loadEmails()

}

func (helper *Helper) status(message string) {

	if helper.statusLabel != nil {
		helper.statusLabel.SetText(validUTF8(message))
	}

}

func (helper *Helper) refreshEmailList() {

	if helper.emailList == nil {
		return
	}

	helper.emailList.Clear()
	helper.visible = make([]int, 0)

	filter := 0

	// The filter drop-down is created locally, so use the stored index from
	// its change callback. Default to "All" when no callback has fired yet.
	if helper.filterValue >= 0 {
		filter = helper.filterValue
	}

	selected_row := -1

	for index, analysis := range helper.analyses {

		if matchesFilter(analysis, filter) == false {
			continue
		}

		helper.visible = append(helper.visible, index)
		appendWrappedRow(helper.emailList, emailRowText(analysis))

		if index == helper.selected {
			selected_row = len(helper.visible) - 1
		}

	}

	if selected_row >= 0 {
		helper.emailList.SelectRow(selected_row)
	}

	if helper.selected >= 0 {

		found := false

		for _, index := range helper.visible {

			if index == helper.selected {
				found = true
				break
			}

		}

		if found == false {
			helper.selected = -1
		}

	}

	if helper.selected == -1 {

		if len(helper.visible) > 0 {

			helper.selected = helper.visible[0]
			helper.emailList.SelectRow(0)

		}

		helper.updateDetail()

	}

}

func matchesFilter(analysis *structs.Analysis, filter int) bool {

	if filter == 0 {
		return true
	} else if filter == 1 {
		return analysis.Verdict == structs.VerdictUnknown
	} else if filter == 2 {
		return analysis.Verdict == structs.VerdictSpam
	} else if filter == 3 {
		return analysis.Verdict == structs.VerdictUnblockable
	} else if filter == 4 {
		return analysis.Match != nil
	}

	return true

}

func emailRowText(analysis *structs.Analysis) string {

	badge := "[UNKNOWN]"

	if analysis.Verdict == structs.VerdictSpam {
		badge = "[SPAM]"
	} else if analysis.Verdict == structs.VerdictUnblockable {
		badge = "[UNBLOCKABLE]"
	} else if analysis.Verdict == structs.VerdictNotSpam {
		badge = "[NOT SPAM]"
	}

	subject := ""
	sender := ""

	if analysis.Email != nil {
		subject = analysis.Email.Subject
		sender = analysis.Email.From
	}

	if subject == "" {
		subject = "(no subject)"
	}

	if sender == "" {
		sender = filepath.Base(analysis.File)
	}

	row := badge + " " + subject + "\n" + sender

	if analysis.Match != nil {
		row += " — " + analysis.Match.Domain
	}

	return row

}

func (helper *Helper) selectVisible(index int) {

	if index < 0 || index >= len(helper.visible) {
		return
	}

	helper.selected = helper.visible[index]
	helper.selectedIndicator = -1
	helper.updateDetail()

}

func (helper *Helper) currentAnalysis() *structs.Analysis {

	if helper.selected < 0 || helper.selected >= len(helper.analyses) {
		return nil
	}

	return helper.analyses[helper.selected]

}

func (helper *Helper) currentIndicator() *structs.Indicator {

	analysis := helper.currentAnalysis()

	if analysis == nil || helper.selectedIndicator < 0 || helper.selectedIndicator >= len(analysis.Indicators) {
		return nil
	}

	return analysis.Indicators[helper.selectedIndicator]

}

func (helper *Helper) syncSelection() {

	if helper.emailList == nil {
		return
	}

	for row, index := range helper.visible {

		if index == helper.selected {
			helper.emailList.SelectRow(row)
			return
		}

	}

}

func (helper *Helper) updateDetail() {

	helper.syncSelection()
	helper.resetTarget()

	analysis := helper.currentAnalysis()

	if analysis == nil {

		helper.summary.SetText("No email selected")
		helper.indicators.Clear()
		helper.message.SetText("")
		helper.rebuildDomains(nil)
		helper.refreshRoutes()
		helper.updatePending()

		return

	}

	email := analysis.Email
	verdict := verdictText(analysis)

	summary := "<b>File:</b> " + escapeMarkup(filepath.Base(analysis.File)) + "\n"

	if email != nil {

		summary += "<b>From:</b> " + escapeMarkup(email.From) + "\n"
		summary += "<b>To:</b> " + escapeMarkup(email.To) + "\n"
		summary += "<b>Subject:</b> " + escapeMarkup(email.Subject) + "\n"
		summary += "<b>Date:</b> " + escapeMarkup(email.Date.Format("2006-01-02 15:04:05")) + "\n"
		summary += "<b>Domains:</b> " + escapeMarkup(strings.Join(email.Domains, ", ")) + "\n"
		summary += "<b>Sender IPs:</b> " + escapeMarkup(strings.Join(email.SenderIPs, ", ")) + "\n"

	}

	summary += "<b>Verdict:</b> " + verdict

	helper.summary.SetMarkup(summary)

	if email != nil {
		helper.message.SetText(validUTF8(email.Message))
	} else {
		helper.message.SetText("")
	}

	helper.indicators.Clear()

	for _, indicator := range analysis.Indicators {
		appendWrappedRow(helper.indicators, indicatorText(indicator))
	}

	helper.rebuildDomains(analysis)
	helper.selectedIndicator = -1

	if len(analysis.Indicators) > 0 {
		helper.selectedIndicator = 0
		helper.indicators.SelectRow(0)
	}

	helper.selectTargetForCurrentIndicator()
	helper.refreshRoutes()
	helper.updatePending()

}

func verdictText(analysis *structs.Analysis) string {

	if analysis.Verdict == structs.VerdictSpam {

		result := "SPAM"

		if analysis.Match != nil {
			result += " via " + analysis.Match.Domain

			if analysis.Match.File != "" {
				result += " (" + analysis.Match.File + ")"
			}

		}

		if len(analysis.Reasons) > 0 {
			result += " because of " + strings.Join(analysis.Reasons, ", ")
		}

		return escapeMarkup(result)

	} else if analysis.Verdict == structs.VerdictUnblockable {

		result := "UNBLOCKABLE (treat as Not Spam)"

		if analysis.Match != nil {
			result += " via " + analysis.Match.Domain
		}

		return escapeMarkup(result)

	} else if analysis.Verdict == structs.VerdictNotSpam {

		return "NOT SPAM"

	}

	return "UNKNOWN"

}

func indicatorText(indicator *structs.Indicator) string {

	parts := make([]string, 0)
	parts = append(parts, "["+strings.ToUpper(indicator.Kind)+"]")
	parts = append(parts, indicator.Value)

	if indicator.Private == true {
		parts = append(parts, "private")
	}

	if indicator.Authoritative == true {
		parts = append(parts, "sender")
	}

	if indicator.Known != nil {
		parts = append(parts, "match="+indicator.Known.Domain)
	}

	if indicator.Unblockable != nil {
		parts = append(parts, "unblockable="+indicator.Unblockable.Domain)
	}

	if indicator.ASN != "" {
		parts = append(parts, "asn="+indicator.ASN)
	}

	if indicator.Organization != "" {
		parts = append(parts, "org="+indicator.Organization)
	}

	if indicator.NetworkName != "" {
		parts = append(parts, "netname="+indicator.NetworkName)
	}

	if indicator.AddressBlock != "" {
		parts = append(parts, "block="+indicator.AddressBlock)
	}

	if indicator.Registry != "" {
		parts = append(parts, "registry="+indicator.Registry)
	}

	if indicator.Country != "" {
		parts = append(parts, "country="+indicator.Country)
	}

	if len(indicator.PTR) > 0 {
		parts = append(parts, "ptr="+strings.Join(indicator.PTR, ","))
	}

	if len(indicator.ForwardConfirmed) > 0 {
		parts = append(parts, "fcrdns="+strings.Join(indicator.ForwardConfirmed, ","))
	}

	if indicator.Error != "" {
		parts = append(parts, "error="+indicator.Error)
	}

	return truncate(validUTF8(strings.Join(parts, "  ")), 360)

}

func appendWrappedRow(list *gtk.ListBox, text string) {

	row := gtk.NewListBoxRow()
	row.SetTextWrap(text)
	list.AppendRow(row)

}

func (helper *Helper) selectIndicator(index int) {

	helper.selectedIndicator = index
	helper.selectTargetForCurrentIndicator()
	helper.refreshRoutes()

}

func (helper *Helper) rebuildTargets() {

	if helper.targetBox == nil {
		return
	}

	helper.targets = make([]*insights.Entry, 0)
	helper.targetEntry = nil
	labels := []string{"(select target)"}

	if helper.dataset != nil {

		entries := make([]*insights.Entry, 0)
		entries = append(entries, helper.dataset.Entries...)

		sort.SliceStable(entries, func(a int, b int) bool {

			if entries[a].Spammer.Domain == entries[b].Spammer.Domain {
				return entries[a].FilePath() < entries[b].FilePath()
			}

			return entries[a].Spammer.Domain < entries[b].Spammer.Domain

		})

		for _, entry := range entries {
			helper.targets = append(helper.targets, entry)
			labels = append(labels, entry.Spammer.Domain+" — "+filepath.Base(entry.FilePath()))
		}

	}

	drop := gtk.NewDropDown(labels)
	drop.SetSizeRequest(340, -1)

	drop.OnChanged(func() {

		if helper.updating == true {
			return
		}

		index := int(drop.GetSelected())

		if index == 0 {
			helper.targetEntry = nil
			helper.currentUnblockable = false
			helper.refreshRoutes()
		} else if index > 0 && (index-1) < len(helper.targets) {
			helper.targetEntry = helper.targets[index-1]
			helper.currentUnblockable = false
			helper.refreshRoutes()
		}

	})

	helper.targetDrop = drop
	helper.targetBox.Clear()
	helper.targetBox.Append(drop.AsPtr())

}

func (helper *Helper) rebuildHosts() {

	if helper.hostsBox == nil {
		return
	}

	helper.hosts = make([]*insights.HostsFile, 0)
	labels := make([]string, 0)

	if helper.dataset != nil {

		for _, file := range helper.dataset.HostsFiles {
			helper.hosts = append(helper.hosts, file)
			labels = append(labels, filepath.Base(file.RelPath))
		}

	}

	if len(labels) == 0 {
		labels = append(labels, "(no hosts files)")
	}

	drop := gtk.NewDropDown(labels)
	drop.SetSizeRequest(130, -1)
	helper.hostsDrop = drop
	helper.hostsBox.Clear()
	helper.hostsBox.Append(drop.AsPtr())

}

func (helper *Helper) rebuildDomains(analysis *structs.Analysis) {

	if helper.domainBox == nil {
		return
	}

	helper.domainValues = make([]string, 0)
	labels := make([]string, 0)

	if analysis != nil && analysis.Email != nil {

		domains := make([]string, 0)

		if strings.Contains(analysis.Email.From, "@") == true {

			from_domain := analysis.Email.From[strings.Index(analysis.Email.From, "@")+1:]

			if from_domain != "" {
				domains = append(domains, from_domain)
			}

		}

		domains = append(domains, analysis.Email.Domains...)

		for _, domain := range domains {

			if helper.dataset != nil && helper.dataset.MatchValue(domain) != nil {
				continue
			}

			found := false

			for _, value := range helper.domainValues {

				if value == domain {
					found = true
					break
				}

			}

			if found == false {
				helper.domainValues = append(helper.domainValues, domain)
				labels = append(labels, truncate(domain, 26))
			}

		}

	}

	if len(labels) == 0 {
		labels = append(labels, "(no domains)")
	}

	drop := gtk.NewDropDown(labels)
	drop.SetSizeRequest(130, -1)
	helper.domainDrop = drop
	helper.domainBox.Clear()
	helper.domainBox.Append(drop.AsPtr())

}

func (helper *Helper) refreshRoutes() {

	if helper.routeBox == nil {
		return
	}

	indicator := helper.currentIndicator()
	helper.routeCandidates = make([]structs.RouteCandidate, 0)

	if helper.dataset == nil || indicator == nil || (indicator.Kind != "ipv4" && indicator.Kind != "ipv6") {

		helper.currentUnblockable = false
		helper.routeBox.Clear()
		helper.routeBox.Append(gtk.NewDropDown([]string{"(select an IP indicator)"}).AsPtr())

		if helper.prefixEntry != nil {
			helper.prefixEntry.SetText("")
		}

		return

	}

	candidates := helper.dataset.RouteCandidates(helper.targetEntry, indicator.Value, indicator.Routes, indicator.AddressBlock)
	helper.routeCandidates = candidates

	labels := make([]string, 0)
	limit := len(candidates)

	if limit > 200 {
		limit = 200
	}

	for index := 0; index < limit; index++ {

		candidate := candidates[index]
		label := candidate.Prefix + " [" + candidate.Source + "]"

		if candidate.ContainsIP == true {
			label += " *"
		}

		if candidate.Missing == false {
			label += " (already covered)"
		} else {
			label += " (" + formatCount(candidate.AddressCount) + " addresses)"
		}

		labels = append(labels, label)

	}

	if len(labels) == 0 {
		labels = append(labels, "(no route candidates)")
	}

	drop := gtk.NewDropDown(labels)
	drop.SetSizeRequest(340, -1)

	drop.OnChanged(func() {

		if helper.updating == true {
			return
		}

		index := int(drop.GetSelected())

		if index >= 0 && index < len(helper.routeCandidates) {
			helper.prefixEntry.SetText(helper.routeCandidates[index].Prefix)
		}

	})

	helper.routeDrop = drop
	helper.routeBox.Clear()
	helper.routeBox.Append(drop.AsPtr())

	if len(candidates) > 0 {

		if helper.currentUnblockable == true {
			helper.prefixEntry.SetText("")
		} else {
			helper.prefixEntry.SetText(candidates[0].Prefix)
		}

	}

}

func (helper *Helper) updatePending() {

	if helper.pendingLabel == nil {
		return
	}

	if helper.dataset == nil {
		helper.pendingLabel.SetText("No dataset loaded")
		return
	}

	pending := helper.dataset.PendingFiles()

	if len(pending) == 0 {
		helper.pendingLabel.SetText("No pending changes")
	} else {
		helper.pendingLabel.SetText("Pending: " + strings.Join(pending, ", "))
	}

}

func (helper *Helper) addPrefix() {

	if helper.loading == true {
		helper.status("Still loading emails...")
		return
	}

	if helper.dataset == nil || helper.targetEntry == nil {
		helper.status("Select a target spammer entry first")
		return
	}

	value := strings.TrimSpace(helper.prefixEntry.Text())

	if value == "" {
		helper.status("Enter a prefix first")
		return
	}

	added, err := helper.dataset.AddNetwork(helper.targetEntry, value)

	if err != nil {
		helper.status(err.Error())
		return
	}

	if len(added) == 0 {
		helper.status("Prefix is already covered by " + helper.targetEntry.Spammer.Domain)
	} else {
		helper.status("Staged " + strings.Join(added, ", ") + " for " + helper.targetEntry.Spammer.Domain)
	}

	helper.updatePending()
	helper.refreshRoutes()

}

func (helper *Helper) addHost() {

	if helper.loading == true {
		helper.status("Still loading emails...")
		return
	}

	if helper.dataset == nil {
		helper.status("No dataset loaded")
		return
	}

	if helper.hostsDrop == nil || helper.domainDrop == nil {
		return
	}

	host_index := int(helper.hostsDrop.GetSelected())
	domain_index := int(helper.domainDrop.GetSelected())

	if host_index < 0 || host_index >= len(helper.hosts) {
		helper.status("Select a hosts file first")
		return
	}

	if domain_index < 0 || domain_index >= len(helper.domainValues) {
		helper.status("Select a domain first")
		return
	}

	file := helper.hosts[host_index]
	domain := helper.domainValues[domain_index]

	if helper.dataset.AddHost(file, domain) == true {
		helper.status("Staged host " + domain + " in " + file.RelPath)
	} else {
		helper.status("Host " + domain + " already exists in " + file.RelPath)
	}

	helper.updatePending()

}

func (helper *Helper) showDiff() {

	if helper.loading == true {
		helper.status("Still loading emails...")
		return
	}

	if helper.dataset == nil {
		helper.status("No dataset loaded")
		return
	}

	diffs := helper.dataset.Diff()

	if len(diffs) == 0 {
		helper.status("No pending dataset changes")
		return
	}

	content := make([]string, 0)

	for _, diff := range diffs {
		content = append(content, diff.Diff)
	}

	window := gtk.NewWindow(helper.app)
	window.SetTitle("Review dataset changes")
	window.SetDefaultSize(1040, 720)
	window.SetModal(true)
	window.SetTransientFor(helper.window)

	root := gtk.NewBox(gtk.OrientationVertical, 8)
	root.SetMarginStart(10)
	root.SetMarginEnd(10)
	root.SetMarginTop(10)
	root.SetMarginBottom(10)

	text := gtk.NewTextView()
	text.SetEditable(false)
	text.SetMonospace(true)
	text.SetCursorVisible(false)
	text.SetWrapModeWordChar()
	text.SetText(strings.Join(content, "\n"))

	scroll := gtk.NewScrolledWindow()
	scroll.SetPolicy(gtk.PolicyAutomatic, gtk.PolicyAutomatic)
	scroll.SetVExpand(true)
	scroll.SetChild(text.AsPtr())

	buttons := gtk.NewBox(gtk.OrientationHorizontal, 8)

	close_button := gtk.NewButton("Close")

	close_button.OnClick(func() {
		window.Close()
	})

	save_button := gtk.NewButton("Save")

	save_button.OnClick(func() {

		if helper.loading == true {
			helper.status("Still loading emails...")
			return
		}

		changed, err := helper.dataset.Save()

		if err != nil {
			helper.status("Cannot save dataset: " + err.Error())
			return
		}

		helper.status("Saved " + strings.Join(changed, ", "))

		helper.updatePending()
		helper.refreshVerdicts()
		helper.refreshEmailList()
		helper.updateDetail()
		helper.refreshRoutes()

		window.Close()

	})

	buttons.Append(close_button.AsPtr())
	buttons.Append(save_button.AsPtr())

	root.Append(scroll.AsPtr())
	root.Append(buttons.AsPtr())

	window.SetChild(root.AsPtr())
	window.Present()

}

func (helper *Helper) refreshVerdicts() {

	if helper.dataset == nil {
		return
	}

	for _, analysis := range helper.analyses {

		analysis.Verdict, analysis.Match, analysis.Reasons = helper.dataset.Classify(analysis.Email)

		for _, indicator := range analysis.Indicators {

			indicator.Known = nil
			indicator.Unblockable = nil

			match := helper.dataset.MatchValue(indicator.Value)

			if match != nil {

				if match.Unblockable == true {
					indicator.Unblockable = match
				} else {
					indicator.Known = match
				}

			}

		}

	}

}

func (helper *Helper) confirmArchive(as_spam bool) {

	if helper.loading == true {
		helper.status("Still loading emails...")
		return
	}

	analysis := helper.currentAnalysis()

	if analysis == nil {
		helper.status("No email selected")
		return
	}

	action := "Not Spam"

	if as_spam == true {
		action = "Spam"
	}

	message := "Archive \"" + filepath.Base(analysis.File) + "\" as " + action + "?"

	if as_spam == true && analysis.Verdict == structs.VerdictUnknown {
		message += "\n\nNo dataset entry matches this email. It will be archived, but future emails will still be unknown until you stage a dataset change."
	}

	if helper.dataset != nil && len(helper.dataset.PendingFiles()) > 0 {
		message += "\n\nThere are unsaved dataset changes. They will stay staged and can be saved later."
	}

	window := gtk.NewWindow(helper.app)
	window.SetTitle("Confirm")
	window.SetDefaultSize(520, 220)
	window.SetModal(true)
	window.SetTransientFor(helper.window)

	root := gtk.NewBox(gtk.OrientationVertical, 8)
	root.SetMarginStart(10)
	root.SetMarginEnd(10)
	root.SetMarginTop(10)
	root.SetMarginBottom(10)

	label := gtk.NewLabel(validUTF8(message))
	label.SetWrap(true)
	label.SetXAlign(0)
	label.SetVExpand(true)

	buttons := gtk.NewBox(gtk.OrientationHorizontal, 8)

	no_button := gtk.NewButton("No")

	no_button.OnClick(func() {
		window.Close()
	})

	yes_button := gtk.NewButton("Archive")

	yes_button.OnClick(func() {

		window.Close()
		helper.archiveCurrent(as_spam)

	})

	buttons.Append(no_button.AsPtr())
	buttons.Append(yes_button.AsPtr())

	root.Append(label.AsPtr())
	root.Append(buttons.AsPtr())

	window.SetChild(root.AsPtr())
	window.Present()

}

func (helper *Helper) showConfirm(title string, message string, accept_label string, on_accept func()) {

	window := gtk.NewWindow(helper.app)
	window.SetTitle(title)
	window.SetDefaultSize(520, 220)
	window.SetModal(true)
	window.SetTransientFor(helper.window)

	root := gtk.NewBox(gtk.OrientationVertical, 8)
	root.SetMarginStart(10)
	root.SetMarginEnd(10)
	root.SetMarginTop(10)
	root.SetMarginBottom(10)

	label := gtk.NewLabel(validUTF8(message))
	label.SetWrap(true)
	label.SetXAlign(0)
	label.SetVExpand(true)

	buttons := gtk.NewBox(gtk.OrientationHorizontal, 8)

	cancel := gtk.NewButton("Cancel")

	cancel.OnClick(func() {
		window.Close()
	})

	accept := gtk.NewButton(accept_label)

	accept.OnClick(func() {

		window.Close()
		on_accept()

	})

	buttons.Append(cancel.AsPtr())
	buttons.Append(accept.AsPtr())

	root.Append(label.AsPtr())
	root.Append(buttons.AsPtr())

	window.SetChild(root.AsPtr())
	window.Present()

}

func (helper *Helper) archiveCurrent(as_spam bool) {

	analysis := helper.currentAnalysis()

	if analysis == nil {
		return
	}

	archive_folder := filepath.Join(helper.mailsRoot, "archive")

	err := os.MkdirAll(archive_folder, 0755)

	if err != nil {
		helper.status("Cannot create archive folder: " + err.Error())
		return
	}

	target := uniquePath(filepath.Join(archive_folder, filepath.Base(analysis.File)))

	err = os.Rename(analysis.File, target)

	if err != nil {
		helper.status("Cannot archive email: " + err.Error())
		return
	}

	helper.analyses = append(helper.analyses[:helper.selected], helper.analyses[helper.selected+1:]...)

	if helper.selected >= len(helper.analyses) {
		helper.selected = len(helper.analyses) - 1
	}

	helper.selectedIndicator = -1

	helper.refreshEmailList()

	if helper.selected == -1 && len(helper.visible) > 0 {
		helper.selected = helper.visible[0]
	}

	helper.updateDetail()

	if as_spam == true {
		helper.status("Archived as spam: " + filepath.Base(analysis.File))
	} else {
		helper.status("Archived as not spam: " + filepath.Base(analysis.File))
	}

}

func (helper *Helper) startLookup() {

	if helper.loading == true {
		helper.status("Still loading emails...")
		return
	}

	if helper.running.Load() == true {
		helper.cancel.Store(true)
		helper.status("Cancelling after the current lookup...")
		return
	}

	analysis := helper.currentAnalysis()

	if analysis == nil {
		helper.status("No email selected")
		return
	}

	targets := make([]*structs.Indicator, 0)

	for _, indicator := range analysis.Indicators {

		if (indicator.Kind == "ipv4" || indicator.Kind == "ipv6") && indicator.Private == false {
			targets = append(targets, indicator)
		}

	}

	if len(targets) == 0 {
		helper.status("No public IP addresses to look up")
		return
	}

	helper.running.Store(true)
	helper.cancel.Store(false)
	helper.analyzeButton.SetLabel("Cancel")
	helper.progress.SetFraction(0)
	helper.status(fmt.Sprintf("Looking up %d IP addresses...", len(targets)))

	go func() {

		for index, indicator := range targets {

			if helper.cancel.Load() == true {
				break
			}

			ip := indicator.Value

			gtk.RunOnMain(func() {
				helper.status("Looking up " + ip + "...")
			})

			lookup := actions.LookupIP(ip, helper.cache)

			current := index + 1
			total := len(targets)

			gtk.RunOnMain(func() {

				helper.progress.SetFraction(float64(current) / float64(total))
				helper.status(fmt.Sprintf("Looked up %d/%d: %s", current, total, ip))
				helper.applyLookup(indicator, lookup)

			})

		}

		gtk.RunOnMain(func() {

			cancelled := helper.cancel.Load()

			helper.running.Store(false)
			helper.cancel.Store(false)
			helper.analyzeButton.SetLabel("Analyze IPs")

			if cancelled == true {
				helper.status("Lookup cancelled")
			}

		})

	}()

}

func (helper *Helper) applyLookup(indicator *structs.Indicator, lookup structs.IPLookup) {

	indicator.LookedUp = true
	indicator.PTR = lookup.PTR
	indicator.ForwardConfirmed = lookup.ForwardConfirmed
	indicator.ASN = lookup.ASN
	indicator.Organization = lookup.Organization
	indicator.NetworkName = lookup.NetworkName
	indicator.Registry = lookup.Registry
	indicator.Country = lookup.Country
	indicator.AddressBlock = lookup.AddressBlock
	indicator.Routes = lookup.Routes

	if len(lookup.Errors) > 0 {
		indicator.Error = strings.Join(lookup.Errors, "; ")
	} else {
		indicator.Error = ""
	}

	query := strings.Join([]string{lookup.Organization, lookup.NetworkName, strings.Join(lookup.PTR, " "), lookup.ASN}, " ")

	indicator.Candidates = helper.dataset.FuzzyCandidates(query)

	analysis := helper.currentAnalysis()

	if analysis != nil {

		helper.indicators.Clear()

		for _, item := range analysis.Indicators {
			appendWrappedRow(helper.indicators, indicatorText(item))
		}

		if helper.selectedIndicator >= 0 && helper.selectedIndicator < len(analysis.Indicators) {
			helper.indicators.SelectRow(helper.selectedIndicator)
		}

	}

	helper.selectTargetForCurrentIndicator()
	helper.refreshRoutes()

}

func (helper *Helper) selectBestTarget(candidates []structs.FuzzyCandidate) {

	if helper.targetDrop == nil || len(helper.targets) == 0 {
		return
	}

	for _, candidate := range candidates {

		if candidate.Unblockable == true {
			continue
		}

		for index, entry := range helper.targets {

			if entry.Spammer.Domain == candidate.Domain && entry.FilePath() == candidate.File {

				helper.updating = true
				helper.targetDrop.SetSelected(uint(index + 1))
				helper.targetEntry = entry
				helper.updating = false

				return

			}

		}

	}

}

func (helper *Helper) resetTarget() {

	helper.targetEntry = nil
	helper.currentUnblockable = false

	if helper.targetDrop != nil {

		helper.updating = true
		helper.targetDrop.SetSelected(0)
		helper.updating = false

	}

}

func (helper *Helper) selectTargetForCurrentIndicator() {

	indicator := helper.currentIndicator()

	helper.currentUnblockable = false

	if indicator == nil {
		return
	}

	if indicator.Unblockable != nil {
		helper.currentUnblockable = true
	}

	if len(indicator.Candidates) > 0 {

		top := indicator.Candidates[0]

		if top.Unblockable == true {

			helper.currentUnblockable = true
			helper.targetEntry = nil

			if helper.targetDrop != nil {

				helper.updating = true
				helper.targetDrop.SetSelected(0)
				helper.updating = false

			}

			helper.status("Top fuzzy match is unblockable: " + top.Domain + " (treat as Not Spam)")
			return

		}

	}

	if helper.currentUnblockable == true {

		helper.targetEntry = nil

		if helper.targetDrop != nil {

			helper.updating = true
			helper.targetDrop.SetSelected(0)
			helper.updating = false

		}

		helper.status("Unblockable provider network (treat as Not Spam)")
		return

	}

	if len(indicator.Candidates) == 0 && indicator.LookedUp == true {
		helper.status("No fuzzy spammer match; select a target entry or use the hosts fallback")
		return
	}

	if helper.targetEntry == nil {
		helper.selectBestTarget(indicator.Candidates)
	}

}

func truncate(value string, length int) string {

	runes := []rune(value)

	if len(runes) <= length {
		return value
	}

	return string(runes[0:length-1]) + "…"

}

func validUTF8(value string) string {
	return strings.ToValidUTF8(value, "\uFFFD")
}

func escapeMarkup(value string) string {

	value = validUTF8(value)
	value = strings.ReplaceAll(value, "&", "&amp;")
	value = strings.ReplaceAll(value, "<", "&lt;")
	value = strings.ReplaceAll(value, ">", "&gt;")

	return value

}

func formatCount(value uint64) string {

	if value == ^uint64(0) {
		return ">=2^64"
	}

	return fmt.Sprintf("%d", value)

}

func uniquePath(path string) string {

	if _, err := os.Stat(path); os.IsNotExist(err) == true {
		return path
	}

	for index := 1; index < 10000; index++ {

		candidate := fmt.Sprintf("%s.%d", path, index)

		if _, err := os.Stat(candidate); os.IsNotExist(err) == true {
			return candidate
		}

	}

	return path + "." + fmt.Sprintf("%d", os.Getpid())

}
