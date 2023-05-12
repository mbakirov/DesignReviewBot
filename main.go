package main

import (
	//"backend-rest-api/src/handlers"
	"backend-rest-api/src/models"
	"backend-rest-api/src/telegram"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	//"github.com/gorilla/mux"
	"log"
)

//func init() {
//	// loads values from .env into the system
//	if err := godotenv.Load(); err != nil {
//		fmt.Println("No .env file found")
//		os.Exit(255)
//	}
//}

var Users models.Users
//BotToken := "bot2132831960:asdasdasd"

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
		os.Exit(255)
	}

	go Users.Fetch()
}

func main() {
	log.Printf("Users: %#v\n", Users)

	time.Sleep(time.Second * 4)

	// find reviews by issue key
	rc := models.ReviewCollection{}
	reviews := rc.FindForIssue("DSGN-2119")

	// parse task from message
	// author := Users.FindByTgUsername("mbakirov")
	// message := "#задача Посмотрите пожалуйста истории про электронный рецепт\nhttps://clck.ru/Z6kGD\nhttps://jira.eapteka.ru/browse/DSGN-2119"
	// task := models.Task{}
	// task = *task.ParseTask(message, author)

	// generate and save review
	// task.Size = models.TaskSizeMedium
	// reviews := task.CreateReview(Users)
	// for _, r := range reviews {
	//	 r.Save()
	// }

	val, _ := json.Marshal(reviews)
	log.Printf("Reviewers: %#v", string(val))

	//task.Save()

	os.Exit(2)

	//task := models.Task{}
	//task = *task.FindTask("https://jira.eapteka.ru/browse/DSGN-2084")

	/*router := mux.NewRouter()

	router.HandleFunc("/health", handlers.Health).Methods("GET")
	router.HandleFunc("/updates/{token}", handlers.Updates).Methods("POST")

	log.Printf("Starting server on the port %s...\n", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router))*/

	bot := telegram.NewBot(telegram.Token, true)
	//bot.SetWebhook("https://design-review-bot.mbakirov.ru/" + token)

	for u := range bot.PullUpdates("/" + telegram.Token) {
		update := telegram.Update{u}
		log.Printf("%+v\n", update)
		log.Printf("%+v\n", update.Message.Text)
		log.Println("isTaskUpdate", update.IsTask())

		if update.IsTask() {
			author := Users.FindByTgUsername(update.Message.From.UserName)
			log.Printf("UserName: %+v\n", update.Message.From.UserName)
			log.Printf("MessageId: %+v\n", update.Message.MessageID)
			log.Printf("User: %+v\n", author)

			task := models.Task{}
			task = *task.ParseTask(update.Message.Text, author)

			log.Printf("Task: %+v\n", task)
			log.Printf("IsComplete: %+v\n", task.IsComplete())

			task.Save()
		}
		//fmt.Println(matched) // false
	}
}
