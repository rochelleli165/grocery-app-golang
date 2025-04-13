package archive

import (
	"backend/main/config"
	"backend/main/models"
	"context"

	//"encoding/json"
	"fmt"
)

type FirebaseStore struct {
    Name     string    `json:"name"`
    Location string `json:"address"`
    RecentAd string `json:"recent-ad"`
}

func main2() {
    config.ConnectPostgreSQL()
    config.InitFirebase()
    config.InitLogger()
    ctx := context.Background()
    client, err := config.FirebaseApp.Database(ctx)
    if err != nil {
        fmt.Print(err)
        return
    }

    //  user := *models.NewUser(nil, "Rochelle Li", "rochelleli165@gmail.com")
    
    //  model := models.NewUserModel(config.PostgreSQL, *config.Logger)
    //  model.CreateUser(user)
    raw_stores := client.NewRef("stores")

    var result map[string]FirebaseStore
    if err := raw_stores.Get(ctx, &result); err != nil {
        fmt.Print(err)
        return
    }
    stores := []models.Store{}

    model := models.NewStoreModel(config.PostgreSQL, *config.Logger)

    for _, raw_item := range result {
        store := models.NewStore(nil, &raw_item.Name, raw_item.Location)
        stores = append(stores, *store)
    }
    
    for _, store := range stores {
        model.CreateStore(store)
    }
    
}
