package main

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/nitishm/go-rejson"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
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

	rh.SetRedigoClient(client)

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


func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/articles", returnAllArticles).Methods("GET")
	myRouter.HandleFunc("/articles/{id}", returnSingleArticle).Methods("GET")
	myRouter.HandleFunc("/articles/{id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/articles/{id}", updateArticle).Methods("PUT")
	myRouter.HandleFunc("/articles", createNewArticle).Methods("POST")

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func updateArticle(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint Hit: updateArticle")

	vars := mux.Vars(r)
	id := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article

	json.Unmarshal(reqBody, &article)

	if id != article.Id {
		return
	}

	_, err := rh.JSONSet("article:" + id, ".", article)

	if err != nil {
		log.Fatalf("Failed to JSONSet" + err.Error())
	}

	json.NewEncoder(w).Encode(article)
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")

	results, err := client.Do("SCAN", "0", "MATCH", "article:*")

	if err != nil {

	}

	arrayOfByteArrays := reflect.ValueOf(results).Index(1)
	article := Article{}

	var jsonList []Article

	for i := 0; i < arrayOfByteArrays.Elem().Len(); i++ {

		key, _ := redis.String(arrayOfByteArrays.Elem().Index(i).Elem().Interface(), err)

		byteArray, _ := redis.Bytes(rh.JSONGet(key, "."))

		err = json.Unmarshal(byteArray, &article)
		jsonList = append(jsonList, article)
	}
	json.NewEncoder(w).Encode(jsonList)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	res, err := redis.Bytes(rh.JSONGet("article:" + key, "."))
	if err != nil {
		panic(err)
	}

	article := Article{}
	err = json.Unmarshal(res, &article)
	if err != nil {
		log.Fatalf("Failed to JSON Unmarshal")
		return
	}

	fmt.Println("Endpoint Hit: returnSingleArticle")
	json.NewEncoder(w).Encode(article)
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewArticle")

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article

	json.Unmarshal(reqBody, &article)

	res, err := client.Do("INCR", "articles:count")
	if err != nil {

	}

	id := fmt.Sprintf("%v", res)
	article.Id = id
	_, err = rh.JSONSet("article:" + id, ".", article)

	if err != nil {
		log.Fatalf("Failed to JSONSet" + err.Error())
	}

	json.NewEncoder(w).Encode(article)

}

func deleteArticle(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id := vars["id"]

	fmt.Println("Endpoint Hit: Delete Article")
	fmt.Println(id)
	_, err := rh.JSONDel("article:" + id, ".")

	if err != nil {
		log.Fatalf("Failed to JSONDel" + err.Error())
	}


}
