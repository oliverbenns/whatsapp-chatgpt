package prompt

type Prompter interface {
	Prompt(msg string) (string, error)
}
