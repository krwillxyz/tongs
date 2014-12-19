package config

import (
	"./goconfig/"
	"fmt"
	"os"
	"strings"
)

//Loads string of userids from the config file based on the key provided
//If no group is given or the reviewers key is not in the group, we will
//attempt to return the default values.
func LoadUserIds(key string) (reviewers string) {
	_, _, reviewers, _, _ = getConfigAtSection("reviewers", key, "string")
	return reviewers
}

func LoadProjectKey(key string) (projectKey string) {
	_, _, projectKey, _, _ = getConfigAtSection("project-key", key, "string")
	return projectKey
}

func LoadDuration(key string) (duration int64) {
	_, _, _, duration, _ = getConfigAtSection("duration", key, "int")
	return duration
}

func LoadReviewTitle(key string) (reviewTitle string) {
	_, _, reviewTitle, _, _ = getConfigAtSection("title", key, "string")
	return reviewTitle
}

func LoadBaseUrl() (baseUrl string) {
	_, _, baseUrl, _, _ = getConfigAtSection("crucible-baseurl", "settings", "string")
	return baseUrl
}

func LoadUsername() (bool, string) {

	if ok, token := LoadToken(); ok {
		if strings.Contains(token, ":") {
			return true, strings.Split(token, ":")[0]
		}
	}
	return false, ""
}

func LoadToken() (bool, string) {
	if ok, _, token, _, _ := getConfigAtSection("crucible-token",
		"settings", "string"); ok && token != "" {

		return true, token
	} else {
		return false, ""
	}
}

func LoadTemplates() {
	c, _ := readTongsConfig(true)
	fmt.Println("Templates:")
	for _, element := range c.GetSections() {
		if element != "default" && element != "settings" {
			fmt.Println(element)
		}
	}
}

func TemplateExists(template string) bool {
	c, _ := readTongsConfig(true)
	for _, element := range c.GetSections() {
		if element != "default" && element != "settings" {
			if element == template {
				return true
			}
		}
	}
	return false
}

func SaveToken(token string) bool {
	return writeConfigAtSection("settings", "crucible-token", token)
}

func SaveBaseUrl(baseUrl string) bool {
	return writeConfigAtSection("settings", "crucible-baseurl", baseUrl)
}

func ClearToken() bool {
	return writeConfigAtSection("settings", "crucible-token", "")
}

func writeConfigAtSection(section string, option string, value string) bool {
	c, _ := readTongsConfig(true)
	c.AddOption(section, option, value)
	c.WriteConfigFile("tongs.cfg", 0644, "")
	return true
}

func getConfigAtSection(option string, section string,
	datatype string) (bool, string, string, int64, bool) {

	c, _ := readTongsConfig(true)

	if c.HasOption(section, option) {

	} else {
		return false, "", "", 0, false
	}
	if datatype == "string" {
		value, _ := c.GetString(section, option)
		return true, section, value, 0, false
	}
	if datatype == "int" {
		value, _ := c.GetInt64(section, option)
		return true, section, "", value, false
	}
	if datatype == "bool" {
		value, _ := c.GetBool(section, option)
		return true, section, "", 0, value
	}
	return false, "", "", 0, false
}

func readTongsConfig(exit bool) (*goconfig.ConfigFile, error) {
	c, err := goconfig.ReadConfigFile("tongs.cfg")
	if err != nil && exit == true {
		exitError("Could not read tongs.cfg", err)
	}
	return c, err
}

func CreateConfigFile() bool {
	_, err := readTongsConfig(false)
	if err != nil {
		c := goconfig.NewConfigFile()
		c.AddSection("default")
		c.AddOption("default", "project-key", "PROJECT-KEY")
		c.AddOption("default", "duration", "5")
		c.AddOption("default", "reviewers", "user1, user2")
		c.AddOption("default", "title", "My Default Code Review Title")
		c.AddSection("my-team")
		c.AddOption("my-team", "title", "My Team Code Review Template Title")
		c.AddSection("settings")
		c.AddOption("settings", "crucible-baseurl", "http://crucible06.mycompany.com")
		c.AddOption("settings", "crucible-token", "")
		c.WriteConfigFile("tongs.cfg", 0644, "")
		return true
	}
	return false
}

func exitError(message string, err error) {
	fmt.Println(message + ": ")
	fmt.Println(err)
	os.Exit(0)
}
