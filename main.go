package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"io/ioutil"
	_ "time"
	// "strings"
	_ "strconv"
	"os"
)

type shoppingCart struct {
	ShopCartID   int    `json:"shopCartID"`
	ProductID	 int	`json:"productID"`
	Quantity	 int    `json:"quantity"`
}

type shoppingCartUser struct {
	ShopCartID   int    `json:"shopCartID"`
	UserID	     int 	`json:"userID"`
	IsCheckout 	 bool 	`json:"isCheckout"`
}

type checkout struct {
	ShopCartID   int    	`json:"shopCartID"`
	EmailAddr	 string		`json:"emailAddr"`
	Shipping	 string 	`json:"shipping"`
	PostalCode 	 int 		`json:"postalCode"`
	CreditCard	 string 	`json:"creditCard"`
	TotalPayment float64 	`json:"totalPayment"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/shoppingCart", shoppingCartItemEndpoint).Methods("GET","POST")
	router.HandleFunc("/api/v1/shoppingCartUser", shoppingCartCreateEndpoint).Methods("GET","POST")
	router.HandleFunc("/api/v1/checkout", checkoutEndpoint).Methods("GET","POST", "OPTIONS")
	fmt.Println("Listening at port 5000")

	log.Fatal(http.ListenAndServe(":5000", router))
}

func shoppingCartCreateEndpoint(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	dbApiKey := os.Getenv("API_KEY")
	dbReadApiKey := os.Getenv("READ_API_KEY")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")

	if r.Method == "POST" {
		querystringmap := r.URL.Query()
		userID := querystringmap.Get("UserID")

		//Opening database connection
		db, err := sql.Open("mysql", dbUser + ":" + dbPassword + "@tcp(" + dbApiKey + ")/" + dbName)
		// handle error upon failure
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		defer db.Close()
		
		//inserting values into passenger table
		_, err = db.Exec("insert into shopping_cart_user (UserID) values(?)", userID)
		//Handling error of SQL statement
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusAccepted)
	} else if r.Method == "GET" {
		querystringmap := r.URL.Query()
		userID := querystringmap.Get("UserID")
		//Opening database connection
		db, err := sql.Open("mysql", dbUser + ":" + dbPassword + "@tcp(" + dbReadApiKey + ")/" + dbName)
		// handle error upon failure
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		defer db.Close()
		
		var shoppingCartUsers []shoppingCartUser
		results, err := db.Query("select * from shopping_cart_user where UserID = ?", userID)
		//Handling error of SQL statement
		if err != nil {
			http.Error(w, "Missing data", http.StatusBadRequest)
			panic(err.Error())
		}
		for results.Next() {
			var shoppingCartUser shoppingCartUser
			err = results.Scan( &shoppingCartUser.ShopCartID, &shoppingCartUser.UserID, &shoppingCartUser.IsCheckout )
				if err != nil {
					http.Error(w, "Missing data", http.StatusBadRequest)
				} else {
					shoppingCartUsers = append(shoppingCartUsers, shoppingCartUser)
				}
		}
		
		output, _ := json.Marshal(shoppingCartUsers)
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, string(output))

	} else{
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}

func checkoutEndpoint(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	dbApiKey := os.Getenv("API_KEY")
	dbReadApiKey := os.Getenv("READ_API_KEY")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")

	if r.Method == "POST"{
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var checkout checkout
			if err := json.Unmarshal(body, &checkout); err == nil{
				//Opening database connection
				db, err := sql.Open("mysql", dbUser + ":" + dbPassword + "@tcp(" + dbApiKey + ")/" + dbName)
				// db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:5221)/etiAssign")
				// handle error upon failure
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
				defer db.Close()
				
				//inserting values into passenger table
				_, err = db.Exec("insert into checkout (ShopCartID, EmailAddr, Shipping,PostalCode, CreditCard,TotalPayment) values(?,?,?,?,?,?)", checkout.ShopCartID, checkout.EmailAddr, checkout.Shipping,checkout.PostalCode, checkout.CreditCard, checkout.TotalPayment)
				//Handling error of SQL statement
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}

				//inserting values into passenger table
				_, err = db.Exec("update shopping_cart_user set IsCheckout = 1 where ShopCartID = ?", checkout.ShopCartID)
				//Handling error of SQL statement
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
			} else { fmt.Println(err)}
		}else { fmt.Println(err)}
	} else if r.Method == "GET"{
		querystringmap := r.URL.Query()
		userID := querystringmap.Get("ShopCartID")
		//Opening database connection
		db, err := sql.Open("mysql", dbUser + ":" + dbPassword + "@tcp(" + dbReadApiKey + ")/" + dbName)
		//db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:5221)/etiAssign")
		// handle error upon failure
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		defer db.Close()
		
		var checkoutCart []checkout
		results, err := db.Query("select * from checkout where ShopCartID = ?", userID)
		//Handling error of SQL statement
		if err != nil {
			http.Error(w, "Missing data", http.StatusBadRequest)
			panic(err.Error())
		}
		for results.Next() {
			var checkoutDetails checkout
			err = results.Scan( &checkoutDetails.ShopCartID, &checkoutDetails.EmailAddr, &checkoutDetails.Shipping,  &checkoutDetails.PostalCode, &checkoutDetails.CreditCard,&checkoutDetails.TotalPayment)
				if err != nil {
					http.Error(w, "Missing data", http.StatusBadRequest)
				} else {
					checkoutCart = append(checkoutCart, checkoutDetails)
				}
		}
		
		output, _ := json.Marshal(checkoutCart)
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, string(output))
	} else if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	} else{
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}

func shoppingCartItemEndpoint(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	dbApiKey := os.Getenv("API_KEY")
	dbReadApiKey := os.Getenv("READ_API_KEY")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")

	//Getting the items inside of the shopping cart
	if r.Method == "GET"{
		querystringmap := r.URL.Query()
		inputShopCartID := querystringmap.Get("ShopCartID")
		
		//Opening database connection
		db, err := sql.Open("mysql", dbUser + ":" + dbPassword + "@tcp(" + dbReadApiKey + ")/" + dbName)
		// handle error upon failure
		if err != nil {
			http.Error(w, "Unable to connect", http.StatusBadRequest)
		}
		defer db.Close()
		var ShoppingCart []shoppingCart
		cartItemResults, err := db.Query("select * from shopping_cart where ShopCartID = ?", inputShopCartID)
		//Handling error of SQL statement
		if err != nil {
			http.Error(w, "Missing data", http.StatusBadRequest)
			panic(err.Error())
		}
		for cartItemResults.Next() {
			var ShoppingCartItem shoppingCart
			err = cartItemResults.Scan( &ShoppingCartItem.ShopCartID, &ShoppingCartItem.ProductID, &ShoppingCartItem.Quantity )
				if err != nil {
					http.Error(w, "Missing data", http.StatusBadRequest)
				} else {
					ShoppingCart = append(ShoppingCart, ShoppingCartItem)
				}
		}

		output, _ := json.Marshal(ShoppingCart)
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, string(output))

	} else if r.Method =="POST"{
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var newShopItem shoppingCart
			// newShopItem.ShopCartID = 1
			// newShopItem.ProductID = 1
			// newShopItem.Quantity = 1
			if err := json.Unmarshal(body, &newShopItem); err == nil{
				//Opening database connection
				db, err := sql.Open("mysql", dbUser + ":" + dbPassword + "@tcp(" + dbApiKey + ")/" + dbName)
				// handle error upon failure
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
				defer db.Close()
				
				//inserting values into passenger table
				_, err = db.Exec("insert into shopping_cart (ShopCartID, productID, Quantity) values(?, ?, ?)", newShopItem.ShopCartID, newShopItem.ProductID, newShopItem.Quantity)
				//Handling error of SQL statement
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
				w.WriteHeader(http.StatusAccepted)
			}
		}
	} else {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

}
