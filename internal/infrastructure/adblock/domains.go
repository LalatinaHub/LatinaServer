package adblock

// Predefined domain suffixes for advertising, telemetry, malware, and gambling.
var (
	// AdAndTrackingDomains contains high-frequency ad networks and user tracking domains.
	AdAndTrackingDomains = []string{
		"doubleclick.net",
		"googleadservices.com",
		"googlesyndication.com",
		"adservice.google.com",
		"adnxs.com",
		"ads-twitter.com",
		"ads.facebook.com",
		"advertising.apple.com",
		"taboola.com",
		"outbrain.com",
		"popads.net",
		"propellerads.com",
		"criteo.com",
		"rubiconproject.com",
		"pubmatic.com",
		"openx.net",
		"scorecardresearch.com",
		"quantserve.com",
		"zedo.com",
		"inmobi.com",
		"flurry.com",
		"appsflyer.com",
		"branch.io",
		"adjust.com",
		"unityads.unity3d.com",
		"vungle.com",
		"applovin.com",
		"ironsrc.com",
		"adcolony.com",
		"chartboost.com",
		"admob.com",
	}

	// MalwareAndPhishingDomains contains known threat, malware command-and-control, and phishing domains.
	MalwareAndPhishingDomains = []string{
		"malware-traffic-analysis.net",
		"phishtank.com",
		"cybercrime-tracker.net",
		"telemetry.microsoft.com",
		"vortex.data.microsoft.com",
	}

	// GamblingAndScamDomains contains online gambling and scam portals.
	GamblingAndScamDomains = []string{
		"slot88.com",
		"sbobet.com",
		"maxbet.com",
		"pragmaticplay.com",
		"pgsoft.com",
		"judionline.net",
		"judibola.com",
		"rajatogel.com",
		"taruhanbola.com",
	}
)

// GetAllAdblockDomains returns combined, deduplicated domain suffixes.
func GetAllAdblockDomains() []string {
	seen := make(map[string]bool)
	var result []string

	for _, d := range AdAndTrackingDomains {
		if !seen[d] {
			seen[d] = true
			result = append(result, d)
		}
	}

	for _, d := range MalwareAndPhishingDomains {
		if !seen[d] {
			seen[d] = true
			result = append(result, d)
		}
	}

	for _, d := range GamblingAndScamDomains {
		if !seen[d] {
			seen[d] = true
			result = append(result, d)
		}
	}

	return result
}
