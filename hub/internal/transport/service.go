package transport

import (
	"context"
	"errors"
	"io"
	"time"

	"connectrpc.com/connect"
	buildnetv5 "github.com/buildnethq/buildnet/hub/gen/buildnet/v5"
)

// Stub TransportService (PLAINTEXT framing only).
// Purpose: routing + stream plumbing; crypto core comes next.
type Service struct{}

func New() *Service { return &Service{} }

func (s *Service) OpenSession(ctx context.Context, stream *connect.BidiStream[buildnetv5.TransportFrame, buildnetv5.TransportFrame]) error {
	_ = ctx

	for {
		in, err := stream.Receive()
		if err != nil {
			// connect-go typically returns io.EOF when the client closes cleanly
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		if in == nil {
			continue
		}

		switch f := in.Frame.(type) {
		case *buildnetv5.TransportFrame_ClientHello:
			out := &buildnetv5.TransportFrame{
				Frame: &buildnetv5.TransportFrame_ServerHello{
					ServerHello: &buildnetv5.ServerHello{
						Server: &buildnetv5.NodeIdentity{
							NodeId:      []byte("stub-server"),
							IdentityPub: []byte("stub"),
							DhPub:       []byte("stub"),
						},
						SelectedSuite:   buildnetv5.CryptoSuite_SUITE_X25519_CHACHA20POLY1305_BLAKE3,
						SelectedPattern: buildnetv5.HandshakePattern_PATTERN_XX,
						ServerNonce:     []byte("stub-nonce"),
						EffectiveFlow:   f.ClientHello.GetRequestedFlow(),
					},
				},
			}
			if err := stream.Send(out); err != nil {
				return err
			}

		case *buildnetv5.TransportFrame_Envelope:
			// Echo envelope back (loopback test).
			if err := stream.Send(in); err != nil {
				return err
			}

		case *buildnetv5.TransportFrame_PingMs:
			if err := stream.Send(&buildnetv5.TransportFrame{
				Frame: &buildnetv5.TransportFrame_PongMs{PongMs: uint64(time.Now().UnixMilli())},
			}); err != nil {
				return err
			}

		default:
			_ = f
		}
	}
}
