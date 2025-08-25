package misc

func Xor(plain, key []byte) ([]byte) {
	keyCopy := key
	xorResult := make([]byte, len(plain))
	for i := range len(plain) {
		xorResult[i] = plain[i] ^ keyCopy[(i+1)%len(keyCopy)]
		keyCopy[(i+2)%len(keyCopy)] = keyCopy[i%len(keyCopy)]

	}

	return xorResult
	
}