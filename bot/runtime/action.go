package runtime

import (
	"bot-daedalus/bot/command"
	"bot-daedalus/models"
	"bytes"
	"html/template"
)

type Action interface {
	Run(p ChatProvider, t *models.Token, s *State, c command.Command) error
}

type SendTextMessage struct {
	Text string
}

func (a *SendTextMessage) Run(
	p ChatProvider,
	t *models.Token,
	s *State,
	c command.Command,
) error {
	tmpl, err := template.New("test").Parse(a.Text)
	if err != nil {
		return err
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, t); err != nil {
		return err
	}

	result := tpl.String()
	err = p.SendTextMessage(result)
	if err != nil {
		return err
	}

	return nil
}
