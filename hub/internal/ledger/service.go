package ledger

import (
	"context"
	"sync"
	"time"

	"connectrpc.com/connect"
	buildnetv5 "github.com/buildnethq/buildnet/hub/gen/buildnet/v5"
)

// Stub LedgerService (in-memory).
// Next step: SQLite persistence + signature/hash-chain validation.
type Service struct {
	mu     sync.RWMutex
	events []*buildnetv5.LedgerEvent
}

func New() *Service { return &Service{} }

func (s *Service) AppendEvents(ctx context.Context, req *connect.Request[buildnetv5.AppendEventsRequest]) (*connect.Response[buildnetv5.AppendEventsResponse], error) {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()

	evs := req.Msg.GetEvents()
	s.events = append(s.events, evs...)

	return connect.NewResponse(&buildnetv5.AppendEventsResponse{
		Accepted: uint32(len(evs)),
		Rejected: 0,
	}), nil
}

func (s *Service) StreamEvents(ctx context.Context, req *connect.Request[buildnetv5.StreamEventsRequest], stream *connect.ServerStream[buildnetv5.LedgerEvent]) error {
	_ = req
	s.mu.RLock()
	snap := append([]*buildnetv5.LedgerEvent(nil), s.events...)
	s.mu.RUnlock()

	for _, e := range snap {
		if err := stream.Send(e); err != nil {
			return err
		}
	}
	<-ctx.Done()
	return ctx.Err()
}

func (s *Service) GetSnapshot(ctx context.Context, req *connect.Request[buildnetv5.GetSnapshotRequest]) (*connect.Response[buildnetv5.LedgerSnapshot], error) {
	_ = ctx
	_ = req
	return connect.NewResponse(&buildnetv5.LedgerSnapshot{
		HeadByNode:    map[string][]byte{},
		GeneratedAtMs: uint64(time.Now().UnixMilli()),
	}), nil
}
