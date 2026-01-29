package resources

import (
	"context"
	"sync"

	"connectrpc.com/connect"
	buildnetv5 "github.com/buildnethq/buildnet/hub/gen/buildnet/v5"
)

// Stub ResourceService (in-memory catalog + watch).
// Next step: SQLite catalog + worker discovery plugins.
type Service struct {
	mu        sync.RWMutex
	byID      map[string]*buildnetv5.ResourceDescriptor
	watchSubs map[chan *buildnetv5.ResourceEvent]struct{}
}

func New() *Service {
	return &Service{
		byID:      map[string]*buildnetv5.ResourceDescriptor{},
		watchSubs: map[chan *buildnetv5.ResourceEvent]struct{}{},
	}
}

func (s *Service) UpsertResources(ctx context.Context, req *connect.Request[buildnetv5.UpsertResourcesRequest]) (*connect.Response[buildnetv5.UpsertResourcesResponse], error) {
	_ = ctx
	accepted := 0

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, ev := range req.Msg.GetEvents() {
		switch ch := ev.Change.(type) {
		case *buildnetv5.ResourceEvent_Upsert:
			r := ch.Upsert
			if r == nil || len(r.ResourceId) == 0 {
				continue
			}
			s.byID[string(r.ResourceId)] = r
			accepted++
			s.broadcastLocked(ev)
		case *buildnetv5.ResourceEvent_TombstoneResourceId:
			delete(s.byID, string(ch.TombstoneResourceId))
			accepted++
			s.broadcastLocked(ev)
		}
	}

	return connect.NewResponse(&buildnetv5.UpsertResourcesResponse{
		Accepted: uint32(accepted),
		Rejected: 0,
	}), nil
}

func (s *Service) ListResources(ctx context.Context, req *connect.Request[buildnetv5.ListResourcesRequest]) (*connect.Response[buildnetv5.ListResourcesResponse], error) {
	_ = ctx
	typePrefix := req.Msg.GetTypeUriPrefix()
	provider := req.Msg.GetProviderNodeId()
	minVis := req.Msg.GetMinVisibility()

	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]*buildnetv5.ResourceDescriptor, 0, len(s.byID))
	for _, r := range s.byID {
		if typePrefix != "" && (len(r.TypeUri) < len(typePrefix) || r.TypeUri[:len(typePrefix)] != typePrefix) {
			continue
		}
		if len(provider) > 0 && string(r.ProviderNodeId) != string(provider) {
			continue
		}
		if minVis != buildnetv5.ResourceVisibility_RESOURCE_VISIBILITY_UNSPECIFIED && r.Visibility < minVis {
			continue
		}
		out = append(out, r)
	}

	return connect.NewResponse(&buildnetv5.ListResourcesResponse{Resources: out}), nil
}

func (s *Service) WatchResources(ctx context.Context, req *connect.Request[buildnetv5.WatchResourcesRequest], stream *connect.ServerStream[buildnetv5.ResourceEvent]) error {
	typePrefix := req.Msg.GetTypeUriPrefix()
	minVis := req.Msg.GetMinVisibility()

	ch := make(chan *buildnetv5.ResourceEvent, 64)

	s.mu.Lock()
	s.watchSubs[ch] = struct{}{}
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.watchSubs, ch)
		close(ch)
		s.mu.Unlock()
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev := <-ch:
			if typePrefix != "" {
				if up := ev.GetUpsert(); up != nil {
					if len(up.TypeUri) < len(typePrefix) || up.TypeUri[:len(typePrefix)] != typePrefix {
						continue
					}
				}
			}
			if minVis != buildnetv5.ResourceVisibility_RESOURCE_VISIBILITY_UNSPECIFIED {
				if up := ev.GetUpsert(); up != nil && up.Visibility < minVis {
					continue
				}
			}
			if err := stream.Send(ev); err != nil {
				return err
			}
		}
	}
}

func (s *Service) broadcastLocked(ev *buildnetv5.ResourceEvent) {
	for sub := range s.watchSubs {
		select {
		case sub <- ev:
		default:
			// drop if subscriber is slow
		}
	}
}
