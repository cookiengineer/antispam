#!/bin/bash

TARGET="";
GO="$(which go)";
ROOT="$(pwd)";

if [[ "${GO}" != "" ]]; then

	echo "${ROOT}";
	cd "${ROOT}/bin";

	${GO} run "./main.go" "${ROOT}/build";

	if [[ "$?" == "0" ]]; then
		echo -e "Everything okay. [\e[32mok\e[0m]";
		exit 0;
	else
		echo -e "Unexpected error occured :( [\e[31mfail\e[0m]";
		exit 1;
	fi;

else
	echo -e "Please install go(lang) compiler. [\e[31mfail\e[0m]";
	exit 1;
fi;
