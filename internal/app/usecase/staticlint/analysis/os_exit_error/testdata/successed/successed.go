package successed

import (
	"fmt"
)

func withoutExitOSExit() {
	fmt.Println("without exit")
}

func errCheckFunc() {
	withoutExitOSExit() // ""
}
