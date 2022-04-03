package main

import (
	"fmt"
	"net/http"

	"github.com/Comments/internal/comment"
	"github.com/Comments/internal/database"
	transportHTTP "github.com/Comments/internal/transport/http"
)

// App - contains pointers to db connections
type App struct{}

// Run sets application
func (a *App) Run() error {
	fmt.Println("Setting app")

	var err error
	db, err := database.NewDatabase()
	if err != nil {
		return nil
	}
	err = database.MigrateDB(db)
	if err != nil {
		return err
	}

	CommentService := comment.NewService(db)

	handler := transportHTTP.NewHandler(CommentService)
	handler.SetupRoutes()

	if err := http.ListenAndServe(":8080", handler.Router); err != nil {
		fmt.Println("Failed to start server")
		return err
	}

	return nil
}

func main() {
	fmt.Println("Go REST API !!!")
	app := App{}
	if err := app.Run(); err != nil {
		fmt.Println("Error starting application!")
		fmt.Println(err)
	}
}
