package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	name := flag.String("name", "", "username (required)")
	peer := flag.String("peer", "", "peer address host:port to connect to (optional)")
	listen := flag.String("listen", ":50051", "listen address host:port (used when acting as server)")
	flag.Parse()

	if *name == "" {
		fmt.Fprintln(os.Stderr, "-name is required")
		os.Exit(2)
	}

	// Placeholder: networking will be implemented next.
	fmt.Printf("p2p-grpc-chat (stub) name=%s listen=%s peer=%s\n", *name, *listen, *peer)
}
