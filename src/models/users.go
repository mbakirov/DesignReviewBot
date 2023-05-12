package models

import (
	"backend-rest-api/src/spreadsheets"
	"encoding/json"
	"log"
	"math/rand"
	"time"
)

const UserLevelLead = "lead"
const UserLevelSenior = "senior"
const UserLevelMiddle = "middle"
const UserLevelJunior = "junior"

type User struct {
	Name             string  `json:"name"`
	Level            string  `json:"level"`
	TgUsername       string  `json:"telegram_username"`
	JiraUsername     string  `json:"jira_username"`
	WeightUxTask     float64 `json:"weight_ux_task"`
	WeightDesignTaks float64 `json:"weight_design_task"`
	TasksCount		 int 	 `json:"tasks_count"`
	Random			 float64 `json:"random"`
	Score			 float64 `json:"score"`
}

type Users []*User

func (this *Users) Fetch() {
	sleepTime := 3600 * time.Second

	for {
		response := spreadsheets.GetUsers()

		err := json.Unmarshal(response, &this)
		if err != nil {
			log.Printf("[Users fetch] Error while encode response: %s\n", err)
		}

		log.Printf("Users: sleep %#v\n", this)
		time.Sleep(sleepTime)
	}
}

func (this *Users) FindByTgUsername(username string) *User {
	var result *User

	for _, user := range *this {
		if user.TgUsername == username {
			result = user
			break
		}
	}

	return result
}

func (this *User) GetWeight(weightName string) float64 {
	if weightName == "Design Task" {
		return this.WeightDesignTaks
	} else {
		return this.WeightUxTask
	}
}

func (this *User) GetScoreFor(taskType string) float64 {
	rand.Seed(time.Now().UTC().UnixNano())
	return float64(-1 * this.TasksCount) + this.GetWeight(taskType) * rand.Float64()
}