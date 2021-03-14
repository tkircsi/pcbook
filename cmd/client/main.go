package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/tkircsi/pcbook/client"
	"github.com/tkircsi/pcbook/pb"
	"github.com/tkircsi/pcbook/samples"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	username        = "admin1"
	password        = "secret"
	refreshDuration = 30 * time.Second
)

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials")
	}

	log.Printf("dial server %s", *serverAddress)

	cc1, err := grpc.Dial(
		*serverAddress,
		grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	authClient := client.NewAuthClient(cc1, username, password)
	//interceptor, err := client.NewAuthInterceptor(authClient, authMethods(), refreshDuration)
	interceptor, err := client.NewAuthInterceptor(authClient, client.WithAuthMethods(authMethods()), client.WithTokenRefreshDuration(refreshDuration))
	if err != nil {
		log.Fatal("cannot create auth interceptor: ", err)
	}

	cc2, err := grpc.Dial(
		*serverAddress,
		grpc.WithTransportCredentials(tlsCredentials),
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
	)
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}
	laptopClient := client.NewLaptopClient(cc2)
	//testUploadImage(laptopClient)
	testRateLaptop(laptopClient)

}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load the certificate of the CA who signed the certificate of the server
	caCert, err := os.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to add server's CA certificate")
	}

	// Load client certificate and private key
	clientCert, err := tls.LoadX509KeyPair("cert/client-cert.pem", "cert/client-key.pem")
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(config), nil
}

func authMethods() map[string]bool {
	const laptopServicePath = "/pcbook.pb.LaptopService/"
	return map[string]bool{
		laptopServicePath + "CreateLaptop": true,
		laptopServicePath + "UploadImage":  true,
		laptopServicePath + "RateLaptop":   true,
	}
}

func testCreateLpatop(laptopClient *client.LaptopClient) {
	laptopClient.CreateLaptop(samples.NewLaptop())
}

func testSearchLaptop(laptopClient *client.LaptopClient) {
	for i := 0; i < 10; i++ {
		laptopClient.CreateLaptop(samples.NewLaptop())
	}

	filter := &pb.Filter{
		MaxPriceUsd: 3000,
		MinCpuCores: 4,
		MinCpuGhz:   2.5,
		MinRam:      &pb.Memory{Value: 8, Unit: pb.Memory_GIGYBYTE},
	}

	laptopClient.SearchLaptop(filter)
}

func testUploadImage(laptopClient *client.LaptopClient) {
	laptop := samples.NewLaptop()
	laptopClient.CreateLaptop(laptop)
	laptopClient.UploadImage(laptop.GetId(), "tmp/laptop.jpg")
}

func testRateLaptop(laptopClient *client.LaptopClient) {
	n := 3
	laptopIDs := make([]string, n)

	for i := 0; i < n; i++ {
		laptop := samples.NewLaptop()
		laptopIDs[i] = laptop.GetId()
		laptopClient.CreateLaptop(laptop)
	}

	scores := make([]float64, n)
	for {
		fmt.Print("rate laptop? (y/n) ")
		var answer string
		fmt.Scan(&answer)

		if strings.ToLower(answer) != "y" {
			break
		}

		for i := 0; i < n; i++ {
			scores[i] = samples.RandomLaptopScore()
		}

		err := laptopClient.RateLaptop(laptopIDs, scores)
		if err != nil {
			log.Fatal(err)
		}
	}
}
