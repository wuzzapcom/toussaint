package app

import (
	"fmt"
	"net/http"
	"net/url"
)

type httpError struct {
	Message string
	Code    int
}

func parseQueryParameter(params url.Values, field string) (string, *httpError) {
	clientTypeStr := params.Get(field)
	if clientTypeStr == "" {
		return "", &httpError{
			Message: fmt.Sprintf("field %s not found in query parameters", field),
			Code:    http.StatusNotAcceptable,
		}
	}

	clientTypeStr, err := url.QueryUnescape(clientTypeStr)
	if err != nil {
		return "", &httpError{
			Message: fmt.Sprintf("failed to unescape %s query parameter: %+v", field, err),
			Code:    http.StatusBadRequest,
		}
	}
	return clientTypeStr, nil
}

func parseClientType(params url.Values) (ClientType, *httpError) {

	clientTypeStr, httpErr := parseQueryParameter(params, "client-type")
	if httpErr != nil {
		return 0, httpErr
	}

	clientType, err := GetClientType(clientTypeStr)
	if err != nil {
		return 0, &httpError{
			Message: fmt.Sprintf("failed to parse field client-type as ClientType: %+v", err),
			Code:    http.StatusNotAcceptable,
		}
	}
	return clientType, nil
}

func parseRequestType(params url.Values) (RequestType, *httpError) {
	requestTypeStr, httpErr := parseQueryParameter(params, "request-type")
	if httpErr != nil {
		return 0, httpErr
	}

	requestType, err := GetRequestType(requestTypeStr)
	if err != nil {
		return 0, &httpError{
			Message: fmt.Sprintf("failed to parse field request-type as RequestType: %+v", err),
			Code:    http.StatusNotAcceptable,
		}
	}
	return requestType, nil
}
