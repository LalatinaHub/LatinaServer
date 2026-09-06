package warp

// Predefined domain lists for AI, Streaming, and Censored services.
var (
	// AIDomains contains domain suffixes for OpenAI, Claude/Anthropic, etc.
	AIDomains = []string{
		"openai.com",
		"chatgpt.com",
		"oaistatic.com",
		"oaiusercontent.com",
		"anthropic.com",
		"claude.ai",
		"meta.ai",
		"cohere.com",
		"sora.com",
	}

	// StreamingDomains contains domain suffixes for Netflix, Disney+, Prime, Hulu, HBO Max, etc.
	StreamingDomains = []string{
		"netflix.com",
		"netflix.net",
		"nflximg.net",
		"nflxvideo.net",
		"nflxso.net",
		"nflxext.com",
		"disneyplus.com",
		"disney-portal.my.onetrust",
		"dssott.com",
		"bamgrid.com",
		"primevideo.com",
		"amazonvideo.com",
		"hulu.com",
		"hulustream.com",
		"hbomax.com",
		"max.com",
	}

	// UnblockDomains contains domains frequently blocked or datacenter-flagged (Reddit, etc.)
	UnblockDomains = []string{
		"reddit.com",
		"redditmedia.com",
		"redd.it",
	}
)

// GetAllUnlockerDomains returns a combined, deduplicated list of domains targeted for WARP routing.
func GetAllUnlockerDomains() []string {
	seen := make(map[string]bool)
	var result []string

	for _, d := range AIDomains {
		if !seen[d] {
			seen[d] = true
			result = append(result, d)
		}
	}

	for _, d := range StreamingDomains {
		if !seen[d] {
			seen[d] = true
			result = append(result, d)
		}
	}

	for _, d := range UnblockDomains {
		if !seen[d] {
			seen[d] = true
			result = append(result, d)
		}
	}

	return result
}
