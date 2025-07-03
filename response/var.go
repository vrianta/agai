package response

var (

	// Codes holds standard HTTP Code codes
	Codes = struct {
		// 1xx - Informational
		Continue           Code
		SwitchingProtocols Code
		Processing         Code
		EarlyHints         Code

		// 2xx - Success
		OK        Code
		Created   Code
		Accepted  Code
		NoContent Code

		// 3xx - Redirection
		MovedPermanently  Code
		Found             Code
		SeeOther          Code
		NotModified       Code
		TemporaryRedirect Code
		PermanentRedirect Code

		// 4xx - Client Errors
		BadRequest       Code
		Unauthorized     Code
		Forbidden        Code
		NotFound         Code
		MethodNotAllowed Code
		Conflict         Code
		TooManyRequests  Code

		// 5xx - Server Errors
		InternalServerError Code
		NotImplemented      Code
		BadGateway          Code
		ServiceUnavailable  Code
		GatewayTimeout      Code
	}{
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
