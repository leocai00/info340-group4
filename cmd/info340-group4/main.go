package main

// These are your imports / libraries / frameworks
import (
	// this is Go's built-in sql library
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	// this allows us to run our web server
	"github.com/gin-gonic/gin"
	// this lets us connect to Postgres DB's
	_ "github.com/lib/pq"
)

var (
	// this is the pointer to the database we will be working with
	// this is a "global" variable (sorta kinda, but you can use it as such)
	db *sql.DB
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	var errd error
	// here we want to open a connection to the database using an environemnt variable.
	// This isn't the best technique, but it is the simplest one for heroku
	db, errd = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if errd != nil {
		log.Fatalf("Error opening database: %q", errd)
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("html/*")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/account.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "account.html", nil)
	})

	router.GET("/newaccount.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "newaccount.html", nil)
	})

	router.GET("/QuserInfo", func(c *gin.Context) {
		rows, err := db.Query("SELECT first_name, last_name, email, phone_number FROM Customer WHERE customer_id = 1;")
		if err != nil {
			// careful about returning errors to the user!
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		var para string
		var first string
		var last string
		var email string
		var phone string 

		for rows.Next() {
			// assign each of them, in order, to the parameters of rows.Scan.
			// preface each variable with &
			rows.Scan(&first, &last, &email, &phone)
			// can't combine ints and strings in Go. Use strconv.Itoa(int) instead
			para += "<p>Name: " + first + " " + last + "</p>"
			para += "<p>Email: " + email + "</p>"
			para += "<p>Phone Number: " + phone + "</p>"
		}
		c.Data(http.StatusOK, "text/html", []byte(para))
	})

	router.GET("/QuserAddr", func(c *gin.Context) {
		rows, err := db.Query("SELECT address, city_name, state_name, zip_code FROM customer_address JOIN city ON city.city_id = customer_address.city_id JOIN state ON state.state_id = customer_address.state_id JOIN zip ON zip.zip_code_id = customer_address.zip_code_id WHERE customer_id = 1;")
		if err != nil {
			// careful about returning errors to the user!
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		
		var para string
		var address string
		var state_name string
		var city_name string
		var zip_code int 

		for rows.Next() {
			// assign each of them, in order, to the parameters of rows.Scan.
			// preface each variable with &
			rows.Scan(&address, &city_name, &state_name, &zip_code)
			// can't combine ints and strings in Go. Use strconv.Itoa(int) instead
			para += "<p>Address: " + address + " " + city_name + ", " + state_name + strconv.Itoa(zip_code) + "</p>"
		}
		c.Data(http.StatusOK, "text/html", []byte(para))
	})

	router.POST("/Qnewaccount", func(c *gin.Context) {
		fname := c.PostForm("fname")
		lname := c.PostForm("lname")
		email := c.PostForm("email")
		phone := c.PostForm("phone")
		password := c.PostForm("password")
		db.Query("SELECT create_customer($5, $1, $2, $3, $4);", fname, lname, email, phone, password)
	})

	// NO code should go after this line. it won't ever reach that point
	router.Run(":" + port)
}