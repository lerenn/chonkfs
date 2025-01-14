package chonkfs

import (
	"fmt"
)

var debugActivated = false

// SetDebug sets the debug mode
func SetDebug(d bool) {
	debugActivated = d
}

// IsDebug returns the debug mode
func IsDebug() bool {
	return debugActivated
}

func debugf(format string, args ...interface{}) {
	if debugActivated {
		fmt.Printf("CALLED: "+format, args...)
	}
}
