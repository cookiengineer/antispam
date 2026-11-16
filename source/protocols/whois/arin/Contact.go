package arin

type Contact struct {
	Identifier string `whois:"OrgAbuseHandle,OrgTechHandle,OrgRoutingHandle"`
	Name       string `whois:"OrgAbuseName,OrgTechName,OrgRoutingName"`
	Phone      string `whois:"OrgAbusePhone,OrgTechPhone,OrgRoutingPhone"`
	Email      string `whois:"OrgAbuseEmail,OrgTechEmail,OrgRoutingEmail"`
	Reference  string `whois:"OrgAbuseRef,OrgTechRef,OrgRoutingRef"`
}
