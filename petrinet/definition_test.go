package petrinet

import (
	"testing"
)

func TestAddPlaces(t *testing.T) {
	d, _ := CreateDefinition(
		nil,
		nil,
		map[string]string{
			"a": "a",
			"b": "b",
		},
	)

	if len(d.Places) != 2 {
		t.Errorf("unexpected places count")
	}

	if d.Places["a"] != "a" || d.Places["b"] != "b" {
		t.Errorf("expected places is not defined in places of Definition %v, places are: ", d.Places)
	}
}

func TestSetInitialPlaces(t *testing.T) {
	_, err := CreateDefinition(
		nil,
		map[string]string{
			"a": "a",
		},
		map[string]string{
			"a": "a",
			"b": "a",
		},
	)

	if err != nil {
		t.Errorf("unexpected error: initial place exist but marked as non-existent")
	}
}

func TestSetInitialPlaceAndPlaceIsNotDefined(t *testing.T) {
	_, err := CreateDefinition(
		nil,
		map[string]string{
			"x": "x",
		},
		map[string]string{
			"a": "a",
			"b": "a",
		},
	)

	if err == nil {
		t.Errorf("non-existent place error expected")
	}
}

func TestTransitionWithFromPlaceNotDefined(t *testing.T) {
	_, err := CreateDefinition(
		[]*Transition{
			&Transition{
				Name: "test",
				From: []string{"x"},
				To:   []string{"b"},
			}},
		nil,
		map[string]string{
			"a": "a",
			"b": "a",
		},
	)

	if err == nil {
		t.Errorf("non-existent place error expected")
	}
}

func TestTransitionWithToPlaceIsNotDefined(t *testing.T) {
	_, err := CreateDefinition(
		[]*Transition{
			&Transition{
				Name: "test",
				From: []string{"a"},
				To:   []string{"x"},
			}},
		nil,
		map[string]string{
			"a": "a",
			"b": "a",
		},
	)

	if err == nil {
		t.Errorf("non-existent place error expected")
	}
}

func TestDefinitionAddPlace(t *testing.T) {
	d, _ := CreateDefinition(nil, nil, nil)
	d.AddPlace("x")

	if len(d.Places) != 1 {
		t.Errorf("exactly 1 place expected, got %d", len(d.Places))
	}

	if d.Places["x"] != "x" {
		t.Errorf("expected place with Name %s exist, places are: %v", "x", d.Places)
	}
}

func TestDefinitionAddTransition(t *testing.T) {
	d, _ := CreateDefinition(nil, nil, nil)
	d.AddPlace("a")
	d.AddPlace("b")
	tr := Transition{
		"test",
		[]string{"a"},
		[]string{"b"},
	}
	err := d.AddTransition(&tr)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(d.Transitions) != 1 {
		t.Errorf("exactly 1 transition expected, got %d", len(d.Transitions))
	}

	if d.Transitions[0].Name != tr.Name {
		t.Errorf("expected transition with Name %s, got %s", d.Transitions[0].Name, tr.Name)
	}

	if d.Transitions[0].From[0] != tr.From[0] {
		t.Errorf("expected transition with froms %s, got %s", d.Transitions[0].From[0], tr.From[0])
	}

	if d.Transitions[0].To[0] != tr.To[0] {
		t.Errorf("expected transition with Name %s, got %s", d.Transitions[0].To[0], tr.To[0])
	}

}
