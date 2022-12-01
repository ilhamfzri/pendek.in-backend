package helper

import (
	"github.com/ilhamfzri/pendek.in/helper/uaparser"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
)

func SocialMediaInteractionsToDeviceAnalytic(socialMediaInteractions *[]domain.SocialMediaInteraction) domain.DeviceAnalytic {
	deviceAnalytic := domain.DeviceAnalytic{}
	for _, socialMediaInteraction := range *socialMediaInteractions {
		userAgent := socialMediaInteraction.UserAgent
		ua := uaparser.Parse(userAgent)

		if ua.Desktop {
			deviceAnalytic.Desktop += 1
		} else if ua.Mobile {
			deviceAnalytic.Mobile += 1
		} else if ua.Tablet {
			deviceAnalytic.Tablet += 1
		} else {
			deviceAnalytic.Other += 1
		}
	}
	return deviceAnalytic
}

func CustomLinkInteractionsToDeviceAnalytic(customLinkInteractions *[]domain.CustomLinkInteraction) domain.DeviceAnalytic {
	deviceAnalytic := domain.DeviceAnalytic{}
	for _, customLinkInteraction := range *customLinkInteractions {
		userAgent := customLinkInteraction.UserAgent
		ua := uaparser.Parse(userAgent)

		if ua.Desktop {
			deviceAnalytic.Desktop += 1
		} else if ua.Mobile {
			deviceAnalytic.Mobile += 1
		} else if ua.Tablet {
			deviceAnalytic.Tablet += 1
		} else {
			deviceAnalytic.Other += 1
		}
	}
	return deviceAnalytic
}
