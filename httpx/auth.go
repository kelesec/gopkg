package httpx

import (
	"encoding/base64"
	"fmt"
)

// BasicAuth basic认证支持
type BasicAuth struct {
	Username string
	Password string
}

func (b *BasicAuth) GetBasicAuth() (header, auth string) {
	return "Authorization", base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s:%s", b.Username, b.Password)))
}

// TODO: JWT等认证支持
