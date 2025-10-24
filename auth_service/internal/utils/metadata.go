package utils

import (
	"fmt"
	"time"
)

// BuildCookieHeader generates a Set-Cookie string
func BuildCookieHeader(name, value string, maxAge time.Duration, httpOnly, secure bool) string {
	cookie := fmt.Sprintf("%s=%s; Path=/; Max-Age=%d", name, value, int(maxAge.Seconds()))
	if httpOnly {
		cookie += "; HttpOnly"
	}
	if secure {
		cookie += "; Secure"
	}
	cookie += "; SameSite=Strict"
	return cookie
}
