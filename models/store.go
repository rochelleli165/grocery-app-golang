package models

import (
	"context"

	"go.uber.org/zap"
	"github.com/jackc/pgx/v5"

)

// StoreModel defines a struct for Store service
type StoreModel struct {
	PostgreSQL *pgx.Conn
	Logger    zap.Logger
}

// Ingredient defines a struct for ingredient data
type Store struct {
	ID       *uint      `json:"id"`
	Name     *string    `json:"name"`
	Location string `json:"location"`
	FlippMerchantName string `json:"flipp_merchant_name"`
}

func NewStoreModel(PostgreSQL *pgx.Conn, logger zap.Logger) *StoreModel {
	return &StoreModel{
		PostgreSQL: PostgreSQL,
		Logger: logger,
	}
}

// Constructor for Ingredient
func NewStore(id *uint, name *string, location string) *Store {
	return &Store{
		ID:       id,
		Name:     name,
		Location: location,
	}
}

// GetStoreByID to find Store by ID from database
func (i *StoreModel) GetStoreByID(id uint) (Store, error) {
	var store Store
	err := i.PostgreSQL.QueryRow(context.Background(), 
		"SELECT id, name, location FROM store WHERE id = $1", id).Scan(&store.ID, &store.Name, &store.Location)
	if err != nil {
		i.Logger.Error("Error getting store by ID", zap.Error(err))
		return Store{}, err
	}
	return store, nil
}

// CreateStore to add Store to database
func (i *StoreModel) CreateStore(store Store) error {
	query := `INSERT INTO store (name, location) VALUES (@StoreName, @StoreLocation)`
	args := pgx.NamedArgs{
		"StoreName": store.Name,
		"StoreLocation": store.Location,
	  }
	_, err := i.PostgreSQL.Exec(context.Background(), query, args)
	if err != nil {
		i.Logger.Error("Error adding store to database", zap.Error(err))
		return err
	}
	return nil
}

func (i *StoreModel) GetExpiredAdStores() (stores []Store, err error) {
	query := "select store_id, max(sale_end), flipp_merchant, location from ad join store on store_id = store.id group by store_id, flipp_merchant, location having max(sale_end) < current_date;"
	rows, err := i.PostgreSQL.Query(context.Background(),
	query)
	if err != nil {
		i.Logger.Error("Error getting stores with expired ads", zap.Error(err), zap.String("function", "GetExpiredAdStores"))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var s Store
		err := rows.Scan(&s.ID, &s.Location, &s.FlippMerchantName)
		if err != nil {
			i.Logger.Error("Error scanning row", zap.Error(err), zap.String("function", "GetExpiredAdStores"))
			return nil, err
		}
		stores = append(stores, s)
	}
	if err := rows.Err(); err != nil {
		i.Logger.Error("Error processing rows", zap.Error(err), zap.String("function", "GetExpiredAdStores"))
		return nil, err
	}
	return stores, nil
}

// SubscribeStore to add to store_subscription
func (i *StoreModel) SubscribeStore(userID uint, storeID uint) error {
	query := `INSERT INTO store_subscription (user_id, store_id) VALUES ($1, $2)`
	_, err := i.PostgreSQL.Exec(context.Background(), query, userID, storeID)
	if err != nil {
		i.Logger.Error("Error subscribing to store", zap.Error(err))
		return err
	}
	return nil
}

// UnsubscribeStore to remove from store_subscription
func (i *StoreModel) UnsubscribeStore(userID uint, storeID uint) error {
	query := `DELETE FROM store_subscription WHERE user_id = $1 AND store_id = $2`
	_, err := i.PostgreSQL.Exec(context.Background(), query, userID, storeID)
	if err != nil {
		i.Logger.Error("Error unsubscribing from store", zap.Error(err))
		return err
	}
	return nil
}

// GetSubscribedStores to get all stores subscribed by user given userID
func (i *StoreModel) GetSubscribedStores(userID uint) ([]Store, error) {
	var stores []Store
	rows, err := i.PostgreSQL.Query(context.Background(),
		"SELECT s.id, s.name, s.location FROM store_subscription ss INNER JOIN store s ON ss.store_id = s.id WHERE ss.user_id = $1", userID)
	if err != nil {
		i.Logger.Error("Error getting subscribed stores", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var store Store
		err := rows.Scan(&store.ID, &store.Name, &store.Location)
		if err != nil {
			i.Logger.Error("Error scanning subscribed store(s)", zap.Error(err))
			return nil, err
		}
		stores = append(stores, store)
	}

	return stores, nil
}


