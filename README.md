
# postfix-spamdb

This is my attempt to fix the postfix spam problem, by having a locally running database
with dead-simple entries. Currently, the process is not really automated and just reflects
the spam that I receive personally on a daily basis.

In the future, it is planned to automatically integrate a WHOIS workflow to be able to
block networks that send spam. As I have zero tolerance for spam, hosting providers that
violate spam policies and don't give a damn about spam when being reached out to via their
`abuse email` entry are also blocked.


# Postfix Integration

In postfix, it's important to restrict your `mail transport agent` to your own networks when
you self-host your email server. My recommendations are the following:

```postfix
# /etc/postfix/main.cf

tls_preempt_cipherlist = no
tls_ssl_options = NO_COMPRESSION
tls_medium_cipherlist = ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384

smtpd_helo_required = yes
smtpd_helo_restrictions = permit_mynetworks permit_sasl_authenticated reject_invalid_helo_hostname reject_non_fqdn_helo_hostname
smtpd_relay_restrictions = permit_mynetworks permit_sasl_authenticated reject_unauth_destination
smtpd_sender_restrictions = permit_sasl_authenticated reject_non_fqdn_sender reject_unknown_sender_domain reject_sender_login_mismatch reject_unlisted_sender reject_rhsbl_sender dbl.spamhaus.org
smtpd_recipient_restrictions = permit_mynetworks permit_sasl_authenticated reject_unknown_recipient_domain reject_non_fqdn_recipient reject_unlisted_recipient reject_rbl_client zen.spamhaus.org check_client_access hash:/etc/postfix/blocked_clients check_sender_access hash:/etc/postfix/blocked_senders
```

As you can see, you are able to integrate local `postmap` lookup tables to block spamming senders.

- `check_client_access` refers to the mail client (using the `smtp` protocol).
- `check_sender_access` refers to the `From:` field in the e-mail "file" itself.


## Usage / Installation

Automated usage using the [install.sh](/install.sh) script:

```bash
# Generate postmap files
bash make.sh;

# Upload postmap files to /etc/postfix/{blocked_clients,blocked_senders} automagically
bash install.sh your.server.tld:2222;
```

Manual usage:

```bash
# Generate postmap files
bash make.sh;

# Upload and generate postmap files
scp ./build/blocked_clients root@your.server.tld:/etc/postfix/blocked_clients;
scp ./build/blocked_senders root@your.server.tld:/etc/postfix/blocked_senders;

# On the server, run postmap and restart postfix
ssh root@your.server.tld;
postmap /etc/postfix/blocked_clients;
postmap /etc/postfix/blocked_senders;

systemctl restart postfix;
```


## Adding a Pull Request / New Entries

Currently, the process for adding new entries goes like this, and is pretty manual as of now. I'm
working on automating this by creating an intelligent `whois` client, but that takes some more work
to be failsafe due to how the protocol doesn't have any schema that any server adheres to.

1. Open the E-Mail headers (`View source`, `view headers` or similar)

2. Identify the spam-sending server by looking for `designates a.b.c.d as permitted sender`
   line, where `a.b.c.d` represents the email-forwarding server's IP. The most top line is
   the last server/relay.

```
spf=pass (myserver.tld: domain of bounce+123456.123456-cookiengineer=myserver.tld@mg-d1.substack.com designates 159.112.244.6 as permitted sender) smtp.mailfrom="bounce+123456.123456-cookiengineer=myserver.tld@mg-d1.substack.com";
```

In the above example this would mean we want to block `substack.com` and the IP CIDR range behind `159.112.244.6`.
If the server is a `bounce.sender.tld` looking domain, chances are pretty high it was sent by either `sendgrid` or
`elasticemail` which share a lot of failover IP ranges, and use proxy companies to hide their ASN associations.

3. Use the `whois` CLI command to find out who owns the IPv4/IPv6 CIDR range

```bash
whois -h whois.iana.org 159.112.244.6;
# registered by ARIN, status LEGACY implying it could be delegated to another RIR
# In case it is, we need to use the correct server

whois -h whois.iana.org 159.112.244.6;
# CIDR: 159.112.240.0/20
# OrgName: Mailgun Technologies Inc.
# OrgId: MT-757
```

4. Edit the [source/provider-name.json](./source) file, where `provider-name` is the correct provider. In this case
   it's `mailgun.json`.

In the file itself, we can also map aliases for known spam-sending services, so that we can identify shared failover
IPs more easily. Those are important if they rotate their spam sending server IPs, which a lot of services do.

```json
[{
    "domain": "mailgun.net",
    "aliases": [
        "mailgun.com",
        "substack.com"
    ],
    "networks": [
        "159.112.240.0/20"
    ]
}]
```

5. If we look up the corresponding `ASN`, we'll find out the IP range is assigned within [AS396479](https://ipinfo.io/AS396479)
   and we can now find out other IP ranges corresponding to this organization.

6. Alternatively, the ASN data can be gathered from the RIR organizations directly, and the organizations can be mapped back
   to the globally unique ASNs. As the data is anonymized, each organization contains a specific hash associating ASNs with
   IPv4/IPv6 IP ranges. Note that IPv4 ranges on RIPE are not prefix-based (and instead use a start IPv4 and an amount of IPs)
   and can potentially overlap.

```bash
wget "https://ftp.afrinic.net/pub/stats/afrinic/delegated-afrinic-extended-latest";
wget "https://ftp.apnic.net/pub/stats/apnic/delegated-apnic-extended-latest";
wget "https://ftp.arin.net/pub/stats/arin/delegated-arin-extended-latest";
wget "https://ftp.lacnic.net/pub/stats/lacnic/delegated-lacnic-extended-latest";
wget "https://ftp.ripe.net/pub/stats/ripencc/delegated-ripencc-extended-latest";
```


# License

WTFPL

