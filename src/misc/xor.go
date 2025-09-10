package misc

func Xor(plain, key []byte) ([]byte) { // better than RSA man, this is very nice 'cause just sniffing into traffic will not help the blue team :D
	keyCopy := key
	xorResult := make([]byte, len(plain))
	for i := range len(plain) {
		xorResult[i] = plain[i] ^ keyCopy[(i+1)%len(keyCopy)]
		keyCopy[(i+2)%len(keyCopy)] = keyCopy[i%len(keyCopy)]

	}

	return xorResult
	
}