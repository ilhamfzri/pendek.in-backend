package helper

import (
	"fmt"
	"strings"
)

func GenerateRedirectLink(domain string, username string, socialMediaName string) string {
	socialMediaUrl := SocialMediaNameToUrlFormat(socialMediaName)
	redirectLink := fmt.Sprintf("%s/%s/%s", domain, username, socialMediaUrl)
	return redirectLink
}

func SocialMediaNameToUrlFormat(socialMediaName string) string {
	url := strings.ToLower(socialMediaName)
	url = strings.Replace(url, " ", "-", -1)
	return url
}

func SocialMediaUrlToNameFormat(socialMediaUrl string) string {
	url := strings.Replace(socialMediaUrl, "-", " ", -1)
	switch url {
	case "linkedin":
		return "LinkedIn"
	case "bereal":
		return "BeReal"
	default:
		return strings.Title(url)
	}
}
