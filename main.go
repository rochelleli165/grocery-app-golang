package main

import (
	"backend/main/config"
	"backend/main/routes"
	"context"

	"fmt"
	"net/http"
)

func main() {
	config.ConnectPostgreSQL()
	config.InitFirebase()
	config.InitLogger()
	config.InitGenAI()
	ctx := context.Background()
	_, err := config.FirebaseApp.Database(ctx)
	if err != nil {
		return
	}

	router := routes.RegisterRoutes()
	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("Failed to start server", err)
		return
	}
	
}