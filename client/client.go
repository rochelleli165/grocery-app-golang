package main

import (
	"context"
	"log"
	"time"
	"fmt"

	"google.golang.org/grpc"
	"backend/main/pb"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserFeedServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.GetUserAdsRequest{UserId: 2}
	res, err := client.GetUserAds(ctx, req)
	if err != nil {
		log.Fatalf("Error calling GetUser: %v", err)
	}

	fmt.Println("Response:", res)
	for _, ad := range res.Ads {
		log.Printf("Ad Item: Ingredient=%s, Name=%s, Price=%f, Sale=%f\n", ad.AdItems[0].Ingredient, ad.AdItems[0].Name, ad.AdItems[0].Price, ad.AdItems[0].Sale)
	}
}
