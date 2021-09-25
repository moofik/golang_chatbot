package petrinet

import (
	"errors"
	"fmt"
	"reflect"
)

type Marking struct {
	Places map[string]bool
}

func (m *Marking) Mark(place string) {
	if m.Places == nil {
		m.Places = map[string]bool{}
	}

	m.Places[place] = true
}

func (m *Marking) Unmark(place string) error {
	_, ok := m.Places[place]

	if ok {
		delete(m.Places, place)
	} else {
		return fmt.Errorf("place not found in mark")
	}

	return nil
}

func (m *Marking) Has(place string) bool {
	_, ok := m.Places[place]
	return ok
}

type MarkingStorage struct {
	singleState  bool
	markingField string
}

// GetMarking get marking from a subject
// subject must be a pointer to a struct
func (s *MarkingStorage) GetMarking(subject interface{}) (*Marking, error) {
	rv := reflect.ValueOf(subject)

	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return nil, errors.New("v must be pointer to struct")
	}

	rv = rv.Elem()
	fv := rv.FieldByName(s.markingField)

	if !fv.IsValid() {
		return nil, fmt.Errorf("not a field Name: %s", s.markingField)
	}

	if s.singleState {
		if fv.Kind() != reflect.String {
			return nil, fmt.Errorf("%s is not a string field", s.markingField)
		}

		return &Marking{map[string]bool{fv.String(): true}}, nil
	}

	places := fv.Interface().(map[string]bool)

	return &Marking{places}, nil
}

func (s *MarkingStorage) SetMarking(subject interface{}, m *Marking) error {
	rv := reflect.ValueOf(subject)

	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return errors.New("v must be pointer to struct")
	}

	rv = rv.Elem()
	fv := rv.FieldByName(s.markingField)

	if !fv.IsValid() {
		return fmt.Errorf("not a field Name: %s", s.markingField)
	}

	if !fv.CanSet() {
		return fmt.Errorf("%s is not a settable field", s.markingField)
	}

	if s.singleState {
		if fv.Kind() != reflect.String {
			return fmt.Errorf("%s is not a string field", s.markingField)
		}

		for name, _ := range m.Places {
			fv.SetString(name)
			return nil
		}
	}

	if fv.IsNil() {
		fv.Set(reflect.MakeMap(fv.Type()))
	}

	for name, flag := range m.Places {
		fv.SetMapIndex(reflect.ValueOf(name), reflect.ValueOf(flag))
	}

	return nil
}
