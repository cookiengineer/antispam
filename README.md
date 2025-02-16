
# antispam

This project is my attempt to fix the LLM spamming and phishing problem.

This repository reflects the spam that I or the networks under my protection receive on
a daily basis. It also includes providers that refuse to abide with the `abuse@` report
policies because they have a conflict of interest to their paying customers (obviously).

![screenshot](./guides/screenshot.png)

## Dataset

The dataset is maintained in the [insights](./source/insights) folder. It differs between the
following types of entries:

- [hosts](./source/insights/hosts) which are specific domains that cannot be blocked by other means (yet).
- [spammers](./source/insights/spammers) which are known to offer spam campaigns as a service (usually under the umbrella of AI targeted marketing campaigns).
- [spammers/unblockable](./source/insights/spammers/unblockable) which are known to send a lot of spam via e.g. Google Forms Scams or Microsoft Azure to bypass Microsoft Exchange Filters.
- [phishers](./source/insights/phishers) which are known to send phishing and malware campaigns.

The abstractions behind the scenes use a form of longest-prefix hashset maps which allow this
to be computable in a faster manner than with tries. The standalone library is available as
the [golpm](https://github.com/cookiengineer/golpm) project.


## Usage

The usage of this tool is intended to be used via `cronjob`s or via `exec` or `filter` event.
It can be used manually on the `eml` files or directly on incoming email buffers via `stdin`.
See the below Postfix and Dovecot Configuration sections for more details.

```bash
# View usage help
antispam;

# View an email and check for spam indicators
antispam view path/to/mail.eml;

# If an email is spam, return exit code 1
antispam classify path/to/mail.eml;

# Mark an email as spam, output new spammer ready for pull-request
antispam mark --json path/to/mail.eml;
```


## Building

The build [toolchain](./toolchain) is implemented in `go`, so you only need to install `go` first.

```bash
# Install go compiler/framework
sudo pacman -S go;

# Build binary
cd /path/to/antispam/toolchain;
go run build.go;

# Execute binary
cd /path/to/antispam/build;
./antispam-linux-amd64;
```


## Toolchain

The following tools are available to manage large folders of spam/malware/phishing emails. In order
to use them, copy your email files to the [mails](./mails) folder and run the scripts afterwards.

```bash
# Cleanup spam
cd /path/to/antispam/toolchain;
go run cleanup.go --spam;

# Cleanup from allowlisted domains
go run cleanup.go --from="@example.com";

# Show whether E-Mails are classified as Spam or NotSpam
go run learn.go;
```

[build.go](./toolchain/build.go):

- `go run build.go` generates binaries for `linux/amd64` into the [build](./build) folder.
- `go run build.go --debug` generates binaries with debug symbols into the [build](./build) folder.

[cleanup.go](./toolchain/cleanup.go):

The cleanup script removes all files in [mails](./mails) that match its selection criteria.

This is intended to skip emails from trusted sources, no matter whether they are classified
or spam or not. Take a look at the [cleanup.sh](./cleanup.sh) for examples.

- `go run cleanup.go --from="johndoe@example.com"`
- `go run cleanup.go --from="@example.com"`
- `go run cleanup.go --domain="example.com"`
- `go run cleanup.go --spam`

[discover.go](./toolchain/discover.go):

(Currently work-in-progress, so use with manual oversight)

The discover script tries to combat the shell game that certain cyber terrorist nations are playing
to avoid international sanctions. This is intented to discover e.g. Russian, Iranian or Chinese
phishing/fraud companies that are part of international scam and malware campaigns.

- `go run discover.go --domain="example.com"` discovers neighboring ASNs that are likely spam providers, too.

[learn.go](./toolchain/learn.go):

- `go run learn.go` classifies all files in [mails](./mails) and shows the reasons.

[postfix.go](./toolchain/postfix.go):

- `go run postfix.go` generates `postmap` compatible files into the [build](./build) folder.


## Postfix Configuration

The Postfix configuration is documented in [POSTFIX.md](./guides/POSTFIX.md) and uses
external `postmap` blocklists to block network prefixes and domains. Postmaps files
are shrinked, meaning the shortest prefix length prevails (e.g. `1.2.3.0/24` will be
removed if `1.2.0.0/16` is blocked anyways).


## Dovecot Configuration

The Dovecot configuration is documented in [DOVECOT.md](./guides/DOVECOT.md) and uses
an external `sieve` script to pipe incoming mails to the [antispam-sieve](./source/cmds/antispam-sieve/main.go) wrapper.


## Pull Requests

Pull Requests are certainly welcome! I don't like spam, and so do you, I guess?
So let's fight spam together!

If you want to contribute a new Spammer entry (generated via `antispam mark --json <file>`,
please make sure to use the same naming scheme for the files.

Each spammer organization has a separate JSON file, containing an Array of [structs.Spammer](./source/structs/Spammer.go).
For example, [Amazon](./source/insights/spammers/amazon.json) contains the `structs.Spammer`
instances for Amazon US, Amazon EU, Amazon JP etc.


# License

AGPL-3

