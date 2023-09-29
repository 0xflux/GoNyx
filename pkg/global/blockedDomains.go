package global

/*
Hold any blocked domains for privacy reasons, e.g. automatic packets to the
Mozzila foundation.
*/

// TODO implement
func GetBlockedDomains() []string {
	blockedDomains := []string{
		"services.mozilla.com",
	}

	return blockedDomains
}
