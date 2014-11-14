package crucible

/*
	CreateReview
	------------
	Go API for creating a new review based on the passed in elements
	https://docs.atlassian.com/fisheye-crucible/latest/wadl/crucible.html#d2e2
*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

//struct definitions of the JSON structure to use the Crucible Rest API
type Person struct {
	UserName string `json:"userName"`
}

type ReviewData struct {
	Name                 string `json:"name"`
	Type                 string `json:"type"`
	Creator              Person `json:"creator"`
	Author               Person `json:"author"`
	Moderator            Person `json:"moderator"`
	ProjectKey           string `json:"projectKey"`
	AllowReviewersToJoin bool   `json:"allowReviewersToJoin"`
	DueDate              string `json:"dueDate"`
}

type CodeReview struct {
	ReviewData ReviewData `json:"reviewData"`
}

var msg = map[string]string{
	"update-review":           "Attempting to Update Review...\n",
	"review-title":            "%s (%s)\n",
	"review-created-id":       "Review Created!\nId: %s\n",
	"review-created-url":      "Url: %s/cru/%s\n",
	"review-created-template": "Template: %s\n",
	"review-created-due":      "Due Date: %s\n",
	"review-created-title":    "Title: %s\n",
	"review-created-key":      "Project Key: %s\n",
	"review-created-reviewer": "Adding Reviewer: %s\n",
	"review-not-found":        "Code review with id %s not found...\n",
	"not-author":              "You are not the author of this code review...\n",
	"bad-response":            "Unable to read Crucible response\n",
	"bad-http-response":       "Error. Unable to read Crucible HTTP response\n",
	"no-crucible":             "Unable to reach Crucible\n",
	"crucible-maintenance":    "Crucible is doing maintenance. Please try again in a bit.\n",
	"unable-to-authenticate":  "Unable to authenticate with Crucible\n",
	"unable-to-create":        "Unable to create code review\n",
}

var regex = map[string]string{
	"extract-author":       "<author>.*<userName>(.*)</userName>",
	"extract-title":        "<name>(.*)</name>",
	"extract-token":        "<token>(.*)</token>",
	"crucible-maintenance": "Crucible Maintenance",
	"status-code":          "\"status-code\":(.*)}",
	"perma-id":             "\"permaId\":{\"id\":\"(.*)\"},\"permaIdHistory\"",
}

var urls = map[string]string{
	"crucible-login":         "%s/rest-service/auth-v1/login",
	"crucible-add-reviewer":  "%s/rest-service/reviews-v1/%s/reviewers?FEAUTH=%s",
	"crucible-review":        "%s/rest-service/reviews-v1/%s",
	"crucible-create-review": "%s/rest-service/reviews-v1?FEAUTH=%s",
}

func UpdateReview(templateName string,
	userName string, baseUrl string, token string,
	reviewers_raw string, projectId string) {

	fmt.Printf(msg["update-review"])
	title := confirmReviewOwnership(baseUrl, projectId, userName)
	fmt.Printf(msg["update-review-title"], title, projectId)
	reviewers := strings.Split(reviewers_raw, ",")
	for _, reviewer := range reviewers {
		clean_reviewer := strings.TrimSpace(strings.ToLower(reviewer))
		if len(clean_reviewer) > 0 {
			fmt.Printf(msg["adding-reviewer"], clean_reviewer)
			addReviewersPost(token, baseUrl, projectId, clean_reviewer)
		}
	}
}

func CreateReview(reviewName string, templateName string,
	reviewLength int64, userName string, baseUrl string,
	token string, reviewers_raw string,
	projectKey string) (bool, string) {

	dueDate := calculateDueDate(reviewLength)

	json := createReviewByteArray(reviewName, "REVIEW",
		strings.ToLower(userName), strings.ToLower(userName),
		dueDate, userName, projectKey, true)

	if id, ok := createReviewPost(json, token, baseUrl); ok {

		fmt.Printf(msg["review-created-id"], id)
		fmt.Printf(msg["review-created-url"], baseUrl, id)
		fmt.Printf(msg["review-created-template"], templateName)
		fmt.Printf(msg["review-created-due"], dueDate)
		fmt.Printf(msg["review-created-title"], reviewName)
		fmt.Printf(msg["review-created-key"], projectKey)
		reviewers := strings.Split(reviewers_raw, ",")
		for _, reviewer := range reviewers {
			clean_reviewer := strings.TrimSpace(strings.ToLower(reviewer))
			if len(clean_reviewer) > 0 {
				fmt.Printf(msg["review-created-reviewer"], clean_reviewer)
				addReviewersPost(token, baseUrl, id, clean_reviewer)
			}
		}
		return true, id
	}
	return false, ""
}
func addReviewersPost(token string, baseUrl string,
	permaId string, userName string) {

	restUrl := fmt.Sprintf(urls["crucible-add-reviewer"], baseUrl, permaId, token)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", restUrl, bytes.NewBuffer([]byte(userName)))
	client.Do(req)
}

func confirmReviewOwnership(baseUrl string,
	permaId string, userName string) string {

	restUrl := fmt.Sprintf(urls["crucible-create-review"], baseUrl, permaId)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", restUrl, nil)
	req.Header.Add("Accept", "application/xml")
	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != 200 {
		exitError(fmt.Sprintf(msg["review-not-found"], permaId), nil)
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		exitError(msg["bad-response"], err)
	}
	title := extractTitle(string(body))

	if userName == extractAuthor(string(body)) {
		return title
	}
	fmt.Printf(msg["review-title"], title, permaId)
	exitError(msg["not-author"], nil)
	return ""
}

func extractAuthor(body string) (author string) {
	if hasAuthor, _ := regexp.MatchString("author", string(body)); hasAuthor {
		re := regexp.MustCompile(regex["extract-author"])
		author := re.FindStringSubmatch(string(body))
		if author != nil && len(author) > 1 {
			return author[1]
		}
	}
	return ""
}

func extractTitle(body string) (name string) {
	if hasName, _ := regexp.MatchString("name", string(body)); hasName {
		re := regexp.MustCompile(regex["extract-title"])
		name := re.FindStringSubmatch(string(body))
		if name != nil && len(name) > 1 {
			return name[1]
		}
	}
	return ""
}

func createReviewPost(json []byte, token string, baseUrl string) (string, bool) {

	restUrl := fmt.Sprintf(urls["crucible-create-review"], baseUrl, token)
	client := &http.Client{}
	req, err := http.NewRequest("POST", restUrl, bytes.NewBuffer(json))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		exitError(msg["unable-to-create"], err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		exitError(msg["bad-response"], err)
	}
	if hasStatus, _ := regexp.MatchString("status-code", string(body)); hasStatus {
		re := regexp.MustCompile(regex["status-code"])
		status := re.FindStringSubmatch(string(body))
		if status != nil && len(status) > 1 {
			return status[1], false
		}
	}

	if hasPermaId, _ := regexp.MatchString("permaId", string(body)); hasPermaId {
		re := regexp.MustCompile(regex["perma-id"])
		permaId := re.FindStringSubmatch(string(body))
		if permaId != nil && len(permaId) > 1 {
			return permaId[1], true
		}
	}

	return "", false
}

func createReviewByteArray(reviewName string, reviewType string,
	authorUsername string, creatorUsername string,
	dueDate string, moderatorUsername string, projectKey string,
	allowReviewerToJoin bool) []byte {

	author := Person{
		UserName: authorUsername,
	}
	creator := Person{
		UserName: creatorUsername,
	}
	moderator := Person{
		UserName: moderatorUsername,
	}

	reviewData := ReviewData{
		Name:                 reviewName,
		Type:                 reviewType,
		Author:               author,
		Creator:              creator,
		DueDate:              dueDate,
		Moderator:            moderator,
		ProjectKey:           projectKey,
		AllowReviewersToJoin: allowReviewerToJoin,
	}

	codeReview := CodeReview{
		ReviewData: reviewData,
	}

	b, err := json.Marshal(codeReview)
	if err != nil {
		return nil
	}
	return b
}

func calculateDueDate(days int64) string {
	const layout = "2006-01-02T15:04:05.001-0700"
	t := time.Now().AddDate(0, 0, int(days))
	return t.Format(layout)
}

func Login(userName string, password string, baseUrl string) string {
	restUrl := fmt.Sprintf(urls["crucible-login"], baseUrl)

	resp, err := http.PostForm(restUrl,
		url.Values{"userName": {strings.ToLower(userName)},
			"password": {password}})

	if err != nil {
		exitError(msg["no-crucible"], err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		exitError(msg["bad-http-response"], err)
	}
	rawToken := string(body[1:])

	isDown, _ := regexp.MatchString(regex["crucible-maintenance"], rawToken)
	if isDown {
		exitError(msg["crucible-maintenance"], nil)
	}

	hasToken, _ := regexp.MatchString(regex["extract-token"], rawToken)
	if hasToken == false {
		exitError(msg["unable-to-authenticate"], nil)
	}

	re := regexp.MustCompile(regex["extract-token"])
	token := re.FindStringSubmatch(rawToken)
	if len(token) == 2 {
		return token[1]
	}
	exitError(msg["unable-to-authenticate"], nil)
	return ""
}

func exitError(message string, err error) {
	if err != nil {
		fmt.Println(message + ": ")
		fmt.Println(err)
	} else {
		fmt.Println(message)
	}
	os.Exit(0)
}
