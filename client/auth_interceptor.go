package client

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AuthInterceptor is a client interceptor for authentication
type AuthInterceptor struct {
	authClient  *AuthClient
	authMethods map[string]bool
	accessToken string
}

func NewAuthInterceptor(
	authClient *AuthClient,
	authMethods map[string]bool,
	refreshDuration time.Duration,
) (*AuthInterceptor, error) {
	interceptor := &AuthInterceptor{
		authClient:  authClient,
		authMethods: authMethods,
	}

	err := interceptor.scheduleRefreshToken(refreshDuration)
	if err != nil {
		return nil, err
	}
	return interceptor, nil
}

func (ai *AuthInterceptor) scheduleRefreshToken(refreshDuration time.Duration) error {
	err := ai.refreshToken()
	if err != nil {
		return err
	}
	go func() {
		wait := refreshDuration
		tk := time.NewTicker(wait)
		for {
			<-tk.C
			// <-time.Tick(wait)
			err := ai.refreshToken()
			if err != nil {
				tk.Reset(time.Second)
			} else {
				tk.Reset(wait)
			}
		}
	}()
	return nil
}

func (ai *AuthInterceptor) refreshToken() error {
	accessToken, err := ai.authClient.Login()
	if err != nil {
		return err
	}
	ai.accessToken = accessToken
	log.Printf("token refreshed: %v", accessToken)
	return nil
}

func (ai *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req,
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		log.Printf("--> unary interceptor: %s", method)

		if ai.authMethods[method] {
			return invoker(ai.attachToken(ctx), method, req, reply, cc, opts...)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (ai *AuthInterceptor) Stream() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		log.Printf("--> stream interceptor: %s", method)

		if ai.authMethods[method] {
			return streamer(ai.attachToken(ctx), desc, cc, method, opts...)
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}

func (ai *AuthInterceptor) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", ai.accessToken)
}
