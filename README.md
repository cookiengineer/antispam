
# antispam

This project is my attempt to fix the spam problem.

Currently this repository reflects the spam that I or my customers receive on a daily
basis, whereas it also includes providers that refuse to abide with the `abuse` report
policies because they have a conflict of interest to their paying customers (obviously).

![screenshot](./guides/screenshot.png)

## Usage

The usage of this tool is intended to be used via `cronjob`s or via `exec` on event.

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


## Postfix Configuration

The Postfix configuration is documented in [POSTFIX.md](./guides/POSTFIX.md) and uses
external `postmap` blocklists to block network prefixes and domains.

```bash
cd /path/to/antispam/toolchain;

# Upload and install postmap files
go run postfix.go install root@your.server.tld:2222;
```


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

