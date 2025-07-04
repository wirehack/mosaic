package types

type TestModuleProxy interface {
	Sum(a, b int) int
}

type D map[string]any

type SendRecipient interface {
	To() string
	ToName() string
}

type MailModuleProxy interface {
	Send(recipient SendRecipient, subject, template string, params *D) error
}
