package winapi

import (
	"errors"
	"fmt"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

const (
	FORMAT_MESSAGE_IGNORE_INSERTS = 0x00000200
	FORMAT_MESSAGE_FROM_STRING    = 0x00000400
	FORMAT_MESSAGE_FROM_HMODULE   = 0x00000800
	FORMAT_MESSAGE_FROM_SYSTEM    = 0x00001000
	FORMAT_MESSAGE_ARGUMENT_ARRAY = 0x00002000
	FORMAT_MESSAGE_MAX_WIDTH_MASK = 0x000000FF
)

func FormatMessage(flags uint32, msgsrc interface{}, msgid uint32, langid uint32, args *byte) (string, error) {
	var b [300]uint16
	n, err := _FormatMessage(flags, msgsrc, msgid, langid, &b[0], 300, args)
	if err != nil {
		return "", err
	}
	for ; n > 0 && (b[n-1] == '\n' || b[n-1] == '\r'); n-- {
	}
	return string(utf16.Decode(b[:n])), nil
}

func _FormatMessage(flags uint32, msgsrc interface{}, msgid uint32, langid uint32, buf *uint16, nSize uint32, args *byte) (n uint32, err error) {
	r0, _, e1 := syscall.Syscall9(procFormatMessage.Addr(), 7,
		uintptr(flags), uintptr(0), uintptr(msgid), uintptr(langid),
		uintptr(unsafe.Pointer(buf)), uintptr(nSize),
		uintptr(unsafe.Pointer(args)), 0, 0)
	n = uint32(r0)
	if n == 0 {
		err = fmt.Errorf("winapi._FormatMessage error: %d", uint32(e1))
	}
	return
}

/*
typedef struct _SECURITY_ATTRIBUTES {
    DWORD  nLength;
    void   *pSecurityDescriptor;
    BOOL   bInheritHandle;
} SECURITY_ATTRIBUTES;
*/
type SECURITY_ATTRIBUTES struct {
	Length             uint32
	SecurityDescriptor uintptr
	InheritHandle      int32
}

func CreateNamedPipe(
	name string,
	openMode uint32,
	pipeMode uint32,
	maxInstances uint32,
	outBufferSize uint32,
	inBufferSize uint32,
	defaultTimeOut uint32,
	sa *SECURITY_ATTRIBUTES) (h HANDLE, err error) {
	pName, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return
	}

	h, err = _CreateNamedPipe(pName, openMode, pipeMode, maxInstances,
		outBufferSize, inBufferSize, defaultTimeOut, sa)
	return
}

/*
HANDLE WINAPI CreateNamedPipe(
  _In_     LPCTSTR               lpName,
  _In_     DWORD                 dwOpenMode,
  _In_     DWORD                 dwPipeMode,
  _In_     DWORD                 nMaxInstances,
  _In_     DWORD                 nOutBufferSize,
  _In_     DWORD                 nInBufferSize,
  _In_     DWORD                 nDefaultTimeOut,
  _In_opt_ LPSECURITY_ATTRIBUTES lpSecurityAttributes
);
*/
func _CreateNamedPipe(pName *uint16, dwOpenMode uint32, dwPipeMode uint32,
	nMaxInstances uint32, nOutBufferSize uint32, nInBufferSize uint32,
	nDefaultTimeOut uint32, pSecurityAttributes *SECURITY_ATTRIBUTES) (h HANDLE, err error) {
	r1, _, e1 := syscall.Syscall9(procCreateNamedPipe.Addr(), 8,
		uintptr(unsafe.Pointer(pName)),
		uintptr(dwOpenMode),
		uintptr(dwPipeMode),
		uintptr(nMaxInstances),
		uintptr(nOutBufferSize),
		uintptr(nInBufferSize),
		uintptr(nDefaultTimeOut),
		uintptr(unsafe.Pointer(pSecurityAttributes)),
		0)
	if h == INVALID_HANDLE_VALUE {
		wec := WinErrorCode(e1)
		if wec != 0 {
			err = wec
		} else {
			err = errors.New("GetModuleHandle failed.")
		}
	} else {
		h = HANDLE(r1)
	}
	return
}
