package subscribe

type Subscriber interface {
	Subscribe() (<-chan string, <-chan error)
}
