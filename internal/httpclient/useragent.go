// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package httpclient

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/mnptu/version"
)

const userAgentFormat = "mnptu/%s"
const uaEnvVar = "TF_APPEND_USER_AGENT"

// Deprecated: Use mnptuUserAgent(version) instead
func UserAgentString() string {
	ua := fmt.Sprintf(userAgentFormat, version.Version)

	if add := os.Getenv(uaEnvVar); add != "" {
		add = strings.TrimSpace(add)
		if len(add) > 0 {
			ua += " " + add
			log.Printf("[DEBUG] Using modified User-Agent: %s", ua)
		}
	}

	return ua
}

type userAgentRoundTripper struct {
	inner     http.RoundTripper
	userAgent string
}

func (rt *userAgentRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if _, ok := req.Header["User-Agent"]; !ok {
		req.Header.Set("User-Agent", rt.userAgent)
	}
	log.Printf("[TRACE] HTTP client %s request to %s", req.Method, req.URL.String())
	return rt.inner.RoundTrip(req)
}

func mnptuUserAgent(version string) string {
	ua := fmt.Sprintf("HashiCorp mnptu/%s (+https://www.mnptu.io)", version)

	if add := os.Getenv(uaEnvVar); add != "" {
		add = strings.TrimSpace(add)
		if len(add) > 0 {
			ua += " " + add
			log.Printf("[DEBUG] Using modified User-Agent: %s", ua)
		}
	}

	return ua
}
