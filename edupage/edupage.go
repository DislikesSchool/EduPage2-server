package edupage

import "net/http"

// Handle is used to access the Edupage API.
type Handle struct {
	hc     *http.Client
	server string
}
