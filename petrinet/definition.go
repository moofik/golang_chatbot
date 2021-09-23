package petrinet

import "fmt"

type Definition struct {
	Transitions   []*Transition
	InitialPlaces map[string]string
	Places        map[string]string
}

func (d *Definition) AddPlace(p string) {
	if len(d.Places) == 0 {
		d.InitialPlaces[p] = p
	}

	d.Places[p] = p
}

func (d *Definition) AddTransition(transition Transition) error {
	for _, from := range transition.From {
		if _, ok := d.Places[from]; !ok {
			return fmt.Errorf(
				"place %s referenced in transition %s does not exist",
				from,
				transition.Name,
			)
		}
	}

	for _, to := range transition.To {
		if _, ok := d.Places[to]; !ok {
			return fmt.Errorf(
				"place %s referenced in transition %s does not exist",
				to,
				transition.Name,
			)
		}
	}

	d.Transitions = append(d.Transitions, &transition)
	return nil
}

func CreateDefinition(t []*Transition, InitialPlaces map[string]string, Places map[string]string) (*Definition, error) {
	//check non-existent initial place
	for name, _ := range InitialPlaces {
		if _, ok := Places[name]; !ok {
			return nil, &NonExistentPlaceError{name}
		}
	}
	//check transition doesnt have a to place that is not defined in definition places
	for _, transition := range t {
		for _, to := range transition.To {
			if _, ok := Places[to]; !ok {
				return nil, &NonExistentPlaceError{to}
			}
		}
	}
	//check transition doesnt have a from place that is not defined in definition places
	for _, transition := range t {
		for _, from := range transition.From {
			if _, ok := Places[from]; !ok {
				return nil, &NonExistentPlaceError{from}
			}
		}
	}

	return &Definition{t, InitialPlaces, Places}, nil
}

type NonExistentPlaceError struct {
	name string
}

func (e *NonExistentPlaceError) Error() string {
	return fmt.Sprintf("place %s does not exist in the workflow", e.name)
}
