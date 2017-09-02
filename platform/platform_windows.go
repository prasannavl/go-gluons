package platform

import (
	win "github.com/prasannavl/go-grab/platform"
)

func Init() {
	win.SetupVirtualTerminal()
}
