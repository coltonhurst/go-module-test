package sdk

import (
	"fmt"
	"unsafe"
)

/*
#include <stdlib.h>
typedef void* clientPtr;
extern char* run_command(const char *command, clientPtr client);
extern clientPtr init(const char *clientSettings);
extern void free_mem(clientPtr client);
*/
import "C"

type clientPointer struct {
	Pointer C.clientPtr
}

type bitwardenLibrary interface {
	init(clientSettings string) (clientPointer, error)
	freeMem(client clientPointer)
	runCommand(command string, client clientPointer) (string, error)
}

type bitwardenLibraryImpl struct{}

func newBitwardenLibrary() bitwardenLibrary {
	return &bitwardenLibraryImpl{}
}

func (b *bitwardenLibraryImpl) init(clientSettings string) (clientPointer, error) {
	ptr := C.init(C.CString(clientSettings))
	if ptr == nil {
		return clientPointer{}, fmt.Errorf("initialization failed")
	}
	return clientPointer{Pointer: ptr}, nil
}

func (b *bitwardenLibraryImpl) freeMem(client clientPointer) {
	C.free_mem(client.Pointer)
}

func (b *bitwardenLibraryImpl) runCommand(command string, client clientPointer) (string, error) {
	cstr := C.run_command(C.CString(command), client.Pointer)
	if cstr == nil {
		return "", fmt.Errorf("run command failed")
	}
	defer C.free(unsafe.Pointer(cstr))
	return C.GoString(cstr), nil
}
