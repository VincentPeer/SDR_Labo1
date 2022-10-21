package server

import "fmt"

// debuggable indicates that the implementing class has a method to check if the implementation outputs debug information.
type debuggable interface {
	isDebug() bool
}

// debug will print a debug message if the implementing class' IsDebug method returns true
func debug(source debuggable, message string) {
	if source.isDebug() {
		fmt.Println("DEBUG: ", message)
	}
}
