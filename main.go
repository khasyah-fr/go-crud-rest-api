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
	Id    int             `json:"id"`
	Code  string          `json:"code"`
	Name  string          `json:"name"`
	Price decimal.Decimal `json:"price" sql:"type:decimal(16,2)"`
}

//Result berbentuk array product
type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {
	db, err = gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/go_rest_api_crud?charset=utf8&parseTime=True")

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

	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/api/products", createProductHandler).Methods("POST")
	router.HandleFunc("/api/products", indexProductHandler).Methods("GET")
	router.HandleFunc("/api/products/{id}", showProductHandler).Methods("GET")
	router.HandleFunc("/api/products/{id}", updateProductHandler).Methods("PUT")
	router.HandleFunc("/api/products/{id}", deleteProductHandler).Methods("DELETE")

	http.ListenAndServe(":9999", router)
}

func homeHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "Welcome to our homepage")
}

func createProductHandler(rw http.ResponseWriter, r *http.Request) {
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

func indexProductHandler(rw http.ResponseWriter, r *http.Request) {
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

func showProductHandler(rw http.ResponseWriter, r *http.Request) {
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

func updateProductHandler(rw http.ResponseWriter, r *http.Request) {
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

func deleteProductHandler(rw http.ResponseWriter, r *http.Request) {
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
