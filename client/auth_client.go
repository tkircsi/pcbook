package client

import (
	"context"
	"time"

	"github.com/tkircsi/pcbook/pb"
	"google.golang.org/grpc"
)

// AuthClient is a client to call authentication RPC
type AuthClient struct {
	service  pb.AuthServiceClient
	username string
	password string
}

// NewAuthClient initiaélizes a new AuthClient object
func NewAuthClient(cc *grpc.ClientConn, username, password string) *AuthClient {
	return &AuthClient{
		service:  pb.NewAuthServiceClient(cc),
		username: username,
		password: password,
	}
}

func (c *AuthClient) Login() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.LoginRequest{
		Username: c.username,
		Password: c.password,
	}

	res, err := c.service.Login(ctx, req)
	if err != nil {
		return "", err
	}

	return res.GetAccessToken(), nil
}
