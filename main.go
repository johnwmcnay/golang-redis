package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mediocregopher/radix"
	_ "github.com/mediocregopher/radix"
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

var Articles []Article

func main() {

	customConnFunc := func(network, addr string) (radix.Conn, error) {
		return radix.Dial(network, addr,
			radix.DialAuthPass("Vny0iYdqnFewcw5iPGzs7e1q0qZlzdkaSEzC9W4zJ6caqaVwLIcda7gq2Fy7ZAqq51IcqTGiQot6pwbvYOoLWoJ13M2UwQkEsyy2DI630TByF6PjOmsYltQjoukGg0SPMOZev9YwyFw7qYcyLaSCZz"),
		)
	}


	// this pool will use our ConnFunc for all connections it creates
	client, err := radix.NewPool("tcp", "127.0.0.1:6379", 10, radix.PoolConnFunc(customConnFunc))

	err = client.Do(radix.Cmd(nil, "SET", "foo", "someval"))
	if err != nil {
		// handle error
	}

	fmt.Println("Rest API v2.0 - Mux Routers")
	Articles = []Article{
		Article{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
		Article{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
	}
	handleRequests()
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(Articles)
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

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	vars := mux.Vars(r)
	key := vars["id"]

	json.Unmarshal(reqBody, &article)
	// update our global Articles array to include
	// our new Article
	for index, art := range Articles {
		if art.Id == key {
			Articles[index] = article
			json.NewEncoder(w).Encode(article)
			break
		}
	}
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	fmt.Println("Endpoint Hit: returnSingleArticle")
	// Loop over all of our Articles
	// if the article.Id equals the key we pass in
	// return the article encoded as JSON
	for _, article := range Articles {
		if article.Id == key {
			json.NewEncoder(w).Encode(article)
		}
	}
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewArticle")
	// get the body of our POST request
	// unmarshal this into a new Article struct
	// append this to our Articles array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	json.Unmarshal(reqBody, &article)
	// update our global Articles array to include
	// our new Article
	Articles = append(Articles, article)

	json.NewEncoder(w).Encode(article)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	// once again, we will need to parse the path parameters
	vars := mux.Vars(r)
	// we will need to extract the `id` of the article we
	// wish to delete
	id := vars["id"]

	fmt.Println("Endpoint Hit: Delete Article")

	// we then need to loop through all our articles
	for index, article := range Articles {
		// if our id path parameter matches one of our
		// articles
		if article.Id == id {
			// updates our Articles array to remove the
			// article
			Articles = append(Articles[:index], Articles[index+1:]...)
			break
		}
	}
}
