package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"

	buildnetv5 "github.com/buildnethq/buildnet/hub/gen/buildnet/v5"
	buildnetv5connect "github.com/buildnethq/buildnet/hub/gen/buildnet/v5/buildnetv5connect"
)

func main() {
	baseURL := "http://127.0.0.1:7444"

	// IMPORTANT: default http.Client does NOT do h2c prior-knowledge.
	// For Connect bidi streaming on cleartext, we need an http2.Transport with AllowHTTP + custom dial.
	h2cTransport := &http2.Transport{
		AllowHTTP: true,
		DialTLSContext: func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, network, addr)
		},
	}

	httpClient := &http.Client{Transport: h2cTransport}

	client := buildnetv5connect.NewTransportServiceClient(httpClient, baseURL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream := client.OpenSession(ctx)

	// Send ClientHello
	if err := stream.Send(&buildnetv5.TransportFrame{
		Frame: &buildnetv5.TransportFrame_ClientHello{
			ClientHello: &buildnetv5.ClientHello{
				Client: &buildnetv5.NodeIdentity{
					NodeId:      []byte("smoke-client"),
					IdentityPub: []byte("stub"),
					DhPub:       []byte("stub"),
				},
				Offer: &buildnetv5.SuiteOffer{
					SupportedSuites:   []buildnetv5.CryptoSuite{buildnetv5.CryptoSuite_SUITE_X25519_CHACHA20POLY1305_BLAKE3},
					SupportedPatterns: []buildnetv5.HandshakePattern{buildnetv5.HandshakePattern_PATTERN_XX},
				},
				ClientNonce: []byte("nonce"),
			},
		},
	}); err != nil {
		panic(fmt.Errorf("Send(ClientHello): %w", err))
	}

	msg, err := stream.Receive()
	if err != nil {
		panic(fmt.Errorf("Receive(): %w", err))
	}

	sh := msg.GetServerHello()
	if sh == nil {
		fmt.Printf("Got non-ServerHello frame: %#v\n", msg)
		return
	}

	fmt.Println("OK: received ServerHello")
	fmt.Printf("  selectedSuite=%v\n", sh.SelectedSuite)
	fmt.Printf("  selectedPattern=%v\n", sh.SelectedPattern)
	fmt.Printf("  server.nodeId=%q\n", string(sh.GetServer().GetNodeId()))
}
