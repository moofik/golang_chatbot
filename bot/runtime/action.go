package runtime

import (
	"bot-daedalus/bot/command"
	"bytes"
	"html/template"
)

type Action interface {
	Run(p ChatProvider, t TokenProxy, s *State, c command.Command) error
}

type SendTextMessage struct {
	Text string
}

func (a *SendTextMessage) Run(
	p ChatProvider,
	t TokenProxy,
	s *State,
	c command.Command,
) error {
	tmpl, err := template.New("test").Parse(a.Text)
	if err != nil {
		return err
	}

	var tpl bytes.Buffer

	if err := tmpl.Execute(&tpl, t.ToPlainStruct()); err != nil {
		return err
	}

	result := tpl.String()
	err = p.SendTextMessage(result)
	if err != nil {
		return err
	}

	return nil
}
