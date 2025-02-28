package gctscript

import (
	"github.com/antonk9021/qocryptotrader/gctscript/modules"
	"github.com/antonk9021/qocryptotrader/gctscript/wrappers/gct"
)

// Setup configures the wrapper interface to use
func Setup() {
	modules.SetModuleWrapper(gct.Setup())
}
