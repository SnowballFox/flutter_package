package flutter_package

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"time"
)

func AskForConfirmation(msg string) bool {
	fmt.Println(msg)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	notOkayResponses := []string{"n", "N", "no", "No", "NO"}
	if ContainsString(okayResponses, response) {
		return true
	} else if ContainsString(notOkayResponses, response) {
		return false
	} else {
		fmt.Println("请输入y或者n")
		return AskForConfirmation(msg)
	}
}

func AskForInformation(msg string) string {
	fmt.Println(msg)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	return response
}

func printCommand(dir string, name string, arg ...string) {
	fmt.Printf("[Current Directory]: %s", dir)
	fmt.Printf("[Command]: %s", name)
	for _, ele := range arg {
		fmt.Print(" " + ele)
	}
	fmt.Println()
}

func ExecuteCommand(dir string, name string, arg ...string) {
	printCommand(dir, name, arg...)

	count := 1
	status := true
	go func() {
		for status {
			fmt.Printf("\r%ds", count)
			time.Sleep(1000 * time.Millisecond)
			count += 1
		}
		fmt.Println()
	}()

	defer func() {
		count = 1
		status = false
	}()

	command := exec.Command(name, arg...)
	command.Dir = dir
	outInfo := bytes.Buffer{}
	command.Stdout = &outInfo
	err := command.Start()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = command.Wait()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println()
	fmt.Println(outInfo.String())
}
