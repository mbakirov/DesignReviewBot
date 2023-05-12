package models

import (
	"backend-rest-api/src/spreadsheets"
	"encoding/json"
	"github.com/andygrunwald/go-jira"
	"github.com/google/go-cmp/cmp"
	"log"
	"os"
	"regexp"
	"sort"
	"time"
)

const TaskSizeSmall = "S"
const TaskSizeMedium = "M"
const TaskSizeLarge = "L"
const TaskUrgentHashtag = `\#срочно`
const TaskJiraRegex = `https:\/\/jira.eapteka.ru\/browse\/([A-Z\-0-9]+)`

//goland:noinspection GoVetStructTag
type Task struct {
	JiraLink   string `json:"JiraLink"`
	Size       string `json:"size"`
	Author     *User  `json:"Author"`
	Reviewers  Users  `json:"Reviewers"`
	Reviewed   Users  `json:"Reviewed"`
	Text       string `json:"Text"`
	isComplete bool   `json:"isComplete"`
	isUrgent   bool
	jiraIssue  *jira.Issue
}

func parseJiraLink(text string) string {
	regex, _ := regexp.Compile(TaskJiraRegex)
	jiraLink := regex.FindString(text)

	regex, _ = regexp.Compile(`[A-Z\-0-9]+$`)
	return regex.FindString(jiraLink)
}

func (t Task) ParseTask(text string, author *User) *Task {
	jiraLink := parseJiraLink(text)

	return &Task{
		Author:   author,
		JiraLink: jiraLink,
		Text:     text,
	}
}

func (t *Task) IsComplete() bool {
	if t.isComplete {
		return t.isComplete
	}

	if len(t.Reviewers) == 0 {
		return false
	}

	if cmp.Equal(t.Reviewers, t.Reviewed) {
		t.isComplete = true
	}

	return t.isComplete
}

func (t *Task) IsUrgent() bool {
	if t.isUrgent {
		return t.isUrgent
	}

	regex, _ := regexp.Compile(TaskUrgentHashtag)
	urgentHash := regex.FindString(t.Text)

	log.Println("urgentHash", urgentHash)

	t.isUrgent = urgentHash != ""

	return t.isUrgent
}

//func (t *Task) GetIncompleteReviewers() Users {
//	result := t.Reviewers
//	reviewed := map[string]bool{}
//
//	for _, user := range t.Reviewed {
//		reviewed[user.TgUsername] = true
//	}
//
//	for i, reviewer := range result {
//		if _, ok := reviewed[reviewer.TgUsername]; ok {
//			result = append(result[:i], result[i+1:]...)
//		}
//	}
//
//	return result
//}

func (t *Task) GetSize() string {
	return t.Size
}

func (t *Task) CreateReview(users []*User) []Review {
	var result []Review
	var reviewers []*User
	taskType := t.GetIssueType()

	log.Println("Searching reviewers for task:", t, "(", taskType, ")", "size", t.Size)

	switch t.GetSize() {
	case TaskSizeLarge:
		leads := filterUsersByLevel(users, UserLevelLead)
		sort.Slice(leads, func(first, second int) bool {
			return leads[first].GetScoreFor(taskType) > leads[second].GetScoreFor(taskType)
		})

		seniors := filterUsersByLevel(users, UserLevelSenior)
		sort.Slice(seniors, func(first, second int) bool {
			return seniors[first].GetScoreFor(taskType) > seniors[second].GetScoreFor(taskType)
		})

		reviewers = append(reviewers, leads[0], seniors[0])
		break
	case TaskSizeMedium:
		reviewers = append(filterUsersByLevel(users, UserLevelSenior), filterUsersByLevel(users, UserLevelMiddle)...)
		break
	default:
		reviewers = append(filterUsersByLevel(users, UserLevelMiddle), filterUsersByLevel(users, UserLevelJunior)...)
		break
	}

	if len(reviewers) > 2 {
		sort.Slice(reviewers, func(first, second int) bool {
			return reviewers[first].GetScoreFor(taskType) > reviewers[second].GetScoreFor(taskType)
		})
	}

	log.Println("First reviewer weight:",
		reviewers[0].GetScoreFor(taskType),
		reviewers[0].Name,
		reviewers[0].Level,
		reviewers[0].TasksCount,
		reviewers[0].GetWeight(taskType))

	log.Println("Second reviewer weight:",
		reviewers[1].GetScoreFor(taskType),
		reviewers[1].Name,
		reviewers[1].Level,
		reviewers[1].TasksCount,
		reviewers[1].GetWeight(taskType))

	for idx, reviewer := range reviewers {
		if idx > 1 {
			break
		}

		review := Review{
			IssueKey: t.GetJiraIssue().Key,
			IsUrgent: t.IsUrgent(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			AssigneeTelegramUsername: reviewer.TgUsername,
			WasReviewed: false,
		}

		result = append(result, review)
	}

	return result
}

func filterUsersByLevel(users []*User, level string) []*User {
	var result []*User

	for _, user := range users {
		if user.Level == level {
			result = append(result, user)
		}
	}

	return result
}

func (t *Task) Save() {
	data, _ := json.Marshal(t)
	spreadsheets.SaveTask(data)
}

func (t *Task) GetJiraIssue() *jira.Issue {
	if t.jiraIssue != nil {
		return t.jiraIssue
	}

	login, _ := os.LookupEnv("JIRA_LOGIN")
	password, _ := os.LookupEnv("JIRA_PASSWORD")
	url, _ := os.LookupEnv("JIRA_URL")

	tp := jira.BasicAuthTransport{
		Username: login,
		Password: password,
	}

	jiraClient, _ := jira.NewClient(tp.Client(), url)
	issue, _, err := jiraClient.Issue.Get(t.JiraLink, nil)

	if err == nil {
		t.jiraIssue = issue
	} else {
		log.Fatal("Error while requesting jira issue", t.JiraLink, err)
	}

	return t.jiraIssue
}

func (t *Task) GetIssueType() string {
	result := "Design Task"
	issue := t.GetJiraIssue()

	if issue.Fields.Type.Name == "UX/UI Task" {
		result = issue.Fields.Type.Name
	}

	log.Println(issue.Fields.Type.Name)

	return result
}