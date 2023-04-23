package publish

import (
	"fmt"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type twilioPublisher struct {
	client   *twilio.RestClient
	sendFrom string
	sendTo   string
}

type NewTwilioPublisherParams struct {
	Client   *twilio.RestClient
	SendFrom string
	SendTo   string
}

func NewTwilioPublisher(params *NewTwilioPublisherParams) *twilioPublisher {
	return &twilioPublisher{
		client:   params.Client,
		sendFrom: params.SendFrom,
		sendTo:   params.SendTo,
	}
}

func (p twilioPublisher) Publish(msg string) error {
	to := fmt.Sprintf("whatsapp:%s", p.sendTo)
	from := fmt.Sprintf("whatsapp:%s", p.sendFrom)

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(from)
	params.SetBody(msg)

	_, err := p.client.Api.CreateMessage(params)
	return err
}

var _ Publisher = twilioPublisher{}
