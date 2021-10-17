package runtime

import (
	"bot-daedalus/bot/command"
	"bytes"
	"html/template"
)

type Action interface {
	Run(p ChatProvider, t TokenProxy, s *State, c command.Command) ActionError
}

type SendTextMessage struct {
	params map[string]string
}

func (a *SendTextMessage) Run(
	p ChatProvider,
	t TokenProxy,
	s *State,
	c command.Command,
) ActionError {
	tmpl, err := template.New("test").Parse(a.params["text"])
	if err != nil {
		return &GenericActionError{innerError: err}
	}

	var tpl bytes.Buffer

	if err := tmpl.Execute(&tpl, t.ToPlainStruct()); err != nil {
		return &GenericActionError{innerError: err}
	}

	result := tpl.String()
	err = p.SendTextMessage(result, ProviderContext{
		State:   s,
		Command: c,
		Token:   t,
	})
	if err != nil {
		return &GenericActionError{innerError: err}
	}

	return nil
}

func CreateAction(name string, params map[string]string, actionRegistry func(string, map[string]string) Action) Action {
	if name == "send_text" {
		return &SendTextMessage{params: params}
	}

	if actionRegistry != nil {
		return actionRegistry(name, params)
	}

	return nil
}

type ActionError interface {
	error
}

type GenericActionError struct {
	innerError error
}

func (m *GenericActionError) Error() string {
	return m.innerError.Error()
}
