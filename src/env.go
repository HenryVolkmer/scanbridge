package main

type Environment struct {
	DefaultRecipient string `json:"default_recipient"`
}

func NewEnvironment(c *Config) *Environment {
	env := &Environment{}
	if len(c.Smtp.Recipient) > 1 {
		env.DefaultRecipient = c.Smtp.Recipient
	}
	return env
}