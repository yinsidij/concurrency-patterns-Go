package main

import (
	"pubsub"
    "fmt"
    "math/rand"
    "time"
)

func test(s pubsub.PubSub) {

	subscriptionChannels := make([]chan pubsub.Event, 10)
	for i := range subscriptionChannels {
		subscriptionChannels[i] = make(chan pubsub.Event)
		s.Subscribe(subscriptionChannels[i])
		fmt.Print("Channel ", i, " is subscribed by Server\n")
	}

	consumer := func(c chan pubsub.Event, consumerId int) {

		for event := range c {
			fmt.Print("event ", event.Id, " is consumed by consumer ", consumerId, "\n")
			delay := rand.Intn(100)
			if consumerId == 0 {
				// consumer 0 is very slow
				delay = 2000
			}

			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
	}

	for i := range subscriptionChannels{
		go consumer(subscriptionChannels[i], i)
	}

	publisher := func(publishCount int) {
		for i := 0; i < publishCount; i++{
			event := pubsub.Event{Id: i}
			s.Publish(event)
			fmt.Print("event ", i, " is published\n")
		}
	}

	go publisher(5)
	time.Sleep(1*time.Second)
	fmt.Print("Exit...\n")

}

func main() {
	s0 := pubsub.Server0{}
	s0.Init()
	fmt.Print("Server0 is initialized\n")
	fmt.Print("======================\n")
	s1 := pubsub.Server1{}
	s1.Init()
	fmt.Print("Server1 is initialized\n")
	fmt.Print("======================\n")
	s2 := pubsub.Server2{}
	s2.Init()
	fmt.Print("Server2 is initialized\n")
	fmt.Print("======================\n")
	s3 := pubsub.Server3{}
	s3.Init()
	fmt.Print("Server3 is initialized\n")
	fmt.Print("======================\n")
	test(&s0)
	test(&s1)
	test(&s2)
	test(&s3)
}