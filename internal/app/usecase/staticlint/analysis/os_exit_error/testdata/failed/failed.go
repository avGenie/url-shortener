package failed

import "os"

func deferOSExit() {
	defer os.Exit(0) // want "os exit call error"
}

func goroutineOSExit() {
	go os.Exit(0) // want "os exit call error"
}

func exitOSExit() {
	os.Exit(0) // want "os exit call error"
}
