package helper

import (
	"fmt"
	"os/exec"
)

func ReloadService(names ...string) {
	for _, name := range names {
		fmt.Println("Reloading", name, "...")
		_, err := exec.Command("systemctl", "reload", name).Output()
		if err != nil {
			if err.Error() == "exit status 1" {
				exec.Command("systemctl", "restart", name)
			} else {
				panic(err)
			}
		}

		fmt.Println(name, "successfully reloaded !")
	}
}

func CatchError(print bool) any {
	message := recover()

	if message != nil && print {
		fmt.Println("[Error]", message)
	}
	return message
}
