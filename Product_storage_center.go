package main


import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)


type Product struct {
	Product_ID        string `json:"product_id"`
	Name              string `json:"name"`
	Quantity_in_Stock string `json:"quantity_in_stock"`
	Unit_Price        string `json:"unit_price"`
}


var db *sql.DB
var err error

func main() {
	
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/Local_Host?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Connected")
	defer db.Close()

	router := mux.NewRouter()

	

	router.HandleFunc("/products", createProduct).Methods("POST")
	router.HandleFunc("/products", getProducts).Methods("GET")
	router.HandleFunc("/products/{id}", getProduct).Methods("GET")
	router.HandleFunc("/products/{id}", updateProduct).Methods("PUT")
	router.HandleFunc("/products/{id}", deleteProduct).Methods("DELETE")

	
	log.Fatal(http.ListenAndServe(":8000", router))
}


func createProduct(w http.ResponseWriter, r *http.Request) {

	statement, err := db.Prepare("INSERT INTO products(product_id,name,quantity_in_stock,unit_price)VALUES(?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	dataMap := make(map[string]string)
	json.Unmarshal(body, &dataMap)
	product_id := dataMap["product_id"]
	name := dataMap["name"]
	quantity_in_stock := dataMap["quantity_in_stock"]
	unit_price := dataMap["unit_price"]

	_, err = statement.Exec(product_id, name, quantity_in_stock, unit_price)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New product added")
}


func getProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	result, err := db.Query("SELECT * FROM products WHERE product_id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	
	var product Product

	for result.Next() {
		err := result.Scan(&product.Product_ID, &product.Name, &product.Quantity_in_Stock, &product.Unit_Price)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(product)
}


func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var products []Product

	result, err := db.Query("SELECT * FROM products")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var product Product
		err := result.Scan(&product.Product_ID, &product.Name, &product.Quantity_in_Stock, &product.Unit_Price)
		if err != nil {
			panic(err.Error())
		}
		products = append(products, product)

	}
	json.NewEncoder(w).Encode(products)
}


func updateProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	statement, err := db.Prepare("UPDATE products SET name = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	dataMap := make(map[string]string)
	json.Unmarshal(body, &dataMap)
	newTitle := dataMap["name"]

	_, err = statement.Exec(newTitle, params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Product %s updated", params["id"])

}


func deleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	statement, err := db.Prepare("DELETE FROM products WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = statement.Exec(params["id"])
	if err != nil {
		panic(err.Error)
	}
	fmt.Fprintf(w, "Products %s deleted", params["id"])
}
