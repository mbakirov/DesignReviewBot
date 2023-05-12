package models

import (
	"backend-rest-api/src/spreadsheets"
	"encoding/json"
	"github.com/andygrunwald/go-jira"
	"log"
	"os"
	"sort"
	"time"
)

type Review struct {
	IssueKey string `json:"issue_key"`
	IsUrgent bool `json:"is_urgent"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	AssigneeTelegramUsername string `json:"assignee"`
	WasReviewed bool `json:"was_reviewed"`
	Text string `json:"text"`
	issue *jira.Issue
}

func (review Review) NewReview(users []*User, t Task) []Review {
	var result []Review
	var reviewers []*User
	taskType := t.GetIssueType()

	switch t.Size {
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

func (review *Review) IsComplete() bool {
	return review.WasReviewed
}

func (review *Review) Save() {
	data, _ := json.Marshal(review)
	spreadsheets.SaveReview(data)
}


//func (review *Review) FindOrCreate() *[]Review {
//	if review.IssueKey == "" {
//		panic("IssueKey is empty")
//	}
//
//	reviews := review.Find()
//
//	if len(*reviews) == 0 {
//
//	}
//}

func (review *Review) GetJiraIssue() *jira.Issue {
	if review.issue != nil {
		return review.issue
	}

	login, _ := os.LookupEnv("JIRA_LOGIN")
	password, _ := os.LookupEnv("JIRA_PASSWORD")
	url, _ := os.LookupEnv("JIRA_URL")

	tp := jira.BasicAuthTransport{
		Username: login,
		Password: password,
	}

	jiraClient, _ := jira.NewClient(tp.Client(), url)
	issue, _, err := jiraClient.Issue.Get(review.IssueKey, nil)

	if err == nil {
		review.issue = issue
	} else {
		log.Fatal("Error while requesting jira issue", review.IssueKey, err)
	}

	return review.issue
}