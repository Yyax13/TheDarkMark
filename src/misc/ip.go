package misc

import "strings"

func ScrapIP(connAddress string) string { // this is useless, used just 1 time and we made a entire file in misc just for that xd
	lastIdx := strings.LastIndex(connAddress, ":")
	if lastIdx != 1 {
		return connAddress[:lastIdx]

	}

	return connAddress
	
}