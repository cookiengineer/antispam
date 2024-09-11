
## How To Add A New Entry

Currently, the process for adding new entries goes like this, and is pretty manual as of
now. I'm working on automating this by creating an intelligent `whois` client, but that
takes some more work to be failsafe due to how the protocol doesn't have any schema that
any server adheres to.

1. Open the E-Mail headers (`View source`, `view headers` or similar)

2. Identify the spam-sending server by looking for `designates a.b.c.d as permitted sender`
   line, where `a.b.c.d` represents the email-forwarding server's IP. The most top line is
   the last server/relay.

```
spf=pass (myserver.tld: domain of bounce+123456.123456-cookiengineer=myserver.tld@mg-d1.substack.com designates 159.112.244.6 as permitted sender) smtp.mailfrom="bounce+123456.123456-cookiengineer=myserver.tld@mg-d1.substack.com";
```

In the above example this would mean we want to block `substack.com` and the IP CIDR range
behind `159.112.244.6`. If the server is a `bounce.sender.tld` looking domain, chances are
pretty high it was sent by either `sendgrid` or `elasticemail` which share a lot of failover
IP ranges, and use proxy companies to hide their ASN associations.

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

4. Edit the `source/insights/spammers/<provider>.json` file, where `<provider>` is the
   correct provider. In this case it's the file `mailgun.json`.

In the file itself, we can also map aliases for known spam-sending services, so that we can
identify shared failover IPs more easily. Those are important if they rotate their spam
sending server IPs, which a lot of services do.

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

5. If we look up the corresponding `ASN`, we'll find out the IP range is assigned within
   [AS396479](https://ipinfo.io/AS396479) and we can now find out other IP ranges corresponding
   to this organization.

6. Alternatively, the ASN data can be gathered from the RIR organizations directly, and
   the organizations can be mapped back to the globally unique ASNs. As the data is
   anonymized, each organization contains a specific hash associating ASNs with IPv4/IPv6
   IP ranges.

Note that IPv4 ranges on RIPE are not prefix-based and instead use a start IPv4 and an
amount of IPs which can potentially overlap, and only be represented by multiple prefixes.

```bash
wget "https://ftp.afrinic.net/pub/stats/afrinic/delegated-afrinic-extended-latest";
wget "https://ftp.apnic.net/pub/stats/apnic/delegated-apnic-extended-latest";
wget "https://ftp.arin.net/pub/stats/arin/delegated-arin-extended-latest";
wget "https://ftp.lacnic.net/pub/stats/lacnic/delegated-lacnic-extended-latest";
wget "https://ftp.ripe.net/pub/stats/ripencc/delegated-ripencc-extended-latest";
```

