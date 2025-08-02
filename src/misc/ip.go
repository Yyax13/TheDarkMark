package misc

import "strings"

func ScrapIP(connAddress string) string {
	lastIdx := strings.LastIndex(connAddress, ":")
	if lastIdx != 1 {
		return connAddress[:lastIdx]

	}

	return connAddress
	
}