package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type User struct{
	Id			int64  `json:"id"`
	FirstName	string `json:"first_name"`
	LastName	string `json:"last_name"`
	Email		string `json:"email"`
}

/*Local data*/
var users = [] User {
	{1,"Anil", "Kumble", "Anl@gm.com"},
	{2,"john", "smith", "john@gm.com"},
	{3,"Adem", "gill kst", "kst@gm.com"},
}



func main() {
	initializeGinRouter()
}
/*----------------gorilla-----------------*/
// Gin function

func RunInfo(context *gin.Context) {
	context.JSON(http.StatusOK,gin.H{
		"message": "hello server running 8000 port",
		"status": true,
	})
}

func GetUsers(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, users)
}

func CreateUser(context *gin.Context) {
	var newUser User

	err := context.BindJSON(&newUser)
	if err != nil {
		context.String(http.StatusBadRequest, err.Error())
		return
	}

	newUser.Id = int64(len(users)+1)
	users = append(users, newUser)
	context.IndentedJSON(http.StatusCreated, newUser)
}


func GetUser(context *gin.Context) {
	id , err := strconv.Atoi(context.Param("id"))
	//check if error occured
	if err != nil{
		//executes if there is any error
		println(err)
		context.IndentedJSON(http.StatusBadRequest,"only accept numbers..")
		return

	}
	for i:=range users{
		println(users[i].Id== int64(id))
		if int64(id) == users[i].Id {
			context.IndentedJSON(http.StatusFound, users[i])
			return
		}
	}
	context.IndentedJSON(http.StatusNotFound,"Invalid id check again...")
}


func Update(context *gin.Context) {
	var updatedUser User

	err := context.BindJSON(&updatedUser)
	if err != nil {
		context.String(http.StatusBadRequest, err.Error())
		return
	}else {
		for i:=range users{
			if updatedUser.Id == users[i].Id {
				users[i] =updatedUser
				context.IndentedJSON(http.StatusOK, updatedUser)
				return
			}
		}
	}
	context.IndentedJSON(http.StatusNotFound, "User Notfound :"+updatedUser.FirstName+" "+updatedUser.LastName)
}

func Delete(context *gin.Context) {
	id , err := strconv.Atoi(context.Param("id"))
	//check if error occured
	if err != nil{
		//executes if there is any error
		println(err)
		context.IndentedJSON(http.StatusBadRequest,"only accept numbers..")
		return
	}

	if len(users) >0{
		for i:=range users{
			if int64(id) == users[i].Id {
				users = RemoveIndex(1+i)
				context.IndentedJSON(http.StatusOK, users)
				return
			}
		}
	}
	context.IndentedJSON(http.StatusNotFound, "User Notfound :")
}


func RemoveIndex(index int) []User {
	return append(users[:index-1], users[index:]...)
}

//Gin

func initializeGinRouter()  {

	r := gin.Default()
	r.GET("/", RunInfo)

	userRoute:= r.Group("/users")
	{
		userRoute.GET("/getAll", GetUsers)
		userRoute.POST("/", CreateUser)
		userRoute.GET("/get/:id", GetUser)
		userRoute.PUT("/update", Update)
		userRoute.DELETE("/delete/:id", Delete)
	}



	err:= r.Run("localhost:8000")

	if err != nil {
		log.Fatal(err)
	}
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
