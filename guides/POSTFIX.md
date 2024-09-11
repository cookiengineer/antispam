
# Postfix Configuration

In postfix, it's important to restrict your `mail transport agent` to your own networks
when you self-host your email server. My recommendations are the following:

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

You are able to integrate local `postmap` lookup tables to block spam mails from a list
of known senders.

- `/etc/postfix/blocked_clients` and `check_client_access` refers to the mail client (using the `smtp` protocol).
- `/etc/postfix/blocked_senders` and `check_sender_access` refers to the `From:` field in the e-mail "file" itself.


## Postmap Generation

```bash
cd /path/to/antispam/toolchain;

# Generate postmap files in /build
go run postfix.go generate;

# Upload and install postmap files
# Then execute "postmap" and "systemctl restart postfix"
go run postfix.go install root@your.server.tld:2222;
```

