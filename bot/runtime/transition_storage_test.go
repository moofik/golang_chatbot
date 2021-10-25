package runtime

import (
	"bot-daedalus/bot/command"
	"bot-daedalus/petrinet"
	"testing"
)

func TestNext(t *testing.T) {
	ts := TransitionStorage{
		index:         0,
		transitionMap: nil,
	}

	ts.Add(&command.UserInputCommand{
		Text: "1",
		Metadata: &command.Metadata{
			Cmd:        "1",
			Place:      "1",
			Uniqueness: "1",
		},
	}, &petrinet.Transition{
		Name: "1",
		From: []string{"1"},
		To:   []string{"2"},
	})

	ts.Add(&command.UserInputCommand{
		Text: "2",
		Metadata: &command.Metadata{
			Cmd:        "2",
			Place:      "2",
			Uniqueness: "2",
		},
	}, &petrinet.Transition{
		Name: "2",
		From: []string{"2"},
		To:   []string{"3"},
	})

	ts.Add(&command.UserInputCommand{
		Text: "3",
		Metadata: &command.Metadata{
			Cmd:        "3",
			Place:      "3",
			Uniqueness: "3",
		},
	}, &petrinet.Transition{
		Name: "3",
		From: []string{"3"},
		To:   []string{"2"},
	})

	tr, _ := ts.Next()
	if tr.Name != "3" {
		t.Errorf("t name expected to be %d, got %s", 3, tr.Name)
	}

	tr, _ = ts.Next()
	if tr.Name != "1" {
		t.Errorf("t name expected to be %d, got %s", 1, tr.Name)
	}

	tr, _ = ts.Next()
	if tr.Name != "2" {
		t.Errorf("t name expected to be %d, got %s", 2, tr.Name)
	}
}
