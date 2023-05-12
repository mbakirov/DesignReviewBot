package handlers

import (
	"backend-rest-api/src/telegram"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io"
	"net/http"
)

func Updates(responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Printf("Incoming Updates request\n")

	params := mux.Vars(request)

	if (params["token"] != telegram.Token) {
		fmt.Printf("Wrong request token: %s\n", params["token"])
		return
	}

	bodyBytes, _ := io.ReadAll(request.Body)
	defer request.Body.Close()

	fmt.Printf("Reqeust from telegram: %s\n", string(bodyBytes))
}