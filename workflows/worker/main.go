package main

import (
	"fmt"
	"backend/main/workflows"
	"backend/main/config"
	"context"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)


func main() {
	config.ConnectPostgreSQL()
	config.InitFirebase()
	config.InitLogger()
	config.InitGenAI()
	ctx := context.Background()
	_, err := config.FirebaseApp.Database(ctx)
	if err != nil {
		fmt.Print(err)
		return
	}

	c, err := client.Dial(client.Options{})
	if err != nil {
		fmt.Println("Failed to create Temporal client", err)
	}
	defer c.Close()

	w := worker.New(c, workflows.AdProcessingTaskQueueName, worker.Options{})
	w0 := worker.New(c, workflows.GetExpiredAdStoresTaskQueueName, worker.Options{})

	w.RegisterWorkflow(workflows.AdProcess)
	w0.RegisterWorkflow(workflows.GetExpiredAdStores)

	w.RegisterActivity(workflows.GetFlippData)
	w.RegisterActivity(workflows.RetrieveTranslations)
	w.RegisterActivity(workflows.GetIngredientNamesAndIds)
	w.RegisterActivity(workflows.AddTranslations)
	w.RegisterActivity(workflows.CreateAd)
	w.RegisterActivity(workflows.CreateNewIngredientByName)
	w0.RegisterActivity(workflows.GetExpiredAdStores)

	err = w0.Run(worker.InterruptCh())
	if err != nil {
		fmt.Println("Failed to start worker", err)
	}
	err = w.Run(worker.InterruptCh())
	if err != nil {
		fmt.Println("Failed to start worker", err)
	}
} 