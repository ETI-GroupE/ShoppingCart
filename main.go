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
	"github.com/joho/godotenv"
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
	ShopCartID   int    `json:"shopCartID"`
	TotalPayment int 	`json:"totalPayment"`
	EmailAddr	 string	`json:"emailAddr"`
	Shipping	 string `json:"shipping"`
	CreditCard	 string `json:"creditCard"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/shoppingCart", shoppingCartItemEndpoint	).Methods("GET","POST")
	router.HandleFunc("/shoppingCartUser", shoppingCartCreateEndpoint).Methods("GET","POST")
	router.HandleFunc("/checkout", checkoutEndpoint).Methods("POST")
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))

	// //loading environmental variables
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
}

func shoppingCartCreateEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		querystringmap := r.URL.Query()
		userID := querystringmap.Get("UserID")

		//Opening database connection
		db, err := sql.Open("mysql","root:root@tcp(127.0.0.1:5221)/etiassign")
		// handle error upon failure
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		defer db.Close()
		
		//inserting values into passenger table
		_, err = db.Exec("insert into shopping_cart_item (UserID) values(?)", userID)
		//Handling error of SQL statement
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusAccepted)
	} else if r.Method == "GET" {
		querystringmap := r.URL.Query()
		userID := querystringmap.Get("UserID")
		//Opening database connection
		db, err := sql.Open("mysql","root:root@tcp(127.0.0.1:5221)/etiassign")
		// handle error upon failure
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		defer db.Close()
		
		var shoppingCartUser shoppingCartUser
		results, err := db.Query("select * from shopping_cart_item where UserID = ?", userID)
		//Handling error of SQL statement
		if err != nil {
			http.Error(w, "Missing data", http.StatusBadRequest)
			panic(err.Error())
		}
		for results.Next() {
			err = results.Scan( &shoppingCartUser.ShopCartID, &shoppingCartUser.UserID, &shoppingCartUser.IsCheckout )
				if err != nil {
					http.Error(w, "Missing data", http.StatusBadRequest)
				} else {
					output, _ := json.Marshal(shoppingCartUser)
					w.WriteHeader(http.StatusAccepted)
					fmt.Fprintf(w, string(output))
				}
		}

	} else{
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}

func checkoutEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST"{
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var checkout checkout
			if err := json.Unmarshal(body, &checkout); err == nil{
				//Opening database connection
				db, err := sql.Open("mysql","root:root@tcp(127.0.0.1:5221)/etiassign")
				// handle error upon failure
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
				defer db.Close()
				
				//inserting values into passenger table
				_, err = db.Exec("insert into checkout (ShopCartID, UserID, EmailAddr, Shipping, CreditCard) values(?,?,?,?,?)", checkout.ShopCartID, checkout.EmailAddr, checkout.Shipping, checkout.CreditCard, checkout.TotalPayment)
				//Handling error of SQL statement
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
				w.WriteHeader(http.StatusAccepted)

				//inserting values into passenger table
				_, err = db.Exec("update checkout set IsCheckout = 1 where ShopCart_ID = ?", checkout.ShopCartID)
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

func shoppingCartItemEndpoint(w http.ResponseWriter, r *http.Request) {
	//Getting the items inside of the shopping cart
	if r.Method == "GET"{
		querystringmap := r.URL.Query()
		inputShopCartID := querystringmap.Get("ShopCartID")
		
		//Opening database connection
		db, err := sql.Open("mysql","root:root@tcp(127.0.0.1:5221)/etiassign")
		// handle error upon failure
		if err != nil {
			http.Error(w, "Unable to connect", http.StatusBadRequest)
		}
		defer db.Close()

		var ShoppingCart shoppingCart
		cartItemResults, err := db.Query("select * from shopping_cart where ShopCart_ID = ?", inputShopCartID)
		//Handling error of SQL statement
		if err != nil {
			http.Error(w, "Missing data", http.StatusBadRequest)
			panic(err.Error())
		}
		for cartItemResults.Next() {
			err = cartItemResults.Scan( &ShoppingCart.ShopCartID, &ShoppingCart.ProductID, &ShoppingCart.Quantity )
				if err != nil {
					http.Error(w, "Missing data", http.StatusBadRequest)
				} else {
					output, _ := json.Marshal(ShoppingCart)
					w.WriteHeader(http.StatusAccepted)
					fmt.Fprintf(w, string(output))
				}
		}

		} else if r.Method =="POST"{
			if body, err := ioutil.ReadAll(r.Body); err == nil {
				var newShopItem shoppingCart
				// newShopItem.ShopCartID = 1
				// newShopItem.ProductID = 1
				// newShopItem.Quantity = 1
				if err := json.Unmarshal(body, &newShopItem); err == nil{
					//Opening database connection
					db, err := sql.Open("mysql","root:root@tcp(127.0.0.1:5221)/etiassign")
					// handle error upon failure
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
					}
					defer db.Close()
					
					//inserting values into passenger table
					_, err = db.Exec("insert into shopping_cart (ShopCart_ID, product_ID, Quantity) values(?, ?, ?)", newShopItem.ShopCartID, newShopItem.ProductID, newShopItem.Quantity)
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
