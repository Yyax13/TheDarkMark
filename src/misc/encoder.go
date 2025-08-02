package misc

import (
	base58 "github.com/mr-tron/base58/base58"
	"encoding/base64"
	"net/url"
	"errors"
	
)

func Encode(s, method string) (v string, e error) {
	switch method {
	case "base58":
		return base58.Encode([]byte(s)), nil

	case "base64":
		return base64.StdEncoding.EncodeToString([]byte(s)), nil

	case "url2":
		return url.QueryEscape(url.QueryEscape(s)), nil

	default:
		return "", errors.New("unsupported method")

	}
}