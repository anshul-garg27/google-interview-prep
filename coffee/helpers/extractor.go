package helpers

import (
	"fmt"
	"regexp"
)

func NewProfileExtractionConfig() *ExtractionConfig {
	return &ExtractionConfig{
		Patterns: []Pattern{
			{
				Platform: "INSTAGRAM",
				Type:     "HANDLE",
				Regex:    `instagram\.com/([A-Za-z0-9_.]+)`,
			},
			{
				Platform: "YOUTUBE",
				Type:     "HANDLE",
				Regex:    `youtube\.com/c/(@?[A-Za-z0-9_-]+)`,
			},
			{
				Platform: "YOUTUBE",
				Type:     "CHANNEL_ID",
				Regex:    `youtube\.com/channels/([A-Za-z0-9_-]+)`,
			},
			{
				Platform: "YOUTUBE",
				Type:     "HANDLE",
				Regex:    `youtube\.com/(@?[A-Za-z0-9_-]+)`,
			},
		},
	}
}

func ExtractPlatformProfileIdentifier(url string) (string, string, string, error) {
	config := NewProfileExtractionConfig()
	for _, pattern := range config.Patterns {
		re := regexp.MustCompile(pattern.Regex)
		matches := re.FindStringSubmatch(url)
		if len(matches) == 2 {
			return pattern.Platform, pattern.Type, matches[1], nil
		}
	}
	return "", "", "", fmt.Errorf("unable to extract shortcode")
}

type Pattern struct {
	Platform string
	Type     string
	Regex    string
}

type ExtractionConfig struct {
	Patterns []Pattern
}

func NewPostExtractionConfig() *ExtractionConfig {
	return &ExtractionConfig{
		Patterns: []Pattern{
			{
				Platform: "INSTAGRAM",
				Regex:    `instagram\.com/p/([A-Za-z0-9_-]+)`,
			},
			{
				Platform: "INSTAGRAM",
				Type:     "REEL",
				Regex:    `instagram\.com/reel/([A-Za-z0-9_-]+)`,
			},
			{
				Platform: "INSTAGRAM",
				Type:     "STORY",
				Regex:    `instagram\.com/stories/[a-zA-Z0-9]+/([A-Za-z0-9_-]+)`,
			},
			{
				Platform: "YOUTUBE",
				Regex:    `youtube\.com/watch\?v=([A-Za-z0-9_-]+)`,
			},
			{
				Platform: "YOUTUBE",
				Regex:    `youtube\.com/watch\?v%3D([A-Za-z0-9_-]+)`,
			},
			{
				Platform: "YOUTUBE",
				Regex:    `youtube\.com/v/([A-Za-z0-9_-]+)`,
			},
			{
				Platform: "YOUTUBE",
				Type:     "SHORT",
				Regex:    `youtube\.com/shorts/([A-Za-z0-9_-]+)`,
			},
			{
				Platform: "YOUTUBE",
				Regex:    `youtu\.be/watch\?v=([A-Za-z0-9_-]+)`,
			},
			{
				Platform: "YOUTUBE",
				Regex:    `youtu\.be/v/([A-Za-z0-9_-]+)`,
			},
			{
				Platform: "YOUTUBE",
				Type:     "SHORT",
				Regex:    `youtu\.be/shorts/([A-Za-z0-9_-]+)`,
			},
			{
				Platform: "YOUTUBE",
				Regex:    `youtu\.be/([A-Za-z0-9_-]+)`,
			},
			{
				Platform: "YOUTUBE",
				Regex:    `youtube\.com/attribution_link\?.*v=([A-Za-z0-9_-]+)`,
			},
			{
				Platform: "YOUTUBE",
				Regex:    `youtube\.com/attribution_link\?.*v%3D([A-Za-z0-9_-]+)`,
			},
			{
				Platform: "YOUTUBE",
				Regex:    `youtube\.com/embed/([A-Za-z0-9_-]+)`,
			},
			{
				Platform: "YOUTUBE",
				Regex:    `youtube\.com/e/([A-Za-z0-9_-]+)`,
			},
		},
	}
}

func ExtractPlatformPostIdentifier(url string) (string, string, error) {
	config := NewPostExtractionConfig()
	for _, pattern := range config.Patterns {
		re := regexp.MustCompile(pattern.Regex)
		matches := re.FindStringSubmatch(url)
		if len(matches) == 2 {
			return pattern.Platform, matches[1], nil
		}
	}
	return "", "", fmt.Errorf("unable to extract shortcode")
}
