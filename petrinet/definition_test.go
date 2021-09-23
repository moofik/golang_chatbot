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
		t.Errorf("expected places is not defined in places of definition %v, places are: ", d.Places)
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

func TestAddTransitionAndFromPlaceIsNotDefined(t *testing.T) {
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

func TestAddTransitionAndToPlaceIsNotDefined(t *testing.T) {
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
