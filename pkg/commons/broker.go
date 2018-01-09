package commons

// Publisher publishes WebResoure
type Publisher interface {
	Publish(wr *WebResource) (success bool)
}

// Subscriber pushes PubSub events
type Subscriber interface {
	Subscribe(out chan *WebResource)
}
