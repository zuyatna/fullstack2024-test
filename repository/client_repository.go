package repository

import (
	"context"
	"fmt"
	"fullstack2024-test/database"
	"fullstack2024-test/model"

	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"
)

type IClientRepository interface {
	CreateClient(client *model.Client) (*model.Client, error)
	GetClientByID(id int) (*model.Client, error)
	UpdateClientByID(id int, client *model.Client) (*model.Client, error)
	DeleteClientByID(id int) error
}

type ClientRepository struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewClientRepository(db *gorm.DB) *ClientRepository {
	return &ClientRepository{
		db:          db,
		redisClient: database.GetRedisClient(),
	}
}

func (r *ClientRepository) CreateClient(client *model.Client) (*model.Client, error) {
	if err := r.db.Create(client).Error; err != nil {
		return nil, err
	}

	ctx := context.Background()
	if err := database.StoreJSON(ctx, client.Slug, client); err != nil {

	}

	return client, nil
}

func (r *ClientRepository) GetClientByID(id int) (*model.Client, error) {
	client := &model.Client{}
	if err := r.db.First(client, id).Error; err != nil {
		return nil, err
	}

	return client, nil
}

func (r *ClientRepository) UpdateClientByID(id int, client *model.Client) (*model.Client, error) {
	existingClient := &model.Client{}
	if err := r.db.First(existingClient, id).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&model.Client{}).Where("id = ?", id).Updates(client).Error; err != nil {
		return nil, err
	}

	updatedClient := &model.Client{}
	if err := r.db.First(updatedClient, id).Error; err != nil {
		return nil, err
	}

	ctx := context.Background()

	if existingClient.Slug != updatedClient.Slug {
		if err := database.DeleteKey(ctx, existingClient.Slug); err != nil {
			fmt.Println("Error deleting old key:", err)
		}
	}

	if err := database.StoreJSON(ctx, updatedClient.Slug, updatedClient); err != nil {
		fmt.Println("Error storing updated client:", err)
	}

	return updatedClient, nil
}

func (r *ClientRepository) DeleteClientByID(id int) error {
	client := &model.Client{}
	if err := r.db.First(client, id).Error; err != nil {
		return err
	}

	if err := r.db.Delete(&model.Client{}, id).Error; err != nil {
		return err
	}

	ctx := context.Background()
	if err := database.DeleteKey(ctx, client.Slug); err != nil {
		fmt.Println("Error deleting client from Redis:", err)
	}

	return nil
}
