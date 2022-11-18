package helper

import "github.com/ilhamfzri/pendek.in/internal/model/web"

func ToWebResponseFailed(err error) web.WebResponseFailed {
	return web.WebResponseFailed{
		Status:  "failed",
		Message: err.Error(),
	}
}
