package main

import (
	"context"
	"fmt"
	"os"

	p "github.com/pulumi/pulumi-go-provider"
	qeycloak "github.com/raushan606/pulumi-qeycloak/provider"
)

// A provider is a program that listens for requests from the Pulumi engine
// to interact with cloud providers using a CRUD-based model.
func main() {

	err := p.RunProvider(context.Background(), qeycloak.Name, qeycloak.Version, qeycloak.Provider())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running provider: %v\n", err)
		os.Exit(1)
	}
}
