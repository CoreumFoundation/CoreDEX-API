package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	grpcMetadata "google.golang.org/grpc/metadata"

	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

type GRPCClient struct {
	token *oauth2.Token
	Conn  *grpc.ClientConn
}

/*
Initialize the client.
Depending on the parameter, the environment is determined to be either in cluster of local by:
localhost:port => local (and expecting a simple TLS less port to work with (Hint: start other required services also on local))
localhost => No port is not local
*/
func InitClient(endpoint string) *GRPCClient {
	var conn *grpc.ClientConn
	contactServerHost := os.Getenv(endpoint)
	if contactServerHost == "" {
		logger.Fatalf("%s is not set", endpoint)
	}

	if !strings.Contains(contactServerHost, ":") {
		// Cloud run connect:
		// We expect the service name without the project dependent addition GRPC_APPEND
		// Get the GRPC_APPEND env var:
		grpcAppend := os.Getenv("GRPC_APPEND")
		if grpcAppend == "" {
			logger.Fatalf("GRPC_APPEND is not set, expecting -abc")
		}
		contactServerHost = contactServerHost + grpcAppend
		address := contactServerHost + ":443"

		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			logger.Fatalf("failed to load system root CA certs: %v", err)
		}
		cred := credentials.NewTLS(&tls.Config{
			RootCAs: systemRoots,
		})

		conn, err = grpc.NewClient(address, grpc.WithAuthority(contactServerHost),
			grpc.WithTransportCredentials(cred))
		if err != nil {
			logger.Fatalf(errors.Wrap(err, fmt.Sprintf("failed to init %s client", endpoint)).Error())
		}
	} else {
		var err error
		conn, err = grpc.NewClient(contactServerHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatalf(errors.Wrap(err, fmt.Sprintf("failed to init %s client", endpoint)).Error())
		}
	}
	g := &GRPCClient{
		token: initAuthContext(contactServerHost),
		Conn:  conn,
	}
	go g.authTokenContextRefresh(contactServerHost)
	return g
}

/*
Auth context is valid for a maximum of 1 hour after initialization.
(That can be altered, however this solution would always be required: If you set it for a year and the service
runs for a year, then the token would still invalidate and still require a refresh)
So this process runs in a loop and refreshes the token every 30 minutes (conservative, we do not want edge cases).
*/
func (g *GRPCClient) authTokenContextRefresh(contactServerHost string) {
	for {
		time.Sleep(30 * time.Minute)
		g.token = initAuthContext(contactServerHost)
	}
}

func initAuthContext(contactServerHost string) *oauth2.Token {
	ctx := context.Background()
	if !strings.Contains(contactServerHost, ":") {
		tokenSource, err := idtoken.NewTokenSource(ctx, "https://"+contactServerHost)
		if err != nil {
			logger.Fatalf("Can not get tokenSource: %v", err)
		}
		token, err := tokenSource.Token()
		if err != nil {
			logger.Fatalf("Can not get token from tokenSource: %v", err)
		}
		return token
	}
	return nil
}

// Returns a IAM enhanced GRPC context for accessing the cloudrun services
func (g *GRPCClient) AuthCtx(ctx context.Context) context.Context {
	if g.token == nil {
		return ctx
	}
	return grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+g.token.AccessToken)
}
