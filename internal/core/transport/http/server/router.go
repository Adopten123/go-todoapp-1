package core_http_server

import "net/http"

type ApiVersion string

var (
	ApiVersion1 ApiVersion = "v1"
	ApiVersion2 ApiVersion = "v2"
	ApiVersion3 ApiVersion = "v3"
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion ApiVersion
}
