package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tkircsi/pcbook/pb"
	"github.com/tkircsi/pcbook/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const (
	secretKey     = "MySecretKey"
	tokenDuration = 15 * time.Minute
	serverCert    = "cert/server-cert.pem"
	serverKey     = "cert/server-key.pem"
	caCert        = "cert/ca-cert.pem"
)

var (
	port       *int
	enableTls  *bool
	serverType *string
)

func main() {
	port = flag.Int("port", 0, "the server port")
	enableTls = flag.Bool("tls", false, "enable SSL/TLS")
	serverType = flag.String("type", "grpc", "type of server (grpc/rest/both)")
	flag.Parse()

	userStore := service.NewInMemoryUserStore()
	err := seedUsers(userStore)
	if err != nil {
		log.Fatal("cannot seed users")
	}

	jwtManager := service.NewJWTManager(secretKey, tokenDuration)
	authServer := service.NewAuthServer(userStore, jwtManager)

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore("img")
	ratingStore := service.NewInMemoryRatingStore()
	laptopServer := service.NewLaptopServer(laptopStore, service.WithImageStore(imageStore), service.WithRatingStore(ratingStore))

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server", err)
	}

	if *serverType == "grpc" {
		err = runGRPCServer(jwtManager, authServer, laptopServer, listener)
	} else if *serverType == "rest" {
		err = runRESTServer(jwtManager, authServer, laptopServer, listener)
	} else {
		// both
	}
	if err != nil {
		log.Fatal(err)
	}

}

func runGRPCServer(jwtManager *service.JWTManager, authServer *service.AuthServer, laptopServer *service.LaptopServer, listener net.Listener) error {

	interceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles())
	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	}

	if *enableTls {
		tlsCredentials, err := loadTLSCredentials()
		if err != nil {
			return fmt.Errorf("cannot load TLS credentials: %v", err)
		}
		serverOptions = append(serverOptions, grpc.Creds(tlsCredentials))
	}

	grpcServer := grpc.NewServer(serverOptions...)

	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	reflection.Register(grpcServer)

	log.Printf("gRPC server started at %s, TLS = %t", listener.Addr().String(), *enableTls)
	return grpcServer.Serve(listener)
}

func runRESTServer(jwtManager *service.JWTManager, authServer *service.AuthServer, laptopServer *service.LaptopServer, listener net.Listener) error {

	mux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := pb.RegisterAuthServiceHandlerServer(ctx, mux, authServer)
	if err != nil {
		return err
	}

	err = pb.RegisterLaptopServiceHandlerServer(ctx, mux, laptopServer)
	if err != nil {
		return err
	}

	log.Printf("REST server started at %s, TLS = %t", listener.Addr().String(), *enableTls)

	if *enableTls {
		return http.ServeTLS(listener, mux, serverCert, serverKey)
	}
	return http.Serve(listener, mux)
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load the certificate of the CA who signed the certificate of the client
	caCert, err := os.ReadFile(caCert)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to add client's CA certificate")
	}

	// Load server certificate and private key
	serverCert, err := tls.LoadX509KeyPair(serverCert, serverKey)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}

func accessibleRoles() map[string][]string {
	const laptopServicePath = "/pcbook.pb.LaptopService/"
	return map[string][]string{
		laptopServicePath + "CreateLaptop": {"admin"},
		laptopServicePath + "UploadImage":  {"admin"},
		laptopServicePath + "RateLaptop":   {"admin", "user"},
	}
}

func seedUsers(userStore service.UserStore) error {
	err := createUser(userStore, "admin1", "secret", "admin")
	if err != nil {
		return err
	}
	return createUser(userStore, "user1", "secret", "user")
}

func createUser(userStore service.UserStore, username, password, role string) error {
	user, err := service.NewUser(username, password, role)
	if err != nil {
		return err
	}

	return userStore.Save(user)
}
