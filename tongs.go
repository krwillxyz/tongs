package main

import (
	"./config/"
	"./crucible/"
	"./interaction/"
	"fmt"
	"os"
	"strings"
)

var msg map[string]string {
	"setup": "Make sure your username, base Crucible url, and Crucible token are setup in your .tongs_config file.",
}

func main() {

	isHelp,
		isCreateReview,
		isUpdateReview,
		isGetToken,
		reviewTitle,
		reviewTemplate,
		isCreateConfig,
		projectId := interaction.ParseCommandLineInputs(os.Args)

	if isHelp {
		interaction.Help()
	} else if isCreateConfig {
		ok := config.CreateConfigFile()
		interaction.ConfigFileCreated(ok)
	} else if isGetToken {
		baseUrl := config.LoadBaseUrl()
		fmt.Println("Crucible Url: ", baseUrl)
		username := interaction.RequestUsername()
		username = strings.ToLower(username)
		password := interaction.RequestPassword()
		token := crucible.Login(username, password, baseUrl)
		fmt.Println("Recieved Token: ", token)
		config.SaveUsername(username)
		config.SaveToken(token)
		fmt.Println("Token Saved Successfully!")
		fmt.Println("(you should now clear your terminal)")
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

			if ok, _ := crucible.CreateReview(reviewTitle, reviewTemplate,
				reviewLength, username, baseUrl, token, reviewers,
				projectKey); ok == false {

				fmt.Println("Unable to create new code review.")
				fmt.Println("Has your password changed?")
			}

		} else {
			fmt.Println(msg["setup"])
			interaction.Help()
		}
	} else if isUpdateReview {
		baseUrl := config.LoadBaseUrl()
		reviewers := config.LoadUserIds(reviewTemplate)
		hasUsername, username := config.LoadUsername()
		hasToken, token := config.LoadToken()
		if hasUsername && hasToken {
			crucible.UpdateReview(reviewTemplate, username,
				baseUrl, token, reviewers, projectId)

		} else {
			fmt.Println(msg["setup"])
			interaction.Help()
		}
	}
	return
}
