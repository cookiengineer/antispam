
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


# License

WTFPL

