package main

import (
	"bytes"
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


type Person struct {
	Id 			string `json:"id"`
	Name 		string `json:"name"`
	Age 		string `json:"age"`
	DateJoined 	string `json:"dateJoined"`
	DateUpdated string `json:"dateUpdated"`
}

type PersonExceptions struct {
	Person
	Other *string // Other won't raise an error
}

type Job struct {
	Id 			string `json:"id"`
	Name 		string `json:"test"`
	Age 		string `json:"yes"`
}

var client redis.Conn
var rh rejson.Handler
var objectMap = map[string]map[string]interface{}{
	"persons": {
		"object": Person{},
		"exception": &PersonExceptions{},
	},
	"jobs": {
		"object": Job{},
		"exception": nil,
	}}

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
	myRouter.HandleFunc("/{object}", returnAllObjects).Methods("GET")
	myRouter.HandleFunc("/{object}/{id}", returnSingleObjects).Methods("GET")
	myRouter.HandleFunc("/{object}/{id}", deleteObjects).Methods("DELETE")
	myRouter.HandleFunc("/{object}/{id}", updateObjects).Methods("PUT")
	myRouter.HandleFunc("/{object}", createNewObjects).Methods("POST")

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func updateObjects(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint Hit: updateArticle")

	vars := mux.Vars(r)
	id := vars["id"]
	objectName := vars["object"]

	if !objectIsAllowed(objectName) {
		http.Error(w, "Bad Request", 400)
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)

	object := objectMap[objectName]["object"]
	exception := objectMap[objectName]["exception"]

	dec := json.NewDecoder(bytes.NewReader(reqBody))
	dec.DisallowUnknownFields()

	if err := dec.Decode(exception); err != nil {
		fmt.Println(err)
		http.Error(w, "Bad Request", 400)
		return
	}

	json.Unmarshal(reqBody, &object)

	m := object.(map[string]interface{})

	if id != m["id"] {
		http.Error(w, "Bad Request", 400)
		return
	}

	_, err := rh.JSONSet(objectName + ":" + id, ".", m)

	if err != nil {
		log.Fatalf("Failed to JSONSet" + err.Error())
	}

	json.NewEncoder(w).Encode(m)
}

func returnAllObjects(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")

	vars := mux.Vars(r)
	objectName := vars["object"]

	results, err := client.Do("SCAN", "0", "MATCH", objectName + ":*")

	if err != nil {

	}

	arrayOfByteArrays := reflect.ValueOf(results).Index(1)
	object := objectMap[objectName]["object"]

	var jsonList []map[string]interface{}

	for i := 0; i < arrayOfByteArrays.Elem().Len(); i++ {

		key, _ := redis.String(arrayOfByteArrays.Elem().Index(i).Elem().Interface(), err)

		byteArray, _ := redis.Bytes(rh.JSONGet(key, "."))

		err = json.Unmarshal(byteArray, &object)
		m := object.(map[string]interface{})

		jsonList = append(jsonList, m)
	}
	json.NewEncoder(w).Encode(jsonList)
}

func returnSingleObjects(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	objectName := vars["object"]

	res, err := redis.Bytes(rh.JSONGet(objectName + ":" + key, "."))
	if err != nil {
		panic(err)
	}

	object := objectMap[objectName]["object"]

	err = json.Unmarshal(res, &object)
	m := object.(map[string]interface{})

	if err != nil {
		log.Fatalf("Failed to JSON Unmarshal")
		return
	}

	fmt.Println("Endpoint Hit: returnSingleArticle")
	json.NewEncoder(w).Encode(m)
}

func objectIsAllowed(objectName string) bool {
	return objectMap[objectName] != nil
}

func createNewObjects(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewArticle")

	vars := mux.Vars(r)
	objectName := vars["object"]

	if !objectIsAllowed(objectName) {
		http.Error(w, "Bad Request", 400)
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)

	object := objectMap[objectName]["object"]
	exception := objectMap[objectName]["exception"]

	dec := json.NewDecoder(bytes.NewReader(reqBody))
	dec.DisallowUnknownFields()

	if err := dec.Decode(exception); err != nil {
		fmt.Println(err)
		http.Error(w, "Bad Request", 400)
		return
	}

	json.Unmarshal(reqBody, &object)

	res, err := client.Do("INCR", "count:" + objectName)
	if err != nil {

	}

	id := fmt.Sprintf("%v", res)
	m := object.(map[string]interface{})
	m["id"] = id

	_, err = rh.JSONSet(objectName + ":" + id, ".", m)

	if err != nil {
		log.Fatalf("Failed to JSONSet" + err.Error())
	}

	json.NewEncoder(w).Encode(m)

}

func deleteObjects(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	obj := vars["object"]

	fmt.Println("Endpoint Hit: Delete Article")

	_, err := rh.JSONDel(obj + ":" + id, ".")

	if err != nil {
		log.Fatalf("Failed to JSONDel" + err.Error())
	}


}
