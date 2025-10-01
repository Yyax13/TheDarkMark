package main

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"sync"
	"unsafe"

	"github.com/Yyax13/onTop-C2/src/fidelius"
	"github.com/Yyax13/onTop-C2/src/types"
)

var (
	fideliusMap = 		make(map[int]*types.Fidelius)
	fideliusIDCounter 	int
	fideliusMutex		sync.Mutex

)

//export CreateEncoder
func CreateEncoder(fideliusName *C.char, fideliusNameLen C.int, paramsJson *C.char, paramsJsonLen C.int) C.int {
	fideliusMutex.Lock()
	defer fideliusMutex.Unlock()

	fideliusNameFromC := C.GoStringN(fideliusName, fideliusNameLen)
	paramsJsonFromC := C.GoStringN(paramsJson, paramsJsonLen)

	params := make(map[string]string)
	if paramsJsonFromC != "" {
		if err := json.Unmarshal([]byte(paramsJsonFromC), &params); err != nil {
			return 0 // Invalid JSON parameters

		}

	}

	creator, exists := fidelius.AvaliableFideliusCreators[fideliusNameFromC]
	if !exists {
		return 0 // Not found the encoder

	}

	fideliusCasting, err := creator(params)
	if err != nil {
		return 0 // Failed to create the encoder with these params

	}

	fideliusIDCounter++
	id := fideliusIDCounter
	fideliusMap[id] = &types.Fidelius{
		Name: fideliusNameFromC,
		Description: fmt.Sprintf("Wrapper for %s fidelius", fideliusNameFromC),
		Fidelius: fideliusCasting,

	}

	return C.int(id)

}

//export DestroyEncoder
func DestroyEncoder(id C.int) C.int {
	fideliusMutex.Lock()
	defer fideliusMutex.Unlock()

	if _, exists := fideliusMap[int(id)]; exists {
		delete(fideliusMap, int(id))
		return 1

	} else {
		return 0 // Attempt to delete a non-existent fidelius

	}
	
}

//export Encode
func Encode(encoderID C.int, data *C.uchar, dataLen C.int, out **C.uchar, outLen *C.int) C.int {
	fideliusMutex.Lock()
	fidelius, exists := fideliusMap[int(encoderID)]
	if !exists {
		*out = nil
		*outLen = 0
		return 0 // Fidelius was not found

	}

	fideliusMutex.Unlock()

	dataFromC := C.GoBytes(unsafe.Pointer(data), dataLen)
	encoded, err := fidelius.Fidelius.Encode(dataFromC)
	if err != nil {
		*out = nil
		*outLen = 0
		return 0 // Some error occurred in Encode

	}

	cEncoded := C.malloc(C.size_t(len(encoded)))
	if cEncoded == nil {
		*out = nil
		*outLen = 0
		return 0 // Failed to allocate memory for encoded data

	}

	C.memcpy(cEncoded, unsafe.Pointer(&encoded[0]), C.size_t(len(encoded)))
	*out = (*C.uchar)(cEncoded)
	*outLen = C.int(len(encoded))

	return 1

}

//export Decode
func Decode(encoderID C.int, data *C.uchar, dataLen C.int, out **C.uchar, outLen *C.int) C.int {
	fideliusMutex.Lock()
	fidelius, exists := fideliusMap[int(encoderID)]
	if !exists {
		*out = nil
		*outLen = 0
		return 0 // Fidelius not found

	}

	fideliusMutex.Unlock()

	dataFromC := C.GoBytes(unsafe.Pointer(data), dataLen)
	decoded, err := fidelius.Fidelius.Decode(dataFromC)
	if err != nil {
		*out = nil
		*outLen = 0
		return 0 // Failed to allocate memory for decoded data

	}

	cDecoded := C.malloc(C.size_t(len(decoded)))
	if cDecoded == nil {
		*out = nil
		*outLen = 0
		return 0 // Some error occurred in mem allocation (aka malloc)

	}

	C.memcpy(cDecoded, unsafe.Pointer(&decoded[0]), C.size_t(len(decoded)))
	*out = (*C.uchar)(cDecoded)
	*outLen = C.int(len(decoded))

	return 1

}

func main() {}