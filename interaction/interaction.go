package interaction

import (
	"fmt"
)

func ParseCommandLineInputs(osArgs []string) (isHelp bool, isCreateReview bool, reviewTitle string, reviewTemplate string, isCreateConfig bool) {

	isHelp = false
	isCreateReview = false
	isCreateConfig = false
	reviewTitle = ""
	reviewTemplate = "default"
	count := len(osArgs)
	 if count > 1 {
	 	for index, arg := range osArgs {
		    	if(index == 1 && arg == "help" || arg == "h"){
		    		isHelp = true
		    	}else if(index == 1 && arg == "--create-config"){
		    		isCreateConfig = true
		    	}else if(index == 1 && arg == "create"){
		    		isCreateReview = true
		    		if(count-1>index && osArgs[index+1]!="--title"){
		    			reviewTemplate=osArgs[index+1]
		    		}
		    	}else if(isCreateReview && arg=="--title" && count-1>index){
		    		reviewTitle = osArgs[index+1]
		   	}			
		}

	} else {
		fmt.Println("Please provide a command. (type 'tongs help' for assistance)")
	}
	return isHelp, isCreateReview, reviewTitle, reviewTemplate, isCreateConfig
}
func RequestUsername() string {
	return requestStandardInput("Crucible Username")
}

func RequestPassword() string {
	fmt.Println("PASSWORD WILL NOT BE MASKED!!!")
	fmt.Println("This will only occur whenever your crucible token is forcefully expired...")
	return requestStandardInput("Crucible Password")
}
func Help() {
	fmt.Println("TODO: Write Help Text.")
	fmt.Println("For now check out the Github repo.")
	fmt.Println("tongs --create-config")
	fmt.Println("tongs create")

	
}
func ConfigFileCreated(ok bool) {
	if(ok == true){
		fmt.Println("Success!")
		fmt.Println(".tongs_config was created successfully your home directory")
		fmt.Println("edit this file with your favorite text editor to get tongs")
		fmt.Println("up and running. Check out the Github documentation for more info.")
	} else {
		fmt.Println("Error...")
		fmt.Println(".tongs_config was unable to be created in your home directory.")
		fmt.Println("This could be due to a file permissions issue, or if the config")
		fmt.Println("file already exists.")
	}
}


func requestStandardInput(prompt string) string {
	var input string
	fmt.Print(prompt, ": ")
	fmt.Scan(&input)
	return input
}