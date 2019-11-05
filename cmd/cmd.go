// Package cmd exports a function to create easily a CLI based on cobra.
package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/danielkvist/botio/client"
	"github.com/danielkvist/botio/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Root creates a root *cobra.Command and then adds to it
// the received *cobra.Commands, then it executes the root command
// returning an error if any.
func Root(commands ...*cobra.Command) error {
	examples := []string{
		"botio server bolt --database ./data/commands.db --collection commands --http :9090 --key mysupersecretkey",
		"botio bot --platform telegram --token <telegram-token> --url :9090 --token <jwt-token>",
		"botio client print --command start --url :9090 --jwt <jwt-token>",
	}

	root := &cobra.Command{
		Use:          "botio",
		Short:        "Botio is a simple and opinionated CLI to create and manage easily bots for differents platforms.",
		Example:      strings.Join(examples, "\n"),
		SilenceUsage: true,
	}

	for _, cmd := range commands {
		root.AddCommand(cmd)
	}

	return root.Execute()
}

func checkURL(url string, prefix bool, suffix bool) (string, error) {
	if url == "" {
		return "", fmt.Errorf("server URL cannot be an empty string")
	}

	if !strings.HasPrefix(url, "https://") && prefix {
		url = "https://" + url
	}

	if !strings.HasSuffix(url, "/api/commands") && suffix {
		url = url + "/api/commands"
	}

	return url, nil
}

func getClient(url string, server string, crt string, key string, ca string) (client.Client, error) {
	u, err := checkURL(url, false, false)
	if err != nil {
		return nil, fmt.Errorf("while parsing URL: %v", err)
	}

	cert, err := tls.LoadX509KeyPair(crt, key)
	if err != nil {
		return nil, fmt.Errorf("while loading client SSL key pair: %v", err)
	}

	certPool := x509.NewCertPool()
	caCert, err := ioutil.ReadFile(ca)
	if err != nil {
		return nil, fmt.Errorf("while reading CA certificate: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		return nil, fmt.Errorf("faile to append CA certificates")
	}

	creds := credentials.NewTLS(&tls.Config{
		ServerName:   server,
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	})

	conn, err := grpc.Dial(u, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("while creating a new Dial for %q: %v", u, err)
	}

	c := client.New(u, conn)
	return c, nil
}

func printCommand(cmd *proto.BotCommand) {
	fmt.Printf("%q: %q\n", cmd.GetCmd().GetCommand(), cmd.GetResp().GetResponse())
}
