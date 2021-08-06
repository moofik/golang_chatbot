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
				"state %s referenced in transition %s does not exist",
				from,
				transition.Name,
			)
		}
	}

	for _, to := range transition.To {
		if _, ok := d.Places[to]; !ok {
			return fmt.Errorf(
				"state %s referenced in transition %s does not exist",
				to,
				transition.Name,
			)
		}
	}

	d.Transitions = append(d.Transitions, &transition)
	return nil
}
