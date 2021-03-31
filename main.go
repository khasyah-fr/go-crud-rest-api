package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shopspring/decimal"
)

var db *gorm.DB
var err error

type Product struct {
	Id          int             `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price" sql:"type:decimal(16,2)"`
	Amount      int             `json:"amount"`
}

//Result adalah sebuah array produk
type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {
	// Open Connection
	db, err = gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/fazztrack?charset=utf8&parseTime=True")

	// Error Handling Connection
	if err != nil {
		log.Fatal("Connection failed to open")
	} else {
		log.Println("Connection established successfully")
	}

	db.AutoMigrate(&Product{})
	handleRequests()
}

func handleRequests() {
	log.Println("Start the development server at http://127.0.0.1:9999")

	router := mux.NewRouter().StrictSlash(true)

	// Handler Not Found
	router.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusNotFound)

		res := Result{Code: 404, Message: "Method not found"}
		response, _ := json.Marshal(res)
		rw.Write(response)
	})

	// Method Not Allowed
	router.MethodNotAllowedHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusMethodNotAllowed)

		res := Result{Code: 403, Message: "Method not allowed"}
		response, _ := json.Marshal(res)
		rw.Write(response)
	})

	router.HandleFunc("/", homePage)
	router.HandleFunc("/api/products", createHandler).Methods("POST")
	router.HandleFunc("/api/products", indexHandler).Methods("GET")
	router.HandleFunc("/api/products/{id:[0-9]+}", showHandler).Methods("GET")
	router.HandleFunc("/api/products/{id:[0-9]+}", updateHandler).Methods("PUT")
	router.HandleFunc("/api/products/{id:[0-9]+}", deleteHandler).Methods("DELETE")

	http.ListenAndServe(":9999", router)
}

func homePage(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "Welcome to the homepage")
}

func createHandler(rw http.ResponseWriter, r *http.Request) {
	payload, _ := ioutil.ReadAll(r.Body)

	var product Product
	json.Unmarshal(payload, &product)

	db.Create(&product)

	res := Result{Code: 200, Data: product, Message: "Product created successfully"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(result)
}

func indexHandler(rw http.ResponseWriter, r *http.Request) {
	products := []Product{}
	db.Find(&products)

	res := Result{Code: 200, Data: products, Message: "Products indexed successfully"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(result)
}

func showHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	var product Product

	db.First(&product, productID)

	res := Result{Code: 200, Data: product, Message: "Product showed successfully"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(result)
}

func updateHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	payload, _ := ioutil.ReadAll(r.Body)

	var productUpdated Product
	json.Unmarshal(payload, &productUpdated)

	var product Product
	db.First(&product, productID)
	db.Model(&product).Updates(productUpdated)

	res := Result{Code: 200, Data: product, Message: "Product updated successfully"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(result)
}

func deleteHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	var product Product

	db.First(&product, productID)
	db.Delete(&product)

	res := Result{Code: 200, Message: "Product deleted successfully"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(result)
}
