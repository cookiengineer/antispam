
# Dovecot Installation Guide

Dovecot uses `pigeonhole` as its `sieve` implementation.

## 1. Setup Sieve

Integrate the `sieve` plugin into the `/etc/dovecot/dovecot.conf` like this:

```conf
plugin {
    sieve_extensions = +vnd.dovecot.execute
    sieve_plugins = sieve_extprograms
    sieve_execute_bin_dir = /etc/dovecot/sieve
}
```

Then create the folder `/etc/dovecot/sieve` and the file `/etc/dovecot/sieve/antispam.sv`
with the following contents:

```sieve
require "vnd.dovecot.execute";

if not execute :pipe /usr/bin/antispam-sieve {
    discard;
    stop;
}
```

## 2. Install Antispam-Sieve

On your local machine, build the binaries:

```bash
cd /path/to/antispam/toolchain;

go run build.go;
```

Afterwards, upload the binary located at `/path/to/antispam/build/antispam-sieve_linux_amd64`
to the server as `/usr/bin/antispam-sieve` and make sure it's executable.


## 3. Restart Dovecot Service

```bash
sudo systemctl restart dovecot;
```
