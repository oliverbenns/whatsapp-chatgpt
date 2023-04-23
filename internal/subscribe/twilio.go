package subscribe

type twilioSubscriber struct {
}

type NewTwilioSubscriberParams struct {
}

func NewTwilioSubscriber(params *NewTwilioSubscriberParams) *twilioSubscriber {
	return &twilioSubscriber{}
}

func (s twilioSubscriber) Subscribe() error {
	return nil
}

var _ Subscriber = twilioSubscriber{}
