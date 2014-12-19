package interaction

import (
	"fmt"
	"strings"
)

func ParseCommandLineInputs(osArgs []string) (isHelp bool,
	isCreateReview bool, isUpdateReview bool, isGetToken bool,
	reviewTemplate string, isCreateConfig bool, projectId string,
	isListTemplates bool) {

	isHelp = false
	isCreateReview = false
	isCreateConfig = false
	isGetToken = false
	isUpdateReview = false
	isListTemplates = false
	reviewTemplate = "default"
	projectId = ""

	count := len(osArgs)
	if count > 1 {
		for index, arg := range osArgs {
			if index == 1 {
				if arg == "help" || arg == "h" {
					isHelp = true
				} else if arg == "setup" {
					isCreateConfig = true
				} else if arg == "token" {
					isGetToken = true
				} else if arg == "templates" {
					isListTemplates = true
					if count-1 > index {
						reviewTemplate = osArgs[index+1]
					}
				} else if arg == "update" {
					isUpdateReview = true
					if count-2 > index {
						reviewTemplate = osArgs[index+1]
						projectId = osArgs[index+2]
					}
				} else if arg == "create" {
					isCreateReview = true
					if count-1 > index {
						reviewTemplate = osArgs[index+1]
					}
				}
			}
		}

	} else {
		isHelp = true
	}
	return isHelp, isCreateReview, isUpdateReview, isGetToken,
		reviewTemplate, isCreateConfig, projectId, isListTemplates
}

func RequestUrl() string {
	return strings.Trim(requestStandardInput("Crucible Url"), " \\/")
}

func RequestUsername() string {
	return strings.Trim(requestStandardInput("Crucible Username"), " ")
}

func RequestPassword() string {
	fmt.Println("*************************************")
	fmt.Println("Password input is not masked!")
	fmt.Println("Clear this terminal after entering")
	fmt.Println("your password for safety.")
	fmt.Println("*************************************")
	return requestStandardInput("Crucible Password")
}
func Help() {
	fmt.Println("Usage: tongs [OPTION] [TEMPLATE] [REVIEW-ID]")
	fmt.Println("")
	fmt.Println("Utility for creating code reviews quickly")
	fmt.Println("in Crucible based on predefined templates.")
	fmt.Println("")
	fmt.Println("Options: ")
	fmt.Println("setup                           ")
	fmt.Println("token                           ")
	fmt.Println("templates                       ")
	fmt.Println("templates [TEMPLATE]            ")
	fmt.Println("create [TEMPLATE]               ")
	fmt.Println("update [TEMPLATE] [REVIEW-ID]   ")

}
func ConfigFileCreated(ok bool) {
	if ok == true {
		fmt.Println("tongs.cfg created!")
	} else {
		fmt.Println("tongs.cfg already exists.")
	}
}

func requestStandardInput(prompt string) string {
	var input string
	fmt.Print(prompt, ": ")
	fmt.Scan(&input)
	return input
}
