package main

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/nitishm/go-rejson"

	//"github.com/mediocregopher/radix"
	"io/ioutil"
	"log"
	"net/http"
)

type Article struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Desc    string `json:"Desc"`
	Content string `json:"Content"`
}

var client redis.Conn
var rh rejson.Handler
func main() {

	var err error

	client, err = redis.Dial("tcp", "127.0.0.1:6379",
		redis.DialPassword("Vny0iYdqnFewcw5iPGzs7e1q0qZlzdkaSEzC9W4zJ6caqaVwLIcda7gq2Fy7ZAqq51IcqTGiQot6pwbvYOoLWoJ13M2UwQkEsyy2DI630TByF6PjOmsYltQjoukGg0SPMOZev9YwyFw7qYcyLaSCZz"))

	defer client.Close()

	// this pool will use our ConnFunc for all connections it creates
	if err != nil {
		// handle error
	}
	fmt.Println("Rest API v2.0 - Mux Routers")

	handleRequests()
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	//json.NewEncoder(w).Encode(Articles)
}
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/all", returnAllArticles)
	myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/article/{id}", updateArticle).Methods("PUT")
	myRouter.HandleFunc("/article/{id}", returnSingleArticle)
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func updateArticle(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint Hit: updateArticle")

	//reqBody, _ := ioutil.ReadAll(r.Body)
	//var article Article
	//vars := mux.Vars(r)
	//key := vars["id"]
	//
	//json.Unmarshal(reqBody, &article)
	//// update our global Articles array to include
	//// our new Article
	//for index, art := range Articles {
	//	if art.Id == key {
	//		Articles[index] = article
	//		json.NewEncoder(w).Encode(article)
	//		break
	//	}
	//}
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//key := vars["id"]
	//
	//var article string
	//
	//err := client.Do(radix.Cmd(&article, "GET", "article:" + key))
	//if err != nil {
	//
	//}
	//
	//fmt.Println("Endpoint Hit: returnSingleArticle")
	//
	//json.NewEncoder(w).Encode(article)
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewArticle")

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article

	rh.SetRedigoClient(client)

	json.Unmarshal(reqBody, &article)

	_, err := client.Do("INCR", "articles:count")
	if err != nil {

	}

	id := "5"

	article.Id = id
	test, err := rh.JSONSet("article:" + id, ".", article)

	if err != nil {
		log.Fatalf("Failed to JSONSet" + err.Error())
	}
	fmt.Println(test)
	json.NewEncoder(w).Encode(article)

}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	// once again, we will need to parse the path parameters
	//vars := mux.Vars(r)
	//// we will need to extract the `id` of the article we
	//// wish to delete
	//id := vars["id"]
	//
	//fmt.Println("Endpoint Hit: Delete Article")
	//
	//// we then need to loop through all our articles
	//for index, article := range Articles {
	//	// if our id path parameter matches one of our
	//	// articles
	//	if article.Id == id {
	//		// updates our Articles array to remove the
	//		// article
	//		Articles = append(Articles[:index], Articles[index+1:]...)
	//		break
	//	}
	//}
}
