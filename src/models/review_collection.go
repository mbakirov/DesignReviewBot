package models

import (
	"backend-rest-api/src/spreadsheets"
	"encoding/json"
	"log"
)

type ReviewCollection []*Review

func (reviewCollection ReviewCollection) FindForIssue(issueKey string) *ReviewCollection {
	result := &reviewCollection

	if issueKey == "" {
		panic("IssueKey is empty")
	}

	response := spreadsheets.GetReviews("issue_key="+issueKey)

	err := json.Unmarshal(response, result)
	if err != nil {
		log.Printf("[Review collection] Error while encode response FindForIssue: %s\n", err)
	}

	return result
}

func (reviewCollection *ReviewCollection) IsComplete() bool {
	result := true

	for _, review := range *reviewCollection {
		if !review.IsComplete() {
			result = false
		}
	}

	return result
}