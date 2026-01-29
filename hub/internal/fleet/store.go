package fleet

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	buildnetv5 "github.com/buildnethq/buildnet/hub/gen/buildnet/v5"
)

type Record struct {
	Info       *buildnetv5.WorkerInfo
	APIKey     string
	LastSeenMs uint64
}

type Store struct {
	mu      sync.RWMutex
	workers map[string]*Record
	ttlMs   uint64
	// alpha epoch: unix millis
}

func NewStore(leaseTTL time.Duration) *Store {
	return &Store{
		workers: make(map[string]*Record),
		ttlMs:   uint64(leaseTTL.Milliseconds()),
	}
}

func (s *Store) LeaseTTLms() uint64 { return s.ttlMs }

func (s *Store) Join(info *buildnetv5.WorkerInfo) (workerID string, apiKey string, serverEpoch uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	workerID = info.GetWorkerId()
	apiKey = newAPIKey(32)
	now := uint64(time.Now().UTC().UnixMilli())

	s.workers[workerID] = &Record{
		Info:       info,
		APIKey:     apiKey,
		LastSeenMs: now,
	}
	return workerID, apiKey, now
}

func (s *Store) Heartbeat(workerID, apiKey string) (ok bool, serverEpoch uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := uint64(time.Now().UTC().UnixMilli())
	rec, exists := s.workers[workerID]
	if !exists {
		return false, now
	}
	if rec.APIKey != apiKey {
		return false, now
	}
	rec.LastSeenMs = now
	return true, now
}

func (s *Store) List() (out []*buildnetv5.WorkerRecord, serverEpoch uint64) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := uint64(time.Now().UTC().UnixMilli())
	out = make([]*buildnetv5.WorkerRecord, 0, len(s.workers))
	for _, rec := range s.workers {
		out = append(out, &buildnetv5.WorkerRecord{

