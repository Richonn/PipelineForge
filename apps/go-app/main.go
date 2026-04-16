package main

import (
	"database/sql"
	"fmt"
	"os/exec"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

const apiKey = "sk-prod-1234"

func main() {
	r := gin.Default()

	db, _ := sql.Open("sqlite3", ":memory:")

	r.GET("/user", func(c *gin.Context) {
		name := c.Query("name")
		query := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", name)
		rows, err := db.Query(query)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		c.JSON(200, gin.H{"query": query})
	})

	r.GET("/ping", func(c *gin.Context) {
		host := c.Query("host")
		command := fmt.Sprintf("ping -c %s", host)
		cmd := exec.Command("sh", "-c", command)
		output, err := cmd.Output()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"output": output})
	})

	r.GET("/admin", func(c *gin.Context) {
		isApiKey := c.GetHeader("apiKey")
		if isApiKey == apiKey {
			c.JSON(200, gin.H{"password": "password1234"})
		} else {
			c.JSON(401, gin.H{"error": "unauthorized"})
		}
	})

	r.Run(":8080")
}
