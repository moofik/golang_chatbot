package petrinet

import "testing"

func TestGetMarkingForMultistate(t *testing.T) {
	subject := struct {
		States map[string]bool
	}{map[string]bool{"place1": true, "place2": true}}

	storage := MarkingStorage{
		markingField: "States",
		singleState:  false,
	}

	marking, err := storage.GetMarking(&subject)

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if len(marking.Places) != 2 {
		t.Errorf("for current multi state marking places must not be greater than 1")
	}

	if !marking.Has("place1") {
		t.Errorf("current multi state marking place should have place with Name 'place1'")
	}

	if !marking.Has("place2") {
		t.Errorf("current multi state marking place should have place with Name 'place2'")
	}
}

func TestGetMarkingForSinglestate(t *testing.T) {
	subject := struct {
		State string
	}{"place"}

	storage := MarkingStorage{
		markingField: "State",
		singleState:  true,
	}

	marking, err := storage.GetMarking(&subject)

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if len(marking.Places) != 1 {
		t.Errorf("for current single state marking places must not be greater than 1")
	}

	if !marking.Has("place") {
		t.Errorf("current single state marking place should have place with Name 'place'")
	}
}

func TestSetMarkingForMultistate(t *testing.T) {
	subject := struct {
		States map[string]bool
	}{map[string]bool{}}

	storage := MarkingStorage{
		markingField: "States",
		singleState:  false,
	}

	marking := &Marking{map[string]bool{"place1": true, "place2": true}}

	err := storage.SetMarking(&subject, marking)

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if len(subject.States) != 2 {
		t.Errorf("for current multi state marking places must not be greater than 1")
	}

	placeFlag, ok := subject.States["place1"]
	if !ok || placeFlag != true {
		t.Errorf("current subject should have place 'place1' with flag = true")
	}

	placeFlag, ok = subject.States["place2"]
	if !ok || placeFlag != true {
		t.Errorf("current subject should have place 'place2' with flag = true")
	}

	if len(subject.States) != 2 {
		t.Errorf("current subject should have only two places")
	}
}

func TestSetMarkingForSinglestate(t *testing.T) {
	subject := struct {
		State string
	}{}

	storage := MarkingStorage{
		markingField: "State",
		singleState:  true,
	}

	marking := &Marking{map[string]bool{"place": true}}

	err := storage.SetMarking(&subject, marking)

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if subject.State != "place" {
		t.Errorf("current subject should have state = \"place\"")
	}
}
