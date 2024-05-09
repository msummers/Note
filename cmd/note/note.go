package main

import (
	"Note/pkg/model/note"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func main() {
	// TODO move gin and gorm init to config package
	// TODO abstract gin & gorm for mocking and testing
	// Init gin
	router := gin.Default()

	// TODO Add User model, login, and session context for security
	// TODO middleware to enforce user identification per REST call

	// Init db
	db, err := gorm.Open(sqlite.Open("notes.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	note.Init(router, db)

	err = router.Run()
	if err != nil {
		log.Fatalf("error running gin: %v", err)
	}
}
