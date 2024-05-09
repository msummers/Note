package note

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Note struct {
	gorm.Model
	// TODO User model
	User  string `json:"user"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

var db *gorm.DB

func Init(router *gin.Engine, gormDb *gorm.DB) {
	// Setup gorm
	db = gormDb
	err := db.AutoMigrate(&Note{})
	if err != nil {
		log.Fatalf("auto migrate note err: %v", err)
	}

	// Setup gin routes
	// Create
	router.POST("/note/:user", create)

	// List
	router.GET("/note/:user", list)

	// Read
	router.GET("/note/:user/:id", read)

	// Update
	router.POST("/note/:user/:id", update)

	// Delete
	router.DELETE("/note/:user/:id", remove)

}

func create(ctx *gin.Context) {
	var note = &Note{}
	// Unmarshal the Note
	err := ctx.ShouldBind(note)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note.User = ctx.Param("user")
	result := db.Create(&note)
	if hasError(result, ctx) {
		return
	}
	ctx.IndentedJSON(http.StatusOK, note)
}

func remove(ctx *gin.Context) {
	result := db.Where("id = ? AND user= ?", ctx.Param("id"), ctx.Param("user")).Delete(&Note{})
	if hasError(result, ctx) {
		return
	}
	ctx.IndentedJSON(http.StatusOK, Note{})
}

func list(ctx *gin.Context) {
	notes := &[]Note{}
	result := db.Where("user = ?", ctx.Param("user")).Find(notes)
	if hasError(result, ctx) {
		return
	}
	ctx.IndentedJSON(http.StatusOK, notes)
}

func read(ctx *gin.Context) {
	note, result := getNote(ctx)
	if hasError(result, ctx) {
		return
	}
	ctx.IndentedJSON(http.StatusOK, note)
}

// getNote DRYs the read function for read & update
func getNote(ctx *gin.Context) (note *Note, result *gorm.DB) {
	note = &Note{}
	result = db.Where("user = ? AND id = ?", ctx.Param("user"), ctx.Param("id")).First(note)
	return note, result
}

func update(ctx *gin.Context) {
	// Get the original Note
	original, result := getNote(ctx)
	if hasError(result, ctx) {
		return
	}

	// Unmarshal the updated Note
	var updates Note
	err := ctx.ShouldBind(&updates)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	applyUpdates(&updates, original)

	// Save the update
	result = db.Save(original)
	if hasError(result, ctx) {
		return
	}
	ctx.IndentedJSON(http.StatusOK, updates)
}

func applyUpdates(src *Note, dst *Note) {
	// TODO refactor using meta programming
	if src.Title != "" {
		dst.Title = src.Title
	}
	if src.Body != "" {
		dst.Body = src.Body
	}
}

func hasError(result *gorm.DB, ctx *gin.Context) bool {
	// TODO refactor to use more specific Status codes
	if result.Error != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return true
	}
	if result.RowsAffected == 0 {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "note not found"})
		return true
	}
	return false
}
