package models

import (
	"context"

	"go.uber.org/zap"
	"github.com/jackc/pgx/v5"

)

// UserModel defines a struct for user service
type UserModel struct {
	PostgreSQL *pgx.Conn
	Logger    zap.Logger
}

// User defines a struct for user data
type User struct {
	ID       *uint      `json:"id"`
	Name     string    `json:"name"`
	Email string `json:"email"`
}

func NewUserModel(PostgreSQL *pgx.Conn, logger zap.Logger) *UserModel {
	return &UserModel{
		PostgreSQL: PostgreSQL,
		Logger: logger,
	}
}

// Constructor for User
func NewUser(id *uint, name string, email string) *User {
	return &User{
		ID:       id,
		Name:     name,
		Email: email,
	}
}

// GetUserByID to find user by ID from database
func (i *UserModel) GetUserByID(id uint) (User, error) {
	var user User
	err := i.PostgreSQL.QueryRow(context.Background(), 
		"SELECT id, name, email FROM grocery_user WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		i.Logger.Error("Error getting user by ID", zap.Error(err))
		return User{}, err
	}
	return user, nil
}

// CreateUser to add user to database
func (i *UserModel) CreateUser(user User) error {
	query := `INSERT INTO grocery_user (name, email) VALUES (@userName, @userEmail)`
	args := pgx.NamedArgs{
		"userName": user.Name,
		"userEmail": user.Email,
	  }
	_, err := i.PostgreSQL.Exec(context.Background(), query, args)
	if err != nil {
		i.Logger.Error("Error adding user to database", zap.Error(err))
		return err
	}
	return nil
}

// AddUserSubscriptionStore to add a user's subscription to a grocery store
func (i *UserModel) AddUserSubscriptionStore(userID uint, storeID uint) error {
	_, err := i.PostgreSQL.Exec(context.Background(), 
		"INSERT INTO store_subscription (user_id, store_id) VALUES ($1, $2)", userID, storeID)
	if err != nil {
		i.Logger.Error("Error adding user subscription to store", zap.Error(err))
		return err
	}
	return nil
}