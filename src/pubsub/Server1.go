package pubsub

type Server1 struct {
	publish   chan Event
	subscribe chan subReq
	cancel    chan subReq
}

type subReq struct {
	c chan<- Event
	ok chan bool
}

func (s *Server1)Init(){
	s.publish = make(chan Event)
	s.subscribe = make(chan subReq)
	s.cancel = make(chan subReq)
	go s.loop()
}

func (s *Server1) loop(){
	sub := make(map[chan<- Event]bool)

	for{
		select {
		case e := <-s.publish:
			for c := range sub {
				c <- e
			}
		case req := <-s.subscribe:
			if sub[req.c] {
				req.ok <- false
				// only break select
				break
			}
			sub[req.c] = true
			req.ok <- true
		case req := <-s.cancel:
			if !sub[req.c] {
				req.ok <- false
				break
			}
			close(req.c)
			delete(sub, req.c)
			req.ok <- true
		}
	}
}

func (s *Server1)Publish(e Event){
	s.publish <- e
}

func (s *Server1)Subscribe(c chan<- Event){
	req := subReq{
		c:  c,
		ok: make(chan bool),
	}
	s.subscribe <- req
	if ! <- req.ok{
		panic("pubsub: already subscribed")
	}
}

func (s *Server1)Cancel(c chan<- Event){
	req := subReq{
		c:  c,
		ok: make(chan bool),
	}
	s.cancel <- req
	if ! <- req.ok{
		panic("pubsub: not subscribed")
	}
}