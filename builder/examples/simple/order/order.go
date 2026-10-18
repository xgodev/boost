// Package order keeps the orders of the example in memory.
package order

import (
	"strconv"
	"sync"
)

// Order is an order of one item.
type Order struct {
	ID   string `json:"id"`
	Item string `json:"item"`
}

// Service stores the orders.
type Service struct {
	mu     sync.Mutex
	orders map[string]Order
}

// NewService provides the order service.
// @Provide(index=0)
func NewService() *Service {
	return &Service{orders: map[string]Order{}}
}

// Create stores a new order of item.
func (s *Service) Create(item string) Order {
	s.mu.Lock()
	defer s.mu.Unlock()

	o := Order{ID: strconv.Itoa(len(s.orders) + 1), Item: item}
	s.orders[o.ID] = o
	return o
}

// Get returns the order with id.
func (s *Service) Get(id string) (Order, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	o, ok := s.orders[id]
	return o, ok
}
