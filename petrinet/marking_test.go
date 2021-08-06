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
		t.Errorf("current multi state marking place should have place with name 'place1'")
	}

	if !marking.Has("place2") {
		t.Errorf("current multi state marking place should have place with name 'place2'")
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
		t.Errorf("current single state marking place should have place with name 'place'")
	}
}
