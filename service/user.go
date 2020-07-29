package service

import (
	"context"
	pb "grpc-gateway-test/model"
	"github.com/mitchellh/mapstructure"
	"log"
	"grpc-gateway-test/driver"
)

type UserRpcService struct{}

func (s *UserRpcService) Get(ctx context.Context, req *pb.UserKeys) (*pb.UserSettings, error) {
	log.Printf("User Get called")
	client := driver.GetFireStoreClient(ctx)
	defer driver.CloseFireStoreClient(client)

	docRef := client.Collection("users").Doc(req.Username).Collection("profile").Doc("settings")
	doc, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}
	data := doc.Data()
	params := new(pb.UserSettings)
	err = mapstructure.Decode(data, params)
	if err != nil {
		return nil, err
	}
	return params, nil
}
