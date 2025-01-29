
# Postfix Installation Guide

# 1. Setup Postfix

Postfix uses several options to implement its `restrictions` dependent on whether postfix is used
as a mail relay (please don't do that) or a separate mail transport agent that's running on a
specific FQDN (fully qualified domain name).

My recommendations for client and sender restrictions are the following:

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

Explanations: It is possible to integrate `postmap` lookup tables to block spam mails from a list
of known IP addresses and/or domains.

- `/etc/postfix/blocked_clients` and `check_client_access` refers to the mail client information that is connected via the `smtp` protocol.
- `/etc/postfix/blocked_senders` and `check_sender_access` refers to the `From:` field in the EML / E-Mail "file" itself.


## 2. Generate Postmap Files

Postmap files can be generated into the build folder, from the [source/insights](../source/insights)
dataset:

```bash
# Save postmap files into /build
cd /path/to/antispam/toolchain;
go run postfix.go;
```

## 3. Upload Postmap Files

If you followed the modifications above of the `main.cf`, you need to now upload the postmap files
to the server. For example, this can be done via `scp`:

```bash
cd /path/to/antispam/build;
scp blocked_clients root@your.server.tld:/etc/postfix/blocked_clients;
scp blocked_senders root@your.server.tld:/etc/postfix/blocked_senders;
```

## 4. Restart Postfix Service

On the server, you need to process the postmap files via the `postmap` command and then restart the
postfix service for the changes to take effect:

```bash
ssh root@your.server.tld;
postmap /etc/postfix/blocked_clients;
postmap /etc/postfix/blocked_senders;

systemctl restart postfix;
```

