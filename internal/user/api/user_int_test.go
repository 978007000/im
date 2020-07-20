package api

import (
	"fmt"
	"im/pkg/pb"
	"testing"

	"google.golang.org/grpc"
)

func getUserIntClient() pb.UserIntClient {
	conn, err := grpc.Dial("localhost:50300", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return pb.NewUserIntClient(conn)
}

func TestUserIntServer_Auth(t *testing.T) {
	_, err := getUserIntClient().Auth(getCtx(), &pb.AuthReq{
		UserId:   3,
		DeviceId: 1,
		Token:    "0",
	})
	fmt.Println(err)
}

func TestUserIntServer_GetUsers(t *testing.T) {
	resp, err := getUserIntClient().GetUsers(getCtx(), &pb.GetUsersReq{
		UserIds: []int64{1, 2, 3},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	for k, v := range resp.Users {
		fmt.Printf("%+-5v  %+v\n", k, v)
	}
}
