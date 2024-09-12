#!/bin/bash

GO="$(which go)";
SSH="$(which ssh)";
ROOT="$(pwd)";
GUEST_HOST="${1}";
GUEST_PORT="22";

if [[ "${GUEST_HOST}" == *:* ]]; then

	tmp1=$(echo "${GUEST_HOST}" | cut -d":" -f1);
	tmp2=$(echo "${GUEST_HOST}" | cut -d":" -f2);

	GUEST_HOST="$tmp1";
	GUEST_PORT="$tmp2";

fi;

if [[ "${GO}" != "" ]] && [[ "${SSH}" != "" ]]; then

	if [[ "${GUEST_HOST}" != "" ]]; then

		cd "${ROOT}/toolchain";

		${GO} run postfix.go generate;

		if [[ "$?" == "0" ]]; then

			cat "${ROOT}/build/blocked_clients" | ssh "root@${GUEST_HOST}" -p "${GUEST_PORT}" "cat > /etc/postfix/blocked_clients";
			cat "${ROOT}/build/blocked_senders" | ssh "root@${GUEST_HOST}" -p "${GUEST_PORT}" "cat > /etc/postfix/blocked_senders";

			if [[ "$?" == "0" ]]; then
				echo -e "- Copied blocklists to server. [\e[32mok\e[0m]";
			else
				echo -e "! Could not copy blocklists to server. [\e[31mfail\e[0m]";
			fi;

			ssh "root@${GUEST_HOST}" -p "${GUEST_PORT}" "postmap /etc/postfix/blocked_clients";
			ssh "root@${GUEST_HOST}" -p "${GUEST_PORT}" "postmap /etc/postfix/blocked_senders";

			if [[ "$?" == "0" ]]; then
				echo -e "- Postmapped blocklists on server. [\e[32mok\e[0m]";
			else
				echo -e "! Could not postmap blocklists on server. [\e[31mfail\e[0m]";
			fi;

			ssh "root@${GUEST_HOST}" -p "${GUEST_PORT}" "systemctl restart postfix";

			if [[ "$?" == "0" ]]; then
				echo -e "- Restarted postfix on server. [\e[32mok\e[0m]";
			else
				echo -e "! Could not restart postfix on server. [\e[31mfail\e[0m]";
			fi;

		fi;

	else
		echo -e "! Missing HOST parameter. [\e[31mfail\e[0m]";
		echo -e "";
	fi;

fi;

