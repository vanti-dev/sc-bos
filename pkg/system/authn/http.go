package authn

import (
	"net/http"
	"sync"
)

// nextOrNotFound calls next.ServerHTTP if next is not nil, otherwise http.NotFound.
type nextOrNotFound struct {
	mu   sync.Mutex
	next http.Handler
}

func (p *nextOrNotFound) Next(next http.Handler) {
	p.mu.Lock()
	p.next = next
	p.mu.Unlock()
}

func (p *nextOrNotFound) Clear() {
	p.Next(nil)
}

func (p *nextOrNotFound) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	p.mu.Lock()
	next := p.next
	p.mu.Unlock()

	if next == nil {
		http.NotFound(writer, request)
		return
	}
	next.ServeHTTP(writer, request)
}
