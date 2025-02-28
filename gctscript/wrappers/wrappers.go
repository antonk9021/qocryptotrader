package wrappers

import (
	"github.com/antonk9021/qocryptotrader/gctscript/modules"
	"github.com/antonk9021/qocryptotrader/gctscript/wrappers/validator"
)

// GetWrapper returns the instance of each wrapper to use
func GetWrapper() modules.GCTExchange {
	if validator.IsTestExecution.Load() == true {
		return validator.Wrapper{}
	}
	return modules.Wrapper
}
