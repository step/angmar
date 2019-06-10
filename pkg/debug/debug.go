// +build debug

package tarutils

import "log"

func Debug(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}
