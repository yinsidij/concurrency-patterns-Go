package pubsub

import "sync"

type Server0 struct {
	mu sync.Mutex
  	// key of map represents the channel 
  	// used to receive the event.
	sub map[chan <- Event]bool
}

func (s *Server0)Init(){
	s.sub = make(map[chan<-Event]bool)
}

func (s *Server0)Publish(e Event){
	s.mu.Lock()
	defer s.mu.Unlock()

	for c := range s.sub{
		c <- e
	}
}

func (s *Server0)Subscribe(c chan<- Event){
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.sub[c]{
		panic("pubsub: already subscribed")
	}
	s.sub[c] = true
}

func (s *Server0)Cancel(c chan<- Event){
	s.mu.Lock()
	s.mu.Unlock()

	if !s.sub[c]{
		panic("pubsub: not subscribed")
	}
	close(c)
	delete(s.sub, c)
}
