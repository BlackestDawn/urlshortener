package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func respondJSON(c *gin.Context, code int, data any) {
	c.Header("Content-Type", "application/json")

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error encoding JSON: %s\n", err)
		c.String(http.StatusInternalServerError, "an error occurd")
	}

	c.String(code, string(jsonData))
}

func respondJSONError(c *gin.Context, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}

	if code >= 500 {
		log.Printf("internal error: %s\n", msg)
	}

	type errVal struct {
		Error string
	}

	respondJSON(c, code, errVal{Error: msg})
}
