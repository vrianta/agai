package agai

import "net/http"

// Redirect Response codes

type (
	code int

	httpStatus struct {
		// 1xx - Informational
		Continue           code
		SwitchingProtocols code
		Processing         code
		EarlyHints         code

		// 2xx - Success
		OK        code
		Created   code
		Accepted  code
		NoContent code

		// 3xx - Redirection
		MovedPermanently  code
		Found             code
		SeeOther          code
		NotModified       code
		TemporaryRedirect code
		PermanentRedirect code

		// 4xx - Client Errors
		BadRequest       code
		Unauthorized     code
		Forbidden        code
		NotFound         code
		MethodNotAllowed code
		Conflict         code
		TooManyRequests  code

		// 5xx - Server Errors
		InternalServerError code
		NotImplemented      code
		BadGateway          code
		ServiceUnavailable  code
		GatewayTimeout      code
	}
)

var (
	HttpStatus = httpStatus{
		// 1xx
		Continue:           100,
		SwitchingProtocols: 101,
		Processing:         102,
		EarlyHints:         103,

		// 2xx
		OK:        200,
		Created:   201,
		Accepted:  202,
		NoContent: 204,

		// 3xx
		MovedPermanently:  301,
		Found:             302,
		SeeOther:          303,
		NotModified:       304,
		TemporaryRedirect: 307,
		PermanentRedirect: 308,

		// 4xx
		BadRequest:       400,
		Unauthorized:     401,
		Forbidden:        403,
		NotFound:         404,
		MethodNotAllowed: 405,
		Conflict:         409,
		TooManyRequests:  429,

		// 5xx
		InternalServerError: 500,
		NotImplemented:      501,
		BadGateway:          502,
		ServiceUnavailable:  503,
		GatewayTimeout:      504,
	}
)

// Redirects to the URI user provided
func (_c *Controller) Redirect(uri string) View {
	http.Redirect(_c.W, _c.R, uri, int(HttpStatus.SeeOther)) // if golang developpers worked so hard to create this why should I do it again :P
	return nil
}

func (_c *Controller) RedirectWithCode(uri string, _code code) View {
	http.Redirect(_c.W, _c.R, uri, int(_code))
	return nil
}
