package linkparser

import "regexp"

var wikiLinkRegex = regexp.MustCompile(`\[\[([^\]]+)\]\]`)

// ExtractLinks extracts all [[wiki-link]] targets from markdown content.
func ExtractLinks(content string) []string {
	matches := wikiLinkRegex.FindAllStringSubmatch(content, -1)
	seen := make(map[string]bool)
	var links []string

	for _, match := range matches {
		if len(match) > 1 {
			link := match[1]
			if !seen[link] {
				seen[link] = true
				links = append(links, link)
			}
		}
	}

	return links
}
