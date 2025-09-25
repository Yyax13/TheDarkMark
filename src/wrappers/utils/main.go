package main

/*
#include <stdlib.h>
typedef struct {
	unsigned char	*data;
	int				len;

} _go_base64;

*/
import "C"

import (
	"encoding/base64"

)

//export _b64_d
func _b64_d(data *C.char) (C._go_base64) {
	goData := C.GoString(data)
	var decoded C._go_base64;
	decodedBytes, err := base64.StdEncoding.DecodeString(goData)
	if err != nil {
		return decoded

	}

	decoded.data = (*C.uchar)(C.CBytes(decodedBytes))
	decoded.len = C.int(len(decodedBytes))

	return decoded

}

func main() {}