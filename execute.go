package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// run command in background and Send result to telegram
func RunCommand(messenger *TelegramMessenger, command []string) {

	// cerate response tmplate
	returnMessageTemplate := fmt.Sprintf("COMMAND: \n%s\n\nOUTPUT:\n\n", strings.Join(command, " ")) + "%s"

	// change dir
	if command[0] == "cd" {
		path := strings.Join(command[1:], "")
		if err := os.Chdir(path); err != nil {
			returnMessage := fmt.Sprintf(returnMessageTemplate, err.Error())
			messenger.Send(returnMessage)
			return
		}
		returnMessage := fmt.Sprintf(returnMessageTemplate, "")
		messenger.Send(returnMessage)
		return
	}

	doneFlag := false
	pDoneFlag := &doneFlag

	// run command in background
	go func() {
		out, err := exec.Command(command[0], command[1:]...).Output()
		var returnMessage string

		if err != nil {
			returnMessage = fmt.Sprintf(returnMessageTemplate, err.Error())
		} else {
			returnMessage = fmt.Sprintf(returnMessageTemplate, string(out[:]))
		}
		messenger.Send(returnMessage)
		*pDoneFlag = true
	}()

	// if command done in 4 second return response:
	for i := 0; i < 4; i++ {
		time.Sleep(time.Second)
		if doneFlag {
			return
		}
	}
	// else:
	returnMessage := fmt.Sprintf(returnMessageTemplate,
		"Your command is running in the background and you will get the results when it is done.")
	messenger.Send(returnMessage)
}
