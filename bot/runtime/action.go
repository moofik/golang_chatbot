package runtime

import (
	"bytes"
	"html/template"
)

type Action interface {
	Run(p ChatProvider, t TokenProxy, s *State, prev *State, c Command) ActionError
	GetName() string
}

type SendTextMessage struct {
	params map[string]string
}

func (a *SendTextMessage) GetName() string {
	return "send_text"
}

func (a *SendTextMessage) Run(
	p ChatProvider,
	t TokenProxy,
	s *State,
	prev *State,
	c Command,
) ActionError {
	tmpl, err := template.New("test").Parse(a.params["text"])
	if err != nil {
		return &GenericActionError{InnerError: err}
	}

	var tpl bytes.Buffer

	data := t.GetExtras()

	if err := tmpl.Execute(&tpl, data); err != nil {
		return &GenericActionError{InnerError: err}
	}

	result := tpl.String()
	err = p.SendTextMessage(result, ProviderContext{
		State:   s,
		Command: c,
		Token:   t,
	})
	if err != nil {
		return &GenericActionError{InnerError: err}
	}

	return nil
}

type RememberInput struct {
	params map[string]string
}

func (a *RememberInput) GetName() string {
	return "remember_input"
}

func (a *RememberInput) Run(
	p ChatProvider,
	t TokenProxy,
	s *State,
	prev *State,
	c Command,
) ActionError {
	extras := t.GetExtras()
	extras[a.params["var"]] = c.GetInput()
	t.SetExtras(extras)
	return nil
}

type RememberCaption struct {
	params map[string]string
}

func (a *RememberCaption) GetName() string {
	return "remember_input"
}

func (a *RememberCaption) Run(
	p ChatProvider,
	t TokenProxy,
	s *State,
	prev *State,
	c Command,
) ActionError {
	newCmd := &ButtonPressedCommand{
		ButtonCommand: c.GetInput(),
		ButtonText:    c.GetCaption(),
		Metadata: &CommandMetadata{
			Cmd:        "button",
			Place:      s.Name,
			Uniqueness: c.GetInput(),
		}}

	cmd, _ := prev.GetCommandByUniqueness(newCmd)

	extras := t.GetExtras()
	extras[a.params["var"]] = cmd.GetCaption()
	t.SetExtras(extras)
	return nil
}

func CreateAction(name string, params map[string]string, actionRegistry func(string, map[string]string) Action) Action {
	if name == "send_text" {
		return &SendTextMessage{params: params}
	}

	if name == "remember_input" {
		return &RememberInput{params: params}
	}

	if name == "remember_caption" {
		return &RememberCaption{params: params}
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
	InnerError error
}

func (m *GenericActionError) Error() string {
	return m.InnerError.Error()
}
