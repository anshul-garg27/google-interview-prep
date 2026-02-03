package helpers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_ExtractPlatformShortcode(t *testing.T) {
	urls := []string{

		"https://www.instagram.com/p/CrzzgYNgNzR",
		"https://instagram.com/p/CrzzgYNgNzR",
		"https://www.youtube.com/watch?v=wy42At0lQvE",
		"https://www.youtube.com/v/wy42At0lQvE",
		"https://www.instagram.com/p/CrzzgYNgNzR#asdf",
		"https://www.youtube.com/watch?v=wy42At0lQvE#asdf",
		"https://www.youtube.com/v/wy42At0lQvE#asdf",
		"https://www.instagram.com/p/CrzzgYNgNzR?asdf=1",
		"https://www.youtube.com/watch?v=wy42At0lQvE?sdf=1",
		"https://www.youtube.com/v/wy42At0lQvE?234a=1",
		"https://youtu.be/v/wy42At0lQvE?234a=1",
		"https://youtu.be/watch?v=wy42At0lQvE?234a=1",
		"https://www.instagram.com/stories/msugamsingh/3124348276069682270/",
		"https://www.instagram.com/reel/CtbcnrxAhm7/",
		"https://www.instagram.com/p/Cr8LGMtB_dC/",
		"https://www.instagram.com/p/Cr6EOPABupF/",
		"https://youtu.be/watch?v=wy42At0lQvE?234a=1",
		"https://www.youtube.com/shorts/wy42At0lQvE?asdf=1",
		"https://youtube.com/shorts/wy42At0lQvE?asdf=1",
		"https://www.youtube.com/watch?v=lalOy8Mbfdc",
		"https://m.youtube.com/watch?v=lalOy8Mbfdc",
		"https://www.youtube.com/v/dQw4w9WgXcQ",
		"https://m.youtube.com/v/dQw4w9WgXcQ",
		"https://youtu.be/-wtIMTCHWuI",
		"https://www.youtube.com/attribution_link?a=JdfC0C9V6ZI&u=%2Fwatch%3Fv%3DEhxJLojIE_o%26feature%3Dshare",
		"https://youtube.com/oembed?url=http%3A//www.youtube.com/watch?v%3D-wtIMTCHWuI&format=json",
		"https://youtube.com/embed/lalOy8Mbfdc",
		"https://youtube.com/e/dQw4w9WgXcQ",
		"https://www.youtube.com/watch?v=lalOy8Mbfdc",
	}

	for _, url := range urls {
		fmt.Printf("URL: %s\n", url)
		platform, shortcode, err := ExtractPlatformPostIdentifier(url)
		if err != nil {
			fmt.Printf("Error extracting shortcode from URL: %v\n", err)
		} else {
			fmt.Printf("Platform: %s, Shortcode: %s\n\n", platform, shortcode)
		}
	}
}

func TestClient_ExtractPlatformHandle(t *testing.T) {
	urls := map[string]string{
		"https://www.instagram.com/virat.kohli#asdf":                       "INSTAGRAM - HANDLE - virat.kohli",
		"https://www.instagram.com/virat.kohli?avbds=true":                 "INSTAGRAM - HANDLE - virat.kohli",
		"https://www.youtube.com/channels/UChBmxf4t3Tz8DKVizooeSvA":        "YOUTUBE - CHANNEL_ID - UChBmxf4t3Tz8DKVizooeSvA",
		"https://www.youtube.com/channels/UChBmxf4t3Tz8DKVizooeSvA#asdf":   "YOUTUBE - CHANNEL_ID - UChBmxf4t3Tz8DKVizooeSvA",
		"https://www.youtube.com/channels/UChBmxf4t3Tz8DKVizooeSvA?abcd=1": "YOUTUBE - CHANNEL_ID - UChBmxf4t3Tz8DKVizooeSvA",
		"https://www.youtube.com/channels/UChBmxf4t3Tz8DKVizooeSvA/about":  "YOUTUBE - CHANNEL_ID - UChBmxf4t3Tz8DKVizooeSvA",
		"https://www.youtube.com/@tseries":                                 "YOUTUBE - HANDLE - @tseries",
		"https://www.youtube.com/c/tseries":                                "YOUTUBE - HANDLE - tseries",
		"https://www.youtube.com/@tseries/about":                           "YOUTUBE - HANDLE - @tseries",
		"https://www.youtube.com/@tseries/about?asdf=1":                    "YOUTUBE - HANDLE - @tseries",
		"https://www.youtube.com/@tseries?asdf=1":                          "YOUTUBE - HANDLE - @tseries",
	}

	for url, expectedOutput := range urls {
		platform, identifierType, identifier, err := ExtractPlatformProfileIdentifier(url)
		if err != nil {
			t.Errorf("Error extracting handle or channel ID from URL: %v", err)
		} else {
			assert.Equal(t, expectedOutput, fmt.Sprintf("%s - %s - %s", platform, identifierType, identifier))
			fmt.Printf("URL: %s\n%s - %s - %s\n\n", url, platform, identifierType, identifier)
		}
	}
}

