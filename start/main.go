package main

import (
	"backend/main/config"
	"backend/main/workflows"

	"context"

	"go.temporal.io/sdk/client"

	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

)

func runWorkflow() {
	fmt.Println("executing workflow at", time.Now())

	c, err := client.Dial(client.Options{})

	if err != nil {
		fmt.Println("Failed to create Temporal client", zap.Error(err))
		return
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:		"get_expired_ad_stores_workflow_" + string(time.Now().GoString()),
		TaskQueue: workflows.GetExpiredAdStoresTaskQueueName,
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, workflows.GetExpiredAdStores)

	if err != nil {
		fmt.Println("Failed to execute workflow", zap.Error(err))
		return
	}
	var expiredAdStores []workflows.AdProcessInput
	err = we.Get(context.Background(), &expiredAdStores)
	if err != nil {
		fmt.Println("Failed to get workflow result", zap.Error(err))
		return
	}
	options = client.StartWorkflowOptions{
		ID:		"ad_workflow_" + string(time.Now().GoString()),
		TaskQueue: workflows.AdProcessingTaskQueueName,
	}
	for _, v := range expiredAdStores {
		we, err = c.ExecuteWorkflow(context.Background(), options, workflows.AdProcess, v.StoreID, v.ZipCode, v.StoreFlippName)
		if err != nil {
			fmt.Println("Failed to execute workflow", zap.Error(err))
			return
		}
	}

	var result string

	err = we.Get(context.Background(), &result)

	if err != nil {
		fmt.Println("Failed to get workflow result", zap.Error(err))
		return
	}
	fmt.Println("Workflow result:", result)
}

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
	runWorkflow()
	
	c := cron.New()
	c.AddFunc("0 3 * * *", runWorkflow)
	c.Start()
	
	
	
	select {}
}