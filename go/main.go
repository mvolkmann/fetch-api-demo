// go get github.com/gin-gonic/gin - for REST server library
// go get github.com/codegangsta/gin - for live reload of Go servers
// go get github.com/lib/pq - PostgreSQL driver
// To run, enter "gin run main.go".
package main

import (
	"database/sql" // to open database connection
	"fmt"
	"net/http" // just for status constants
	"strconv"  // convert between string and int values

	"github.com/gin-gonic/gin" // HTTP web framework
	_ "github.com/lib/pq"      // Postgres driver
)

const badRequest = http.StatusBadRequest
const ok = http.StatusOK
const port = 1919
const serverError = http.StatusInternalServerError

// Dog describes a dog.
type Dog struct {
	ID    int    `json:"id"`
	Breed string `json:"breed"`
	Name  string `json:"name"`
}

// Custom middleware to enable Cross-Origin Resource Sharing (CORS).
func cors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET,HEAD,PUT,POST,DELETE")
	c.Header("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding")
}

func handleError(c *gin.Context, statusCode int, err error) {
	handleErrorMsg(c, statusCode, err.Error())
}

func handleErrorMsg(c *gin.Context, statusCode int, msg string) {
	c.String(statusCode, msg)
	//TODO: Use c.Error(err)?
}

func main() {
	// Connect to database.
	connStr := "user=postgres dbname=survey sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		//TODO: Why do "go vet" want a return value to be captured?
		_ = fmt.Errorf("Error opening database connection: %s\n", err.Error())
		return
	}

	// Configure HTTP request routes.
	router := gin.Default()
	router.Use(cors)

	// Heartbeat
	router.GET("/", func(c *gin.Context) {
		c.String(ok, "I'm alive!")
	})

	// Required for preflight OPTIONS request before POST to /dog.
	router.OPTIONS("/dog", func(c *gin.Context) {
		c.Status(ok)
	})

	// Required for preflight OPTIONS request before PUT to /dog.
	router.OPTIONS("/dog/:id", func(c *gin.Context) {
		c.Status(ok)
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
			handleErrorMsg(c, badRequest, "id must be int")
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
			handleErrorMsg(c, badRequest, "id must be int")
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
