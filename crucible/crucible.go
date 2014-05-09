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
	"time"
	"strings"
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

func CreateReview(reviewName string, templateName string, reviewLength int64, userName string, baseUrl string, token string, reviewers_raw string, projectKey string) (bool, string) {
	
	
	
	dueDate := calculateDueDate(reviewLength)
	
	json := createReviewByteArray(reviewName, "REVIEW", strings.ToLower(userName), strings.ToLower(userName), dueDate, userName, projectKey, true)

	if id, ok := createReviewPost(json, token, baseUrl); ok {
		
		fmt.Println("Review Created:",id)
		fmt.Println("("+baseUrl+"/cru/"+id+")")
		fmt.Println("Using Tongs Template:", templateName)
		fmt.Println("Due Date Calculated:", dueDate)
		fmt.Println("Review Title:", reviewName)
		fmt.Println("Project Key:", projectKey)
		reviewers := strings.Split(reviewers_raw,",")
		for _, reviewer := range reviewers {
			clean_reviewer := strings.TrimSpace(strings.ToLower(reviewer))
			if(len(clean_reviewer)>0){
		    	fmt.Println("Adding Reviewer:",clean_reviewer)
		    	addReviewersPost(token, baseUrl, id, clean_reviewer)
			}
		}
		return true, id
	}
	return false, ""
}
func addReviewersPost(token string, baseUrl string, permaId string, userName string) {
	restUrl := baseUrl + "/rest-service/reviews-v1/" + permaId + "/reviewers?FEAUTH=" + token
	client := &http.Client{}
	req, _ := http.NewRequest("POST", restUrl, bytes.NewBuffer([]byte(userName)))
	client.Do(req)
}

func createReviewPost(json []byte, token string, baseUrl string) (string, bool) {

	restUrl := baseUrl + "/viewer/rest-service/reviews-v1?FEAUTH=" + token

	client := &http.Client{}
	req, err := http.NewRequest("POST", restUrl, bytes.NewBuffer(json))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		exitError("Unable to create code review", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		exitError("Unable to read response", err)
	}
	if hasStatus, _ := regexp.MatchString("status-code", string(body)); hasStatus {
		re := regexp.MustCompile("\"status-code\":(.*)}")
		status := re.FindStringSubmatch(string(body))
		if status != nil && len(status) > 1 {
			return status[1], false
		}
	}

	if hasPermaId, _ := regexp.MatchString("permaId", string(body)); hasPermaId {
		re := regexp.MustCompile("\"permaId\":{\"id\":\"(.*)\"},\"permaIdHistory\"")
		permaId := re.FindStringSubmatch(string(body))
		if permaId != nil && len(permaId) > 1 {
			return permaId[1], true
		}
	}

	return "", false
}

func createReviewByteArray(reviewName string, reviewType string, authorUsername string, creatorUsername string,
	dueDate string, moderatorUsername string, projectKey string, allowReviewerToJoin bool) []byte {

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
	restUrl := baseUrl + "/rest-service/auth-v1/login"

	resp, err := http.PostForm(restUrl, url.Values{"userName": {strings.ToLower(userName)}, "password": {password}})
	if err != nil {
		exitError("Unable to reach Crucible", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		exitError("Error unable to read Crucible http response", err)
	}
	rawToken := string(body[1:])

	isDown, _ := regexp.MatchString("Crucible Maintenance", rawToken)
	if isDown {
		exitError("Crucible is doing maintenance. Please try again in a bit.", nil)
	}

	hasToken, _ := regexp.MatchString("<token>(.*)</token>", rawToken)
	if hasToken == false {
		exitError("Unable to authenticate with Crucible", nil)
	}

	re := regexp.MustCompile("<token>(.*)</token>")
	token := re.FindStringSubmatch(rawToken)
	if len(token) == 2 {
		return token[1]
	}
	exitError("Unable to authenticate with Crucible", nil)
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
