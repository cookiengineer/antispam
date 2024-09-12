#!/bin/bash

ROOT=$(pwd);

SAMPLE_AMAZON="${ROOT}/mails/1713204853.M587868P232608.mail,S=22856,W=23580:2,S";
SAMPLE_SPF="${ROOT}/mails/1719856904.M221548P1609296.mail\,S\=10890\,W\=11052\:2\,S";


cd "${ROOT}/source";

go run cmds/antispam/main.go view "${SAMPLE_AMAZON}";
