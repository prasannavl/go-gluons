package platform

import (
	"runtime"
	"syscall"
)

func SetupVirtualTerminal() {
	if runtime.GOOS != "windows" {
		return
	}
	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	getConsoleWindow := kernel32.MustFindProc("GetConsoleWindow")
	setConsoleMode := kernel32.MustFindProc("SetConsoleMode")

	r1, _, _ := syscall.Syscall(getConsoleWindow.Addr(), 0, 0, 0, 0)
	syscall.Syscall(setConsoleMode.Addr(), 2, r1, uintptr(0x0200), 0) // ENABLE_VIRTUAL_TERMINAL_INPUT
}
