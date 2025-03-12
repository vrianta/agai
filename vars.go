package server

// Global instance of the server
var (
	EmptyTemplateData = struct{}{} // when you do not want to send any data to the template
	srvInstance       *server
	templateRecords   = make(map[string]templates) // keep the reocrd of all the templated which are already templated

	config = Config{
		Http:          false,
		Static_folder: "Static",
		Views_folder:  "Views",
	}

	// ResponseCodes holds standard HTTP ResponseCode codes
	ResponseCodes = struct {
		// 1xx - Informational
		Continue           ResponseCode
		SwitchingProtocols ResponseCode
		Processing         ResponseCode
		EarlyHints         ResponseCode

		// 2xx - Success
		OK        ResponseCode
		Created   ResponseCode
		Accepted  ResponseCode
		NoContent ResponseCode

		// 3xx - Redirection
		MovedPermanently  ResponseCode
		Found             ResponseCode
		SeeOther          ResponseCode
		NotModified       ResponseCode
		TemporaryRedirect ResponseCode
		PermanentRedirect ResponseCode

		// 4xx - Client Errors
		BadRequest       ResponseCode
		Unauthorized     ResponseCode
		Forbidden        ResponseCode
		NotFound         ResponseCode
		MethodNotAllowed ResponseCode
		Conflict         ResponseCode
		TooManyRequests  ResponseCode

		// 5xx - Server Errors
		InternalServerError ResponseCode
		NotImplemented      ResponseCode
		BadGateway          ResponseCode
		ServiceUnavailable  ResponseCode
		GatewayTimeout      ResponseCode
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
