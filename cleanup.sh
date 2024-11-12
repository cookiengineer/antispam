#!/bin/bash

ROOT=$(pwd);

cd "${ROOT}/toolchain";

# private email hosters
go run cleanup.go --from="@gmail.com";
go run cleanup.go --from="@gmx.de";
go run cleanup.go --from="@gmx.net";
go run cleanup.go --from="@outlook.com";
go run cleanup.go --from="@web.de";

# foundations and bug reports
go run cleanup.go --from="@ccc.de";
go run cleanup.go --from="@cert-bund.de";
go run cleanup.go --from="@eff.org";
go run cleanup.go --from="@gulas.ch";
go run cleanup.go --from="@ietf.org";
go run cleanup.go --from="@kernel.org";
go run cleanup.go --from="@mitre.org";
go run cleanup.go --from="@w3c.org";
go run cleanup.go --from="@zendesk.com";

# federal governemnt
go run cleanup.go --from="@bka.bund.de";
go run cleanup.go --from="@bsi.bund.de";
go run cleanup.go --from="@fbi.gov";
go run cleanup.go --from="@interpol.int";

go run cleanup.go --from="@polizei.bw.de";
go run cleanup.go --from="@polizei.bayern.de";
go run cleanup.go --from="@polizei.berlin.de";
go run cleanup.go --from="@polizei.brandenburg.de";
go run cleanup.go --from="@polizei.bremen.de";
go run cleanup.go --from="@polizei.hamburg.de";
go run cleanup.go --from="@polizei.landsh.de";
go run cleanup.go --from="@polizei.nrw.de";
go run cleanup.go --from="@polizei.rlp.de";
go run cleanup.go --from="@polizei.sachsen-anhalt.de";
go run cleanup.go --from="@polizei.slpol.de";
go run cleanup.go --from="@polizei.thueringen.de";

