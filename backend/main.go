package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	Priority  int    `json:"priority"`
}

type addTodo struct {
	Title    string `json:"title"`
	Priority int    `json:"priority"`
}

type updateTodo struct {
	Completed bool `json:"completed"`
}

var todos = []todo{
	{ID: 1, Title: "Task 1", Completed: false, Priority: 2},
	{ID: 2, Title: "Task 2", Completed: false, Priority: 1},
	{ID: 3, Title: "Task 3", Completed: false, Priority: 3},
}

var g_id = 4

func getTodos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, todos)
}

func addTodos(c *gin.Context) {
	var newTodo addTodo

	if err := c.BindJSON(&newTodo); err != nil {
		return
	}

	todo := todo{
		ID:        g_id,
		Title:     newTodo.Title,
		Completed: false,
		Priority:  newTodo.Priority,
	}
	g_id += 1

	todos = append(todos, todo)
	c.IndentedJSON(http.StatusCreated, newTodo)
}

func getTodosByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 0, 0)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	for _, todo := range todos {
		if int64(todo.ID) == id {
			c.IndentedJSON(http.StatusOK, todo)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "todo not found"})
}

func deleteTodosByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 0, 0)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	for i, todo := range todos {
		if int64(todo.ID) == id {
			todos = append(todos[:i], todos[i+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "todo deleted"})
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "todo not found"})
}

func setTodosCompletedByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 0, 0)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	var updatedTodoData updateTodo

	if err := c.BindJSON(&updatedTodoData); err != nil {
		return
	}

	for i, todo := range todos {
		if int64(todo.ID) == id {
			todos[i].Completed = updatedTodoData.Completed
			c.IndentedJSON(http.StatusOK, gin.H{"message": "todo completed"})
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "todo not found"})
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("CORS Middleware")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	var router = gin.Default()
	router.Use(CORSMiddleware())

	router.GET("/todos", getTodos)
	router.POST("/todos", addTodos)
	router.GET("/todos/:id", getTodosByID)
	router.DELETE("/todos/:id", deleteTodosByID)
	router.PATCH("/todos/:id", setTodosCompletedByID)

	router.Run("localhost:3000")
}
