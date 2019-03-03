// go get github.com/gin-gonic/gin - for REST server library
// go get github.com/lib/pq - PostgreSQL driver
// To run, enter "go run main.go".
package main

import (
	"database/sql" // to open database connection
	"errors"
	"fmt"
	"log"
	"net/http" // for status constants
	"strconv"  // to convert between string and int values

	"github.com/gin-gonic/gin" // HTTP web framework
	_ "github.com/lib/pq"      // Postgres driver
)

const allowOrigin = "http://localhost:8080"
const badRequest = http.StatusBadRequest
const forbidden = http.StatusForbidden
const ok = http.StatusOK
const port = 1919
const serverError = http.StatusInternalServerError

// Dog describes a dog.
// We don't want uppercase names in JSON that is produced,
// so alternate names are provided using struct "tags".
type Dog struct {
	ID    int    `json:"id"`
	Breed string `json:"breed"`
	Name  string `json:"name"`
}

func shouldAllow(c *gin.Context) bool {
	origin := c.Request.Header["Origin"][0]
	return origin == allowOrigin
}

// Custom middleware to enable CORS
func cors(c *gin.Context) {
	if shouldAllow(c) {
		c.Header("Access-Control-Allow-Origin", allowOrigin)
	} else {
		c.Status(forbidden)
	}
}

func options(c *gin.Context) {
	c.Header("Access-Control-Allow-Methods", "DELETE,GET,POST,PUT")
	// Must explicitly allow Content-Type header for JSON bodies.
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	c.Status(ok)
}

func handleError(c *gin.Context, statusCode int, err error) {
	c.String(statusCode, err.Error())
}

func main() {
	// Connect to database.
	connStr := "user=postgres dbname=survey sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Configure HTTP request routes.
	router := gin.Default()
	router.Use(cors)

	// Required for preflight OPTIONS request before POST.
	router.OPTIONS("/dog", options)

	// Required for preflight OPTIONS request before PUT.
	router.OPTIONS("/dog/:id", options)

	// Heartbeat
	router.GET("/", func(c *gin.Context) {
		c.String(ok, "I'm alive!")
	})

	// Creates a dog.
	router.POST("/dog", func(c *gin.Context) {
		// Get dog from request body.
		var dog Dog
		if err := c.ShouldBindJSON(&dog); err != nil {
			handleError(c, badRequest, err)
			return
		}

		// Insert dog into database, getting assigned id.
		sql := fmt.Sprintf(
			"insert into dog (breed, name) values ('%s', '%s') returning id",
			dog.Breed,
			dog.Name)
		var id int
		err := db.QueryRow(sql).Scan(&id)
		if err != nil {
			handleError(c, serverError, err)
			return
		}

		dog.ID = id
		c.JSON(ok, dog)
	})

	// Retrieves all the dogs.
	router.GET("/dog", func(c *gin.Context) {
		if !shouldAllow(c) {
			c.Status(forbidden)
			return
		}

		rows, err := db.Query("select id, breed, name from dog")
		if err != nil {
			c.String(serverError, err.Error())
			return
		}

		defer rows.Close()

		dogs := []Dog{}
		var id int
		var breed, name string

		for rows.Next() {
			if err := rows.Scan(&id, &breed, &name); err != nil {
				c.String(serverError, err.Error())
				return
			}
			dogs = append(dogs, Dog{id, breed, name})
		}

		c.JSON(ok, dogs)
	})

	// Updates a dog.
	router.PUT("/dog/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			handleError(c, badRequest, errors.New("id must be int"))
			return
		}

		// Get dog from request body.
		var dog Dog
		if err := c.ShouldBindJSON(&dog); err != nil {
			handleError(c, badRequest, err)
			return
		}

		sql := fmt.Sprintf(
			"update dog set breed='%s', name='%s' where id=%d",
			dog.Breed,
			dog.Name,
			id)
		if _, err := db.Query(sql); err != nil {
			handleError(c, serverError, err)
			return
		}

		c.Status(ok)
	})

	// Deletes a dog.
	router.DELETE("/dog/:id", func(c *gin.Context) {
		id, e := strconv.Atoi(c.Param("id"))
		if e != nil {
			handleError(c, badRequest, errors.New("id must be int"))
			return
		}

		sql := fmt.Sprintf("delete from dog where id=%d", id)
		if _, err := db.Query(sql); err != nil {
			handleError(c, serverError, err)
			return
		}

		c.Status(ok)
	})

	fmt.Printf("listening on port %d\n", port)
	router.Run(":" + strconv.Itoa(port))
}
