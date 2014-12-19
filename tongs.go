package main

import (
	"./config/"
	"./crucible/"
	"./interaction/"
	"fmt"
	"os"
	"strings"
)

var msg = map[string]string{
	"setup": "Make sure your Crucible Base Url, and Crucible Token are setup in your tongs.cfg file.",
}

func main() {

	isHelp,
		isCreateReview,
		isUpdateReview,
		isGetToken,
		reviewTemplate,
		isCreateConfig,
		projectId,
		isListTemplates := interaction.ParseCommandLineInputs(os.Args)

	if isHelp {
		interaction.Help()
	}

	if isListTemplates {
		if reviewTemplate != "default" {
			if config.TemplateExists(reviewTemplate) {
				fmt.Println("Title:    ", config.LoadReviewTitle(reviewTemplate))
				fmt.Println("Key:      ", config.LoadProjectKey(reviewTemplate))
				fmt.Println("Duration: ", config.LoadDuration(reviewTemplate))
				fmt.Println("User IDs: ", config.LoadUserIds(reviewTemplate))
			} else {
				fmt.Println("Template '" + reviewTemplate + "' not found.")
				config.LoadTemplates()
			}
		} else {
			config.LoadTemplates()
		}
	}

	if isCreateConfig {
		ok := config.CreateConfigFile()
		interaction.ConfigFileCreated(ok)

		if ok == false {
			baseUrl := config.LoadBaseUrl()
			fmt.Println("Crucible URL: ", baseUrl)
		} else {
			url := interaction.RequestUrl()
			config.SaveBaseUrl(url)
		}
		isGetToken = true
	}

	if isGetToken {
		baseUrl := config.LoadBaseUrl()
		if isCreateConfig == false {
			fmt.Println("Crucible Url: ", baseUrl)
		}
		username := interaction.RequestUsername()
		username = strings.ToLower(username)
		password := interaction.RequestPassword()
		for i := 0; i < 1000; i++ {
        	fmt.Println("")
    	}
		token := crucible.Login(username, password, baseUrl)
		config.SaveToken(token)
		fmt.Println("Token updated successfully!")
		fmt.Println("(you should now clear your terminal)")
	}

	if isCreateReview {

		if !config.TemplateExists(reviewTemplate){
			fmt.Println("Template '" + reviewTemplate + "' not found.")
			config.LoadTemplates()
			return
		}
		reviewTitle := config.LoadReviewTitle(reviewTemplate)
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
	}

	if isUpdateReview {

		if !config.TemplateExists(reviewTemplate){
			fmt.Println("Template '" + reviewTemplate + "' not found.")
			config.LoadTemplates()
			return
		}

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
