package sharedlogger

import (
	"net/http"
)

/**
This file includes functions for formatting speciffic structures into logs as strings
*/

func getUserIpFromRequest(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

// Basic info log about an incomming http request
func (l *SharedLogger) LogHttpRequest(r *http.Request, uuid string) {
	l.Finfo("Req[%s;%s;%s;%s;%s]", getUserIpFromRequest(r), r.URL.Path, r.Method, r.Body, uuid)
}

// Basic info log about an incomming http request
func (l *SharedLogger) LogHttpResponse(code int, uuid string) {
	l.Finfo("Res[%s;%d]", uuid, code)
}
