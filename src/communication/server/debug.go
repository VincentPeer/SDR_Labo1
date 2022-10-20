package server

import "fmt"

type Debuggable interface {
	IsDebug() bool
}

func Debug(source Debuggable, message string) {
	if source.IsDebug() {
		fmt.Println("DEBUG: ", message)
	}
}
