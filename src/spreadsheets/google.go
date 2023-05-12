package spreadsheets

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const Url = "https://script.google.com/macros/s/AKfycbxnmll13-5HhS6iD86T8bmdBbLP-gA04rrjinhv5CWhtshNYv28VMd2Mg4qyW3rqzAuIw/exec"

func makeQueryString(sheet, query string) string {
	result := "?dsgnbot=true&sheet=" + sheet

	if query != "" {
		result = result + "&" + query
	}

	return result
}

func getRequest(sheet, query string) []byte {
	url := Url + makeQueryString(sheet, query)
	response, err := http.Get(url)

	if err != nil {
		fmt.Printf("[Google] Error while GET request: %s\n", err)
	}

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	log.Println("getRequest", "url=", url, "response:", string(body))

	return body
}

func postRequest(sheet, query string, data []byte) []byte {
	var result []byte
	url := Url + makeQueryString(sheet, query)

	response, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		log.Fatalf("Error while POST-request: %s", err)
		return result
	}

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	log.Println("postRequest", "url=", url, "body", string(data), "response:", string(body))

	return body
}

func GetUsers() []byte {
	log.Println("[GetUsers] Fetching url...")
	return  getRequest("Сотрудники", "action=get_users")
}

func GetReviews(filter string) []byte {
	log.Println("[GetReviews] Fetching url...")
	return getRequest("Review", "action=get_reviews&" + filter)
}

//func FindTask(link string) *json.Decoder {
//	log.Println("[GetUsers] Fetching url...")
//	query := fmt.Sprintf("action=task&task_id=%s", link)
//	return json.NewDecoder(getRequest("DB", query))
//}

func SaveTask(data []byte) string {
	log.Println("[SaveTask] Posting data...")
	return string(postRequest("DB", "", data))
}

func SaveReview(data []byte) string {
	log.Println("[SaveReview] Posting data...")
	return string(postRequest("Review", "", data))
}