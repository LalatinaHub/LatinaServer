package helper

import (
	"fmt"
	"os/exec"
)

func RestartService(names ...string) {
	for _, name := range names {
		fmt.Println("Restarting", name, "...")
		_, err := exec.Command("systemctl", "restart", name).Output()
		if err != nil {
			panic(err)
		}
	}
}

func CatchError(print bool) interface{} {
	message := recover()

	if message != nil && print {
		fmt.Println("[Error]", message)
	}
	return message
}
