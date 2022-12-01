package helper

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func GenerateRedirectLink(host string, username string, socialMediaName string) string {
	socialMediaUrl := SocialMediaNameToUrlFormat(socialMediaName)
	redirectLink := fmt.Sprintf("%s/%s/%s", host, username, socialMediaUrl)
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

func SocialMediaValidator(socialMediaName string, link_or_username string) bool {
	var err error
	validate = validator.New()
	switch socialMediaName {
	case "Tiktok", "Twitter", "Instagram":
		err = validate.Var(link_or_username, "alphanum")
	case "Whatsapp":
		err = validate.Var(link_or_username, "e164")
	default:
		err = validate.Var(link_or_username, "url")
	}
	return err == nil

}

func GenerateLinkResponse(socialMediaName string, link_or_username string) string {
	var linkResponse string
	switch socialMediaName {
	case "Tiktok", "Twitter", "Instagram":
		linkResponse = fmt.Sprintf("https://www.%s.com/%s", strings.ToLower(socialMediaName), link_or_username)
	case "Whatsapp":
		linkResponse = fmt.Sprintf("https://wa.me/%s", link_or_username)
	default:
		linkResponse = link_or_username
	}
	return linkResponse
}

func GetCustomThumbnailUrl(domain string, imageID string) string {
	thumbnailUrl := fmt.Sprintf("%s/%s/%s.jpg", domain, thumbnailResourceEndpointPath, imageID)
	return thumbnailUrl
}
