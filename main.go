package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
)

type User struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

/*Local data*/
var users = []User{
	{1, "Anil", "Kumble", "Anl@gm.com"},
	{2, "john", "smith", "john@gm.com"},
	{3, "Adem", "gill kst", "kst@gm.com"},
}

func main() {
	initializeGinRouter()
}

/*----------------gin-----------------*/
// Gin function

func RunInfo(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "hello server running 8000 port",
		"status":  true,
	})
}

func GetUsers(context *gin.Context) {
	db := DBConnection()
	// Execute the query
	results, err := db.Query("SELECT * FROM user")
	db.Close()
	if err != nil {
		log.Fatal(err.Error()) // proper error handling instead of log in your app
	}
	var users []User
	for results.Next() {
		var user User
		// for each row, scan the result into our tag composite object
		err = results.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email)
		if err != nil {
			log.Print(err.Error()) // proper error handling instead of log in your app
		}
		// and then print out the tag's Name attribute
		users = append(users, user)
	}
	context.IndentedJSON(http.StatusOK, users)
}

func CreateUser(context *gin.Context) {
	var newUser User
	err := context.BindJSON(&newUser)
	if err != nil {
		context.String(http.StatusBadRequest, err.Error())
		return
	}
	db := DBConnection()
	insert, errIn := db.Query("INSERT INTO user VALUES (?, ?, ?, ?)", nil, newUser.FirstName, newUser.LastName, newUser.Email)
	insert.Close()
	db.Close()

	if errIn != nil {
		context.String(http.StatusForbidden, err.Error())
		return
	}
	context.IndentedJSON(http.StatusCreated, newUser)
}

func GetUser(context *gin.Context) {
	db := DBConnection()
	id, err := strconv.Atoi(context.Param("id"))
	//check if error occured
	if err != nil {
		//executes if there is any error
		db.Close()
		log.Print(err.Error())
		context.IndentedJSON(http.StatusBadRequest, "only accept numbers..")
		return
	}
	var user User
	// Execute the query
	err = db.QueryRow("SELECT * FROM user where id = ?", id).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email)
	db.Close()
	if err != nil {
		log.Print(err.Error())
		context.String(http.StatusForbidden, err.Error()) // proper error handling instead of panic in your app
		return
	} else {
		context.IndentedJSON(http.StatusOK, user)
		return
	}
	log.Print("data Not Found")
	context.IndentedJSON(http.StatusNotFound, "Invalid id check again...")
}

func Update(context *gin.Context) {
	var updatedUser User
	err := context.BindJSON(&updatedUser)
	if err != nil {
		context.String(http.StatusBadRequest, err.Error())
		return
	}
	db := DBConnection()
	insert, errIn := db.Query("UPDATE User SET first_name=?, last_name=?, email=? WHERE id=?; ", updatedUser.FirstName, updatedUser.LastName, updatedUser.Email, updatedUser.Id)
	insert.Close()
	db.Close()

	if errIn != nil {
		context.IndentedJSON(http.StatusForbidden, gin.H{
			"message": err.Error(),
			"status":  false,
		})
		return
	}
	context.IndentedJSON(http.StatusCreated, gin.H{
		"message": "Updated Success ...",
		"status":  true,
	})
}

func Delete(context *gin.Context) {
	db := DBConnection()
	id, err := strconv.Atoi(context.Param("id"))
	//check if error occured
	if err != nil {
		//executes if there is any error
		println(err)
		context.IndentedJSON(http.StatusBadRequest, "only accept numbers..")
		return
	}

	db.Exec("DELETE FROM user WHERE id= ?", id)
	db.Close()
	context.IndentedJSON(http.StatusOK, "Success")
}

//Gin

func initializeGinRouter() {

	r := gin.Default()
	r.GET("/", RunInfo)

	userRoute := r.Group("/users")
	{
		userRoute.GET("/getAll", GetUsers)
		userRoute.POST("/create", CreateUser)
		userRoute.GET("/get/:id", GetUser)
		userRoute.PUT("/update", Update)
		userRoute.DELETE("/delete/:id", Delete)
	}

	err := r.Run("localhost:8000")

	if err != nil {
		log.Fatal(err)
	}
}

func DBConnection() *sql.DB {
	// Open up our database connection.
	db, err := sql.Open("mysql", "root:1234@tcp(localhost:3307)/GoTextDb")

	// if there is an error opening the connection, handle it
	if err != nil {
		println("Connection false")
		log.Print(err.Error())
	} else {
		println("Connected")
	}
	return db
}

/*----------------gorilla-----------------*/
// function for gorilla mux
/*
func UpdateUser(writer http.ResponseWriter, request *http.Request) {

}


func DeleteUser(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	writer.WriteHeader(http.StatusOK)
	println(vars["id"])
	fmt.Fprintf(writer, "Id: %v\n", vars["id"])
}


func GetUser(writer http.ResponseWriter, request *http.Request) {

}


func CreateUser(writer http.ResponseWriter, request *http.Request) {

}

func GetUsers(writer http.ResponseWriter, request *http.Request) {
	var user User
	json.NewDecoder(request.Body).Decode(&user)
	fmt.Println(user)
}
*/
//gorilla mux
/*func initializeGorillaRouter()  {

	r := mux.NewRouter()
	r.HandleFunc("/users", GetUsers).Methods("GET")
	r.HandleFunc("/users", CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", GetUser).Methods("GET")
	r.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")
	r.HandleFunc("/users/{id}", UpdateUser).Methods("PUT")

	log.Fatal(http.ListenAndServe(":9000", r))
}
*/
