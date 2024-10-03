package basic_api

type Method string

const (
	GET     Method = "GET"
	HEAD    Method = "HEAD"
	POST    Method = "POST"
	PUT     Method = "PUT"
	PATCH   Method = "PATCH"
	DELETE  Method = "DELETE"
	CONNECT Method = "CONNECT"
	OPTIONS Method = "OPTIONS"
	TRACE   Method = "TRACE"
)
