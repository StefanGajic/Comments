package main

import "fmt"

// App - contains pointers to db connections
type App struct{}

// Run sets application
func (a *App) Run() error {
	fmt.Println("Setting app")
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
