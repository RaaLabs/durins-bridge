package proxy

import (
	"net/http"
	"regexp"

	"golang.org/x/exp/slog"
)

type OverrideHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) bool
}

type OverridesHandler struct {
	Proxy     http.Handler
	Overrides []OverrideHandler
}

func (oh *OverridesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := oh.requestWithoutVersionPrefix(r)

	for _, o := range oh.Overrides {
		if o.ServeHTTP(w, req) {
			slog.Debug("Request handled by override", "url", r.URL)
			return
		}
	}

	slog.Debug("Handing over request to proxy", "url", r.URL)
	oh.Proxy.ServeHTTP(w, r)
}

var versionPrefixPattern = regexp.MustCompile(`^/v[^/]*`)

func (*OverridesHandler) requestWithoutVersionPrefix(r *http.Request) *http.Request {
	match := versionPrefixPattern.FindStringIndex(r.URL.Path)
	if len(match) == 2 {
		req := r.Clone(r.Context())
		req.URL.Path = req.URL.Path[match[1]:]
		//req.URL.RawPath = req.URL.RawPath[match[1]:]

		return req
	}

	return r
}
