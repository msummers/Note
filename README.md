# Note
Simple RESTful CRUDL Note proof of concept that allows for multiple notes for multiple users

## Routes
```go
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

```