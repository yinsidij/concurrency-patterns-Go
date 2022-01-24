package pubsub

import "fmt"

type Server3 struct {
	publish   chan Event
	subscribe chan subReq
	cancel    chan subReq
}

/*
Deal with the slow subscribers, keep the overall program from blocking. 
Main loop goroutine will send the events to the helper3.

If the input channel closes, it indicates that queue is also empty so we can exit
*/
func helper3(in <-chan Event, out chan<- Event){
	var q []Event
	for in != nil || len(q) > 0{
		// Decide whether and what to send.
		var sendOut chan<- Event
		var next Event
		if len(q) > 0{
			sendOut = out
			next = q[0]
		}
		select {
		case e, ok := <-in:
			if !ok{
				in = nil // stop receiving from in
				fmt.Println("helper3 stops")
				break
			}
			q = append(q, e)
			fmt.Println("event ", e.Id, " is queued")
		case sendOut <- next:
			q = q[:1]
			fmt.Println("event ", next.Id, " is out")
		}
	}
	close(out)
}

func (s *Server3)Init(){
	s.publish = make(chan Event)
	s.subscribe = make(chan subReq)
	s.cancel = make(chan subReq)
	go s.loop()
}

func (s *Server3) loop(){
	//map from subscribe channel to helper3 in channel
	sub := make(map[chan<- Event]chan<- Event)
	for{
		select {
		case e := <-s.publish:
			for _, helperIn := range sub {
				helperIn <- e
			}

		case req := <-s.subscribe:
			if sub[req.c] != nil {
				req.ok <- false
				break
			}
			helperIn := make(chan Event)
			go helper3(helperIn, req.c)
			sub[req.c] = helperIn
			req.ok <- true

		case req := <-s.cancel:
			if sub[req.c] != nil{
				req.ok <- false
				break
			}
			close(sub[req.c])
			delete(sub, req.c)
			req.ok <- true
		}
	}
}

func (s *Server3)Publish(e Event){
	s.publish <- e
}

func (s *Server3)Subscribe(c chan<- Event){
	req := subReq{
		c:  c,
		ok: make(chan bool),
	}
	s.subscribe <- req
	if ! <- req.ok{
		panic("pubsub: already subscribed")
	}
}

func (s *Server3)Cancel(c chan<- Event){
	req := subReq{
		c:  c,
		ok: make(chan bool),
	}
	s.cancel <- req
	if ! <- req.ok{
		panic("pubsub: not subscribed")
	}
}

