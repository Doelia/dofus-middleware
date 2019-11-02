package dofusmiddleware

import (
	"bufio"
	"fmt"
	"os"
)

func InputKeyboard() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter your command: ")
		command, _ := reader.ReadString('\n')
		fmt.Print("Command " + command)

	}
}
