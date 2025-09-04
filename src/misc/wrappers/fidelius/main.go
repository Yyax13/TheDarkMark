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
	"github.com/Yyax13/onTop-C2/src/misc"
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
			misc.PanicWarn(fmt.Sprintf("Invalid JSON parameters for '%s': %v\n\n", fideliusNameFromC, err), false)
			return 0

		}

	}

	creator, exists := fidelius.AvaliableFideliusCreators[fideliusNameFromC]
	if !exists {
		misc.PanicWarn(fmt.Sprintf("Not found the encoder %s\n\n", fideliusNameFromC), false)
		return 0

	}

	fideliusCasting, err := creator(params)
	if err != nil {
		misc.PanicWarn(fmt.Sprintf("Failed to create encoder %s with params %v: %v", fideliusNameFromC, params, err), false)
		return 0

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
func DestroyEncoder(id C.int) {
	fideliusMutex.Lock()
	defer fideliusMutex.Unlock()

	if _, exists := fideliusMap[int(id)]; exists {
		delete(fideliusMap, int(id))

	} else {
		misc.SysLog("Attempt to delete a non-existent fidelius\n\n", false)

	}
	
}

//export Encode
func Encode(encoderID C.int, data *C.char, dataLen C.int, out **C.char, outLen *C.int) {
	fideliusMutex.Lock()
	fidelius, exists := fideliusMap[int(encoderID)]
	if !exists {
		*out = nil
		*outLen = 0
		return

	}

	fideliusMutex.Unlock()

	dataFromC := C.GoBytes(unsafe.Pointer(data), dataLen)
	encoded, err := fidelius.Fidelius.Encode(dataFromC)
	if err != nil {
		*out = nil
		*outLen = 0
		return

	}

	cEncoded := C.malloc(C.size_t(len(encoded)))
	if cEncoded == nil {
		misc.SysLog("Failed to allocate memory for encoded data\n\n", false)
		*out = nil
		*outLen = 0
		return

	}

	C.memcpy(cEncoded, unsafe.Pointer(&encoded[0]), C.size_t(len(encoded)))
	*out = (*C.char)(cEncoded)
	*outLen = C.int(len(encoded))

}

//export Decode
func Decode(encoderID C.int, data *C.char, dataLen C.int, out **C.char, outLen *C.int) {
	fideliusMutex.Lock()
	fidelius, exists := fideliusMap[int(encoderID)]
	if !exists {
		*out = nil
		*outLen = 0
		return

	}

	fideliusMutex.Unlock()

	dataFromC := C.GoBytes(unsafe.Pointer(data), dataLen)
	decoded, err := fidelius.Fidelius.Decode(dataFromC)
	if err != nil {
		*out = nil
		*outLen = 0
		return

	}

	cDecoded := C.malloc(C.size_t(len(decoded)))
	if cDecoded == nil {
		misc.SysLog("Failed to allocate memory for decoded data\n\n", false)
		*out = nil
		*outLen = 0
		return

	}

	C.memcpy(cDecoded, unsafe.Pointer(&decoded[0]), C.size_t(len(decoded)))
	*out = (*C.char)(cDecoded)
	*outLen = C.int(len(decoded))
}

//export FreeGoMem
func FreeGoMem(pointer *C.char) {
	if pointer == nil {
		misc.SysLog("Attempt to free a null pointer\n\n", false)
		return

	}

	C.free(unsafe.Pointer(pointer))

}

func main() {}