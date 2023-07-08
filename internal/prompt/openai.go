package prompt

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type openAiPrompter struct {
	client *openai.Client
}

type NewOpenAiPrompterParams struct {
	Client *openai.Client
}

func NewOpenAiPrompter(params *NewOpenAiPrompterParams) *openAiPrompter {
	return &openAiPrompter{
		client: params.Client,
	}
}

func (p openAiPrompter) Prompt(msg string) (string, error) {
	res, err := p.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: msg,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}

	return res.Choices[0].Message.Content, nil

}

var _ Prompter = openAiPrompter{}
