package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/tkircsi/pcbook/pb"
	"github.com/tkircsi/pcbook/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const (
	secretKey     = "MySecretKey"
	tokenDuration = 15 * time.Minute
)

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	log.Printf("server started on port %d", *port)

	userStore := service.NewInMemoryUserStore()
	err := seedUsers(userStore)
	if err != nil {
		log.Fatal("cannot seed users")
	}

	jwtManager := service.NewJWTManager(secretKey, tokenDuration)
	authServer := service.NewAuthServer(userStore, jwtManager)
	interceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles())

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore("img")
	ratingStore := service.NewInMemoryRatingStore()
	laptopServer := service.NewLaptopServer(laptopStore, service.WithImageStore(imageStore), service.WithRatingStore(ratingStore))

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials", err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	)
	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	reflection.Register(grpcServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server", err)
	}

}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load the certificate of the CA who signed the certificate of the client
	caCert, err := os.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to add client's CA certificate")
	}

	// Load server certificate and private key
	serverCert, err := tls.LoadX509KeyPair("cert/server-cert.pem", "cert/server-key.pem")
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