func TestClient_ExtractPlatformShortcodeV2(t *testing.T) {
	urls := map[string]string{
		"https://www.instagram.com/p/CrzzgYNgNzR":                                                               "INSTAGRAM - SHORTCODE - CrzzgYNgNzR",
		"https://instagram.com/p/CrzzgYNgNzR":                                                                   "INSTAGRAM - SHORTCODE - CrzzgYNgNzR",
		"https://www.youtube.com/watch?v=wy42At0lQvE":                                                           "YOUTUBE - SHORTCODE - wy42At0lQvE",
		"https://www.youtube.com/v/wy42At0lQvE":                                                                 "YOUTUBE - SHORTCODE - wy42At0lQvE",
		"https://www.instagram.com/p/CrzzgYNgNzR#asdf":                                                          "INSTAGRAM - SHORTCODE - CrzzgYNgNzR",
		"https://www.youtube.com/watch?v=wy42At0lQvE#asdf":                                                      "YOUTUBE - SHORTCODE - wy42At0lQvE",
		"https://www.youtube.com/v/wy42At0lQvE#asdf":                                                            "YOUTUBE - SHORTCODE - wy42At0lQvE",
		"https://www.instagram.com/p/CrzzgYNgNzR?asdf=1":                                                        "INSTAGRAM - SHORTCODE - CrzzgYNgNzR",
		"https://www.youtube.com/watch?v=wy42At0lQvE?sdf=1":                                                     "YOUTUBE - SHORTCODE - wy42At0lQvE",
		"https://www.youtube.com/v/wy42At0lQvE?234a=1":                                                          "YOUTUBE - SHORTCODE - wy42At0lQvE",
		"https://youtu.be/v/wy42At0lQvE?234a=1":                                                                 "YOUTUBE - SHORTCODE - wy42At0lQvE",
		"https://youtu.be/watch?v=wy42At0lQvE?234a=1":                                                           "YOUTUBE - SHORTCODE - wy42At0lQvE",
		"https://www.instagram.com/stories/msugamsingh/3124348276069682270/":                                    "INSTAGRAM - SHORTCODE - 3124348276069682270",
		"https://www.instagram.com/reel/CtbcnrxAhm7/":                                                           "INSTAGRAM - SHORTCODE - CtbcnrxAhm7",
		"https://www.instagram.com/p/Cr8LGMtB_dC/":                                                              "INSTAGRAM - SHORTCODE - Cr8LGMtB_dC",
		"https://www.instagram.com/p/Cr6EOPABupF/":                                                              "INSTAGRAM - SHORTCODE - Cr6EOPABupF",
		"https://www.youtube.com/shorts/wy42At0lQvE?asdf=1":                                                     "YOUTUBE - SHORTCODE - wy42At0lQvE",
		"https://youtube.com/shorts/wy42At0lQvE?asdf=1":                                                         "YOUTUBE - SHORTCODE - wy42At0lQvE",
		"https://m.youtube.com/watch?v=lalOy8Mbfdc":                                                             "YOUTUBE - SHORTCODE - lalOy8Mbfdc",
		"https://www.youtube.com/v/dQw4w9WgXcQ":                                                                 "YOUTUBE - SHORTCODE - dQw4w9WgXcQ",
		"https://m.youtube.com/v/dQw4w9WgXcQ":                                                                   "YOUTUBE - SHORTCODE - dQw4w9WgXcQ",
		"https://youtu.be/-wtIMTCHWuI":                                                                          "YOUTUBE - SHORTCODE - -wtIMTCHWuI",
		"https://www.youtube.com/attribution_link?a=JdfC0C9V6ZI&u=%2Fwatch%3Fv%3DEhxJLojIE_o%26feature%3Dshare": "YOUTUBE - SHORTCODE - EhxJLojIE_o",
		"https://youtube.com/oembed?url=http%3A//www.youtube.com/watch?v%3D-wtIMTCHWuI&format=json":             "YOUTUBE - SHORTCODE - -wtIMTCHWuI",
		"https://youtube.com/embed/lalOy8Mbfdc":                                                                 "YOUTUBE - SHORTCODE - lalOy8Mbfdc",
		"https://youtube.com/e/dQw4w9WgXcQ":                                                                     "YOUTUBE - SHORTCODE - dQw4w9WgXcQ",
	}

	for url, expectedOutput := range urls {
		fmt.Printf("URL: %s\n", url)
		platform, shortcode, err := ExtractPlatformPostIdentifier(url)
		if err != nil {
			fmt.Printf("Error extracting shortcode from URL: %v\n", err)
		} else {
			assert.Equal(t, expectedOutput, fmt.Sprintf("%s - SHORTCODE - %s", platform, shortcode))
			fmt.Printf("Platform: %s, Shortcode: %s\n\n", platform, shortcode)
		}
	}
}
