package windowmanagement

import (
	"fmt"
	"os/exec"
)

func SwitchToCharacter(name string) {
	fmt.Println("Window.SwithToCharacter:" + name)
	cmd := "/Users/stephane/go/src/dofusmiddleware/windowmanagement/focus_window.sh"
	out := exec.Command("/bin/bash", cmd, name)
	_ = out.Run()
}
