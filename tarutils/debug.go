// +build debug

package tarutils

import "log"

func debug(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}
