package subscribe

type Subscriber interface {
	Subscribe() (<-chan string, error)
}
