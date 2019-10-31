package main

import "os/exec"

func SwitchToCharacter(name string) {
	cmd := "/Users/stephane/go/src/dofus-middleware" + name + ".sh"
	out := exec.Command("/bin/bash", cmd)
	_ = out.Run()
}
