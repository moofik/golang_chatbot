package petrinet

import "testing"

type Subject struct {
	States map[string]bool
}

func (s *Subject) GetMarkingFieldName() string {
	return "States"
}

func createComplexWorkflowDefinition() *Definition {
	places := map[string]string{}
	for i := 'a'; i <= 'g'; i++ {
		places[string(i)] = string(i)
	}

	transitions := []*Transition{
		&Transition{
			"t1",
			[]string{"a"},
			[]string{"b", "c"},
		},
		&Transition{
			"t2",
			[]string{"b", "c"},
			[]string{"d"},
		},
		&Transition{
			"t3",
			[]string{"d"},
			[]string{"e"},
		},
		&Transition{
			"t4",
			[]string{"d"},
			[]string{"f"},
		},
		&Transition{
			"t5",
			[]string{"e"},
			[]string{"g"},
		},
		&Transition{
			"t6",
			[]string{"f"},
			[]string{"g"},
		},
	}

	// The graph looks like:
	// +---+     +----+     +---+     +----+     +----+     +----+     +----+     +----+     +---+
	// | a | --> | t1 | --> | c | --> | t2 | --> | d  | --> | t4 | --> | f  | --> | t6 | --> | g |
	// +---+     +----+     +---+     +----+     +----+     +----+     +----+     +----+     +---+
	//             |                    ^          |                                           ^
	//             |                    |          |                                           |
	//             v                    |          v                                           |
	//           +----+                 |        +----+     +----+     +----+                  |
	//           | b  | ----------------+        | t3 | --> | e  | --> | t5 | -----------------+
	//           +----+                          +----+     +----+     +----+

	d, _ := CreateDefinition(transitions, nil, places)
	return d
}

func createStateMachineDefinition() *Definition {
	places := map[string]string{}
	for i := 'a'; i <= 'd'; i++ {
		places[string(i)] = string(i)
	}

	transitions := []*Transition{
		&Transition{
			"t1",
			[]string{"a"},
			[]string{"b"},
		},
		&Transition{
			"t1",
			[]string{"d"},
			[]string{"b"},
		},
		&Transition{
			"t2",
			[]string{"b"},
			[]string{"c"},
		},
		&Transition{
			"t3",
			[]string{"b"},
			[]string{"d"},
		},
	}

	// The graph looks like:
	//                     t1
	//               +------------------+
	//               v                  |
	// +---+  t1   +-----+  t2   +---+  |
	// | a | ----> |  b  | ----> | c |  |
	// +---+       +-----+       +---+  |
	//               |                  |
	//               | t3               |
	//               v                  |
	//             +-----+              |
	//             |  d  | -------------+
	//             +-----+

	d, _ := CreateDefinition(transitions, nil, places)
	return d
}

func testGetMarkingWithEmptyInitialMarking(t *testing.T) {
	d := createComplexWorkflowDefinition()
	subject := Subject{}
	workflow := DefaultWorkflow{
		Definition: d,
		MarkingStorage: &MarkingStorage{
			singleState:  false,
			markingField: subject.GetMarkingFieldName(),
		},
		Name: "test",
	}
	marking, _ := workflow.GetMarking(subject)

	if !marking.Has("a") {
		t.Errorf("expected marking is in initial place %s", "a")
	}

	if len(marking.Places) != 1 {
		t.Errorf("expected marking is in only 1 place, got %d", len(marking.Places))
	}
}
