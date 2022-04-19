package resolver

import (
	"math/rand"
	"sync"
)

type IDResolver interface {
	SetID(originalID []byte, clientID uint64) ID
	PopID(requestID uint64) (ID, bool)
}

type idResolver struct {
	pool sync.Map
}

type ID struct {
	RequestID  uint64
	OriginalID []byte
	ClientID   uint64
}

func NewIDResolver() IDResolver {
	return &idResolver{}
}

func (r *idResolver) SetID(originalID []byte, clientID uint64) ID {
	id := ID{RequestID: rand.Uint64(), OriginalID: originalID, ClientID: clientID}

	r.pool.Store(id.RequestID, id)
	return id
}

func (r *idResolver) PopID(requestID uint64) (ID, bool) {
	id, ok := r.pool.LoadAndDelete(requestID)
	if !ok {
		return ID{}, ok
	}

	return id.(ID), ok
}
