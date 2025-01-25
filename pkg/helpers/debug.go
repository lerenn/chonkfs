package helpers

import (
	"log"
	"os"
)

// DebugOpenFlags prints the open flags on the logger.
func DebugOpenFlags(logger *log.Logger, flags uint32) {
	if flags&uint32(os.O_RDONLY) != 0 {
		logger.Println("O_RDONLY")
	}
	if flags&uint32(os.O_WRONLY) != 0 {
		logger.Println("O_WRONLY")
	}
	if flags&uint32(os.O_RDWR) != 0 {
		logger.Println("O_RDWR")
	}
	if flags&uint32(os.O_APPEND) != 0 {
		logger.Println("O_APPEND")
	}
	if flags&uint32(os.O_CREATE) != 0 {
		logger.Println("O_CREATE")
	}
	if flags&uint32(os.O_EXCL) != 0 {
		logger.Println("O_EXCL")
	}
	if flags&uint32(os.O_SYNC) != 0 {
		logger.Println("O_SYNC")
	}
	if flags&uint32(os.O_TRUNC) != 0 {
		logger.Println("O_TRUNC")
	}
}
