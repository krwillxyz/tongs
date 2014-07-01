package main

import (
	"./config/"
	"./crucible/"
	"./interaction/"
	"os"
	"strings"
	"fmt"
)

func main() {

	isHelp, isCreateReview, isUpdateReview, isGetToken, reviewTitle, reviewTemplate, isCreateConfig, projectId := interaction.ParseCommandLineInputs(os.Args)

	if isHelp {
		interaction.Help()
	} else if isCreateConfig {
		ok := config.CreateConfigFile() 
		interaction.ConfigFileCreated(ok)
	} else if isGetToken {
		baseUrl := config.LoadBaseUrl()
		fmt.Println("Crucible URL: ", baseUrl)
		username := interaction.RequestUsername()
		username = strings.ToLower(username)
		password := interaction.RequestPassword()
		token := crucible.Login(username, password, baseUrl)
		fmt.Println("Recieved Token: ",token)
		config.SaveUsername(username)
		config.SaveToken(token)
		fmt.Println("Token Saved Successfully!!!")
		fmt.Println("(you can now clear your terminal)")
	} else if isCreateReview {
		if reviewTitle == "" {
			reviewTitle = config.LoadReviewTitle(reviewTemplate)
		}

		baseUrl := config.LoadBaseUrl()
		projectKey := config.LoadProjectKey(reviewTemplate)
		reviewLength := config.LoadDuration(reviewTemplate)
		reviewers := config.LoadUserIds(reviewTemplate)
		hasUsername, username := config.LoadUsername()
		hasToken, token := config.LoadToken()

		if hasUsername && hasToken {
		
			if ok, _ := crucible.CreateReview(reviewTitle, reviewTemplate, reviewLength, username, baseUrl, token, reviewers, projectKey); ok == false {
				fmt.Println("Unknown Error Occured... :(")
			}
			
		} else {
			fmt.Println("Make sure your Username, Base URL & Token are setup in .tongs_config...")
			interaction.Help()
		}
	} else if isUpdateReview {
		baseUrl := config.LoadBaseUrl()
		reviewers := config.LoadUserIds(reviewTemplate)
		hasUsername, username := config.LoadUsername()
		hasToken, token := config.LoadToken()
		if hasUsername && hasToken {
			crucible.UpdateReview(reviewTemplate, username, baseUrl, token, reviewers, projectId)
		} else {
			fmt.Println("Make sure your Username, Base URL & Token are setup in .tongs_config...")
			interaction.Help()
		}
	} 
	return
}
