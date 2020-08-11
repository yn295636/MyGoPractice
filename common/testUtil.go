package common

import "flag"

func IsRunningUnderGoTest() bool {
	if flag.Lookup("test.v") == nil {
		return false
	}
	return true
}
