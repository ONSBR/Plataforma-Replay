package broker

type Broker interface {
	Publish()
	Subscribe()
}
