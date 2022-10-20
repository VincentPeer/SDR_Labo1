package server

import "fmt"

// Debuggable indicates that the implementing class has a method to check if the implementation outputs debug information.
type Debuggable interface {
	IsDebug() bool
}

// Debug will print a debug message if the implementing class' IsDebug method returns true
func Debug(source Debuggable, message string) {
	if source.IsDebug() {
		fmt.Println("DEBUG: ", message)
	}
}
