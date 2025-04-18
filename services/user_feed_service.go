package services

import (
	"fmt"
	"context"
	"backend/main/pb"
	"backend/main/models"
)

type UserFeedService struct {
	pb.UnimplementedUserFeedServiceServer
	UserModel *models.UserModel
	StoreModel *models.StoreModel
	AdModel *models.AdModel
	IngredientModel *models.IngredientModel
}

func NewUserFeedService(userModel *models.UserModel, storeModel *models.StoreModel, adModel *models.AdModel, ingredientModel *models.IngredientModel) *UserFeedService {
	return &UserFeedService{
		UserModel: userModel,
		StoreModel: storeModel,
		AdModel: adModel,
		IngredientModel: ingredientModel,
	}
}

// GetUserAds retrieves ads for a user
func (s (*UserFeedService)) GetUserAds(ctx context.Context, req *pb.GetUserAdsRequest) (*pb.GetUserAdsResponse, error) {
	fmt.Println("GetUserAds called with UserId:", req.UserId)
	stores, err := s.StoreModel.GetSubscribedStores(uint(req.UserId))
	if err != nil {
		return nil, err
	}
	var ads []models.Ad
	for _, store := range stores {
		ad, err := s.AdModel.GetRecentAd(*store.ID)
		if err != nil {
			fmt.Println("Error retrieving ad for store:", store.Name, "Error:", err)
			return nil, err
		}
		ads = append(ads, ad)
	}
	var adList []*pb.Ad
	
	for i, ad := range ads {
		var adItemsList []*pb.AdItemData
		for _, adItem := range ad.Ingredient {
			ingredient, err := s.IngredientModel.GetIngredientByID(uint(adItem.IngredientID))
			if err != nil {
				fmt.Println("Error retrieving ingredient for ad item:", adItem.Name, "Error:", err)
				return nil, err
			}
			adItemsList = append(adItemsList, &pb.AdItemData{
				Ingredient: ingredient.Name,
				Name: adItem.Name,
				Price: adItem.Price,
				Sale: adItem.Sale,
				IngredientType: ingredient.Type.String(),
			})
		}
		adList = append(adList, &pb.Ad{
			StoreName: *stores[i].Name,
			StoreAddress: stores[i].Location,
			AdItems: adItemsList,
		})
	}
	fmt.Println("Ads retrieved successfully for UserId:", req.UserId)
	return &pb.GetUserAdsResponse{Ads: adList}, nil
}