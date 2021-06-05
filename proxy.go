package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/atreya2011/grpc-postgres-crud/postgrescrud"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// end point means the ip and port on which grpc server is listening
var postgrescrudEndPoint = flag.String("endpoint", "localhost:7000", "endpoint of PostgresCrudService")

func run() error {
	// Step 1. define a context with cancel
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Step 2. define the serve mux options here
	muxOpts := []runtime.ServeMuxOption{
		// Todo: check meaning of MIMEWildcard
		// OrigName set to false, json field names will be camelCase
		// EmitDefaults set to true, empty json fields will not be omitted
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: false, EmitDefaults: true}),
	}

	// Step 3. create a new server mux
	mux := runtime.NewServeMux(muxOpts...)
	// Step 4. create grpc dial options for connecting to grpc server
	opts := []grpc.DialOption{grpc.WithInsecure()}
	// Step 5. register the handler from endpoint (function in generated pb.gw.go file)
	// the following are passed into the function
	// context, new server mux, endpoint string, dial options
	err := postgrescrud.RegisterPostgresCrudHandlerFromEndpoint(ctx, mux, *postgrescrudEndPoint, opts)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Listening...")
	// Step 6. Start API Server with ListenAndServe, pass in the serve mux
	// We launch a goroutine for handling requests
	done := make(chan struct{})
	go func() {
		http.ListenAndServe(":8082", mux)
		done <- struct{}{}
	}()

	<-done //wait finish groutine
	log.Println("Program exit")
	return nil
}

func main() {
	// flag.Parse parses the flags passed in the command line
	// if no flag is present, default values are used
	// if there is no flag.Parse, command line flags are not parsed
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatalln(err)
	}
}
