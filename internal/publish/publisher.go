package publish

type Publisher interface {
	Publish(msg string) error
}
