package main

import (
	"./config/"
	"./crucible/"
	"./interaction/"
	"os"
	"strings"
)

func main() {

	isHelp, isCreateReview, reviewTitle, reviewTemplate, isCreateConfig := interaction.ParseCommandLineInputs(os.Args)

	if isHelp {
		interaction.Help()
	} else if isCreateConfig {
		ok := config.CreateConfigFile() 
		interaction.ConfigFileCreated(ok)
	} else if isCreateReview {

		baseUrl := config.LoadBaseUrl()

		if reviewTitle == "" {
			reviewTitle = config.LoadReviewTitle(reviewTemplate)
		}

		projectKey := config.LoadProjectKey(reviewTemplate)

		reviewLength := config.LoadDuration(reviewTemplate)

		reviewers := config.LoadUserIds(reviewTemplate)

		hasUsername, username := config.LoadUsername()
		if hasUsername == false {
			username = interaction.RequestUsername()
			username := strings.ToLower(username)
			config.SaveUsername(username)
		}

		hasToken, token := config.LoadToken()
		if hasToken == false {
			password := interaction.RequestPassword()
			token := crucible.Login(username, password, baseUrl)
			if ok, _ := crucible.CreateReview(reviewTitle, reviewTemplate, reviewLength, username, baseUrl, token, reviewers, projectKey); ok {
				//success - save token
				config.SaveToken(token)
			} else {
				//failed to create review
			}
		} else {
			//user has a token lets see if it works
			if ok, _ := crucible.CreateReview(reviewTitle, reviewTemplate, reviewLength, username, baseUrl, token, reviewers, projectKey); ok == false {
				password := interaction.RequestPassword()
				token := crucible.Login(username, password, baseUrl)
				if ok, _ := crucible.CreateReview(reviewTitle, reviewTemplate, reviewLength, username, baseUrl, token, reviewers, projectKey); ok {
					config.SaveToken(token)
				} else {
					//could not create review
				}

			}
		}
	} 
	return
}
