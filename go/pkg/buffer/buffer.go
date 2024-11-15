package buffer

import "fmt"

type readRequest struct {
	k     string
	aswer chan (string)
}

type writeRequest struct {
	k, v string
}

type Storage struct {
	kv    map[string]string
	read  chan (readRequest)
	write chan (writeRequest)
	die   chan (struct{})
}

func New() *Storage {
	s := &Storage{
		kv:    map[string]string{},
		read:  make(chan readRequest),
		write: make(chan writeRequest),
		die:   make(chan struct{}),
	}
	go func() {
		for {
			select {
			case <-s.die:
				return
			case req := <-s.read:
				req.aswer <- s.kv[req.k]
			case req := <-s.write:
				s.kv[req.k] = req.v
			}
		}
	}()
	return s
}

func (s *Storage) Set(k, v string) {
	s.write <- writeRequest{k: k, v: v}
}

func (s *Storage) Get(k string) string {
	resp := make(chan string)
	s.read <- readRequest{
		k:     k,
		aswer: resp,
	}
	return <-resp
}

func (s *Storage) Die() {
	close(s.die)
}

func main() {
	s := New()
	s.Set("a", "b")
	fmt.Printf("%q %q\n", s.Get("a"), s.Get("q"))
	s.Die()
}
