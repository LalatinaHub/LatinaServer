package helper

import (
	"fmt"
	"os/exec"
)

func ReloadService(names ...string) {
	for _, name := range names {
		fmt.Println("Restarting", name, "...")
		_, err := exec.Command("systemctl", "restart", name).Output()
		if err != nil {
			fmt.Println("Failed restarting", name)
			fmt.Println("[", name, "]", err.Error())
			continue
		}

		fmt.Println(name, "successfully restarted !")
	}
}

func CatchError(print bool) interface{} {
	message := recover()

	if message != nil && print {
		fmt.Println("[Error]", message)
	}
	return message
}
