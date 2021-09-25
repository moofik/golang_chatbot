package petrinet

import "testing"

type Subject struct {
	States map[string]bool
}

func (s *Subject) GetMarkingFieldName() string {
	return "States"
}

func getDefaultMarkingStorage(subject Subject) *MarkingStorage {
	return &MarkingStorage{
		markingField: subject.GetMarkingFieldName(),
		singleState:  false,
	}
}

func makeComplexWorkflowDefinition() *Definition {
	places := map[string]string{}
	for i := 'a'; i <= 'g'; i++ {
		places[string(i)] = string(i)
	}

	transitions := []*Transition{
		{
			"t1",
			[]string{"a"},
			[]string{"b", "c"},
		},
		{
			"t2",
			[]string{"b", "c"},
			[]string{"d"},
		},
		{
			"t3",
			[]string{"d"},
			[]string{"e"},
		},
		{
			"t4",
			[]string{"d"},
			[]string{"f"},
		},
		{
			"t5",
			[]string{"e"},
			[]string{"g"},
		},
		{
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

func makeWorkflowWithSameTransitionName() *Definition {
	places := map[string]string{}
	for i := 'a'; i <= 'g'; i++ {
		places[string(i)] = string(i)
	}

	transitions := []*Transition{
		{
			"a_to_bc",
			[]string{"a"},
			[]string{"b", "c"},
		},
		{
			"b_to_c",
			[]string{"b"},
			[]string{"c"},
		},
		{
			"to_a",
			[]string{"b"},
			[]string{"a"},
		},
		{
			"to_a",
			[]string{"c"},
			[]string{"a"},
		},
	}

	// The graph looks like:
	//   +------------------------------------------------------------+
	//   |                                                            |
	//   |                                                            |
	//   |         +----------------------------------------+         |
	//   v         |                                        v         |
	// +---+     +---------+     +---+     +--------+     +---+     +------+
	// | a | --> | a_to_bc | --> | b | --> | b_to_c | --> | c | --> | to_a | -+
	// +---+     +---------+     +---+     +--------+     +---+     +------+  |
	//   ^                         |                                  ^       |
	//   |                         +----------------------------------+       |
	//   |                                                                    |
	//   |                                                                    |
	//   +--------------------------------------------------------------------+

	d, _ := CreateDefinition(transitions, nil, places)
	return d
}

func makeStateMachineDefinition() *Definition {
	places := map[string]string{}
	for i := 'a'; i <= 'd'; i++ {
		places[string(i)] = string(i)
	}

	transitions := []*Transition{
		{
			"t1",
			[]string{"a"},
			[]string{"b"},
		},
		{
			"t1",
			[]string{"d"},
			[]string{"b"},
		},
		{
			"t2",
			[]string{"b"},
			[]string{"c"},
		},
		{
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

func makeSimpleWorkflowDefinition() *Definition {
	places := map[string]string{}
	for i := 'a'; i <= 'c'; i++ {
		places[string(i)] = string(i)
	}

	transitions := []*Transition{
		{
			"t1",
			[]string{"a"},
			[]string{"b"},
		},
		{
			"t2",
			[]string{"b"},
			[]string{"c"},
		},
	}

	d, _ := CreateDefinition(transitions, nil, places)
	return d
}

func TestGetMarkingEmptyInitialMarking(t *testing.T) {
	d := makeComplexWorkflowDefinition()
	subject := Subject{}
	workflow := DefaultWorkflow{
		Definition: d,
		MarkingStorage: &MarkingStorage{
			singleState:  false,
			markingField: subject.GetMarkingFieldName(),
		},
		Name: "test",
	}
	marking, err := workflow.GetMarking(&subject)

	if err != nil {
		t.Errorf(err.Error())
	}

	if !marking.Has("a") {
		t.Errorf("expected marking is in initial place %s, got %v", "'a'", marking)
	}

	if len(marking.Places) != 1 {
		t.Errorf("expected marking is in only 1 place, got %d", len(marking.Places))
	}
}

func TestGetMarkingImpossiblePlace(t *testing.T) {
	subject := Subject{map[string]bool{"imp": true}}
	storage := getDefaultMarkingStorage(subject)
	workflow := DefaultWorkflow{
		Definition: &Definition{
			Transitions: []*Transition{},
			Places:      map[string]string{},
		},
		MarkingStorage: storage,
		Name:           "test",
	}

	_, err := workflow.GetMarking(&subject)
	expectedMessage := "it seems you forgot to add place 'imp' to the workflow 'test'"

	if err == nil || err.Error() != expectedMessage {
		t.Errorf("error expected (forgot to add place 'imp' to the workflow 'test')")
	}
}

func TestGetMarkingEmptyDefinition(t *testing.T) {
	subject := Subject{}
	storage := getDefaultMarkingStorage(subject)
	workflow := DefaultWorkflow{
		Definition: &Definition{
			Transitions: []*Transition{},
			Places:      map[string]string{},
		},
		MarkingStorage: storage,
		Name:           "test",
	}
	_, err := workflow.GetMarking(&subject)
	expectedMessage := "the Marking is empty and there is no initial place for workflow test"

	if err == nil || err.Error() != expectedMessage {
		t.Errorf("error expected (the Marking is empty and there is no initial place for workflow test)")
	}
}

func TestGetMarkingWithExistentMarking(t *testing.T) {
	d := makeComplexWorkflowDefinition()
	subject := Subject{map[string]bool{"b": true, "c": true}}
	storage := getDefaultMarkingStorage(subject)
	workflow := DefaultWorkflow{
		Definition:     d,
		MarkingStorage: storage,
		Name:           "test",
	}
	marking, err := workflow.GetMarking(&subject)

	if err != nil {
		t.Errorf(err.Error())
	}

	if !marking.Has("b") {
		t.Errorf("expected marking has place b")
	}

	if !marking.Has("c") {
		t.Errorf("expected marking has place c")
	}
}

func TestCanFire(t *testing.T) {
	d := makeComplexWorkflowDefinition()
	subject := Subject{}
	storage := getDefaultMarkingStorage(subject)
	workflow := DefaultWorkflow{
		Definition:     d,
		MarkingStorage: storage,
		Name:           "test",
	}

	if can, _ := workflow.CanFire(&subject, "t1"); !can {
		t.Errorf("expected workflow can fire transition t1 on subject")
	}

	if can, _ := workflow.CanFire(&subject, "t2"); can {
		t.Errorf("expected workflow can not fire transition t2 on subject")
	}

	subject.States = map[string]bool{"b": true}

	if can, _ := workflow.CanFire(&subject, "t1"); can {
		t.Errorf("expected workflow can not fire transition t1 on subject")
	}

	if can, _ := workflow.CanFire(&subject, "t1"); can {
		t.Errorf("expected workflow can not fire transition t1 on subject")
	}

	subject.States = map[string]bool{"b": true, "c": true}

	if can, _ := workflow.CanFire(&subject, "t1"); can {
		t.Errorf("expected workflow can not fire transition t1 on subject")
	}

	if can, _ := workflow.CanFire(&subject, "t2"); !can {
		t.Errorf("expected workflow can fire transition t1 on subject")
	}

	subject.States = map[string]bool{"f": true}

	if can, _ := workflow.CanFire(&subject, "t5"); can {
		t.Errorf("expected workflow can not fire transition t1 on subject")
	}

	if can, _ := workflow.CanFire(&subject, "t6"); !can {
		t.Errorf("expected workflow can fire transition t1 on subject")
	}
}

func TestCanFireWithNonExistentTransition(t *testing.T) {
	d := makeComplexWorkflowDefinition()
	subject := Subject{}
	storage := getDefaultMarkingStorage(subject)
	workflow := DefaultWorkflow{
		Definition:     d,
		MarkingStorage: storage,
		Name:           "test",
	}

	if can, _ := workflow.CanFire(&subject, "foobar"); can {
		t.Errorf("expected workflow can not fire non existent transition foobar on subject")
	}
}

func TestCanFireWithSameTransitionName(t *testing.T) {
	d := makeWorkflowWithSameTransitionName()
	subject := Subject{}
	storage := getDefaultMarkingStorage(subject)
	workflow := DefaultWorkflow{
		Definition:     d,
		MarkingStorage: storage,
		Name:           "",
	}

	if can, _ := workflow.CanFire(&subject, "a_to_bc"); !can {
		t.Errorf("expected workflow can fire transition a_to_bc on subject")
	}

	if can, _ := workflow.CanFire(&subject, "b_to_c"); can {
		t.Errorf("expected workflow can not fire transition b_to_c on subject")
	}

	if can, _ := workflow.CanFire(&subject, "to_a"); can {
		t.Errorf("expected workflow can not fire transition to_a on subject")
	}

	subject.States = map[string]bool{"b": true}

	if can, _ := workflow.CanFire(&subject, "a_to_bc"); can {
		t.Errorf("expected workflow can not fire transition a_to_bc on subject")
	}

	if can, _ := workflow.CanFire(&subject, "b_to_c"); !can {
		t.Errorf("expected workflow can fire transition b_to_c on subject")
	}

	if can, _ := workflow.CanFire(&subject, "to_a"); !can {
		t.Errorf("expected workflow can fire transition to_a on subject")
	}
}

func TestBuildTransitionBlockerListReturnsUndefTransition(t *testing.T) {
	d := makeSimpleWorkflowDefinition()
	subject := Subject{}
	storage := getDefaultMarkingStorage(subject)
	workflow := DefaultWorkflow{
		Definition:     d,
		MarkingStorage: storage,
		Name:           "test",
	}

	_, err := workflow.BuildTransitionBlockerList(&subject, "nf")
	expectedMessage := "transition name nf is not defined for workflow test"
	if err == nil || err.Error() != expectedMessage {
		t.Errorf("expected error: " + expectedMessage)
	}
}

func TestWorkflowBuildTransitionBlockerList(t *testing.T) {
	d := makeComplexWorkflowDefinition()
	subject := Subject{}
	storage := getDefaultMarkingStorage(subject)
	workflow := DefaultWorkflow{
		Definition:     d,
		MarkingStorage: storage,
		Name:           "test",
	}

	if l, _ := workflow.BuildTransitionBlockerList(&subject, "t1"); !l.empty() {
		t.Errorf("expected blocker list is empty")
	}

	if l, _ := workflow.BuildTransitionBlockerList(&subject, "t2"); l.empty() {
		t.Errorf("expected blocker list is not empty")
	}

	subject.States = map[string]bool{"b": true}

	if l, _ := workflow.BuildTransitionBlockerList(&subject, "t1"); l.empty() {
		t.Errorf("expected blocker list not empty")
	}

	if l, _ := workflow.BuildTransitionBlockerList(&subject, "t2"); l.empty() {
		t.Errorf("expected blocker list is not empty")
	}

	subject.States = map[string]bool{"b": true, "c": true}

	if l, _ := workflow.BuildTransitionBlockerList(&subject, "t1"); l.empty() {
		t.Errorf("expected blocker list not empty")
	}

	if l, _ := workflow.BuildTransitionBlockerList(&subject, "t2"); !l.empty() {
		t.Errorf("expected blocker list is empty")
	}

	subject.States = map[string]bool{"f": true}

	if l, _ := workflow.BuildTransitionBlockerList(&subject, "t5"); l.empty() {
		t.Errorf("expected blocker list not empty")
	}

	if l, _ := workflow.BuildTransitionBlockerList(&subject, "t6"); !l.empty() {
		t.Errorf("expected blocker list is empty")
	}
}

func TestBlockerListReasonsReturned(t *testing.T) {
	d := makeComplexWorkflowDefinition()
	subject := Subject{}
	storage := getDefaultMarkingStorage(subject)
	workflow := DefaultWorkflow{
		Definition:     d,
		MarkingStorage: storage,
		Name:           "test",
	}
	blockers, err := workflow.BuildTransitionBlockerList(&subject, "t2")

	if err != nil {
		t.Errorf(err.Error())
	}

	if blockers.count() != 1 {
		t.Errorf("expected blockers count of 1")
	}

	expectedMessage := "Transition is prohibited by marking"
	if blocker, _ := blockers.next(); blocker.message != expectedMessage {
		t.Errorf("expected error message: " + expectedMessage)
	}

	expectedCode := "not-enabled"
	if blocker := blockers.current(); blocker.code != expectedCode {
		t.Errorf("expected error code: " + expectedCode)
	}
}

func TestFireWithNonExistentTransition(t *testing.T) {
	d := makeComplexWorkflowDefinition()
	subject := Subject{}
	storage := getDefaultMarkingStorage(subject)
	workflow := DefaultWorkflow{
		Definition:     d,
		MarkingStorage: storage,
		Name:           "test",
	}

	_, err := workflow.Fire(&subject, "nf")
	expectedMessage := "transition nf is not defined for workflow test"
	if err == nil || err.Error() != expectedMessage {
		t.Errorf("error expected: " + expectedMessage)
	}
}

func TestFireWithNotEnabledTransition(t *testing.T) {
	d := makeComplexWorkflowDefinition()
	subject := Subject{}
	storage := getDefaultMarkingStorage(subject)
	workflow := DefaultWorkflow{
		Definition:     d,
		MarkingStorage: storage,
		Name:           "test",
	}

	_, err := workflow.Fire(&subject, "t2")

	expectedMessage := "transition t2 is not enabled for workflow test"
	if err == nil || err.Error() != expectedMessage {
		t.Errorf("expected error: " + expectedMessage)
	}

	b, _ := err.GetBlockerList().next()
	expectedMessage = "Transition is prohibited by marking"
	if b.message != expectedMessage {
		t.Errorf("expected error message: " + expectedMessage)
	}

	if b.code != "not-enabled" {
		t.Errorf("expected error message: not-enabled")
	}

	if err.GetTransitionName() != "t2" {
		t.Errorf("expected transition name in error: " + err.GetTransitionName())
	}
}

func TestFire(t *testing.T) {
	d := makeComplexWorkflowDefinition()
	subject := Subject{}
	storage := getDefaultMarkingStorage(subject)
	workflow := DefaultWorkflow{
		Definition:     d,
		MarkingStorage: storage,
		Name:           "test",
	}

	marking, err := workflow.Fire(&subject, "t1")

	if err != nil {
		t.Errorf(err.Error())
	}

	if marking.Has("a") {
		t.Errorf("expected marking has no place 'a'")
	}

	if !marking.Has("b") {
		t.Errorf("expected marking has place 'b'")
	}

	if !marking.Has("c") {
		t.Errorf("expected marking has place 'c'")
	}
}

func TestFireWithSameTransitionName(t *testing.T) {
	d := makeWorkflowWithSameTransitionName()
	subject := Subject{}
	storage := getDefaultMarkingStorage(subject)
	workflow := DefaultWorkflow{
		Definition:     d,
		MarkingStorage: storage,
		Name:           "test",
	}
	marking, err := workflow.Fire(&subject, "a_to_bc")

	if err != nil {
		t.Errorf(err.Error())
	}

	if ok := marking.Has("a"); ok {
		t.Errorf("expected marking has no place 'a'")
	}

	if ok := marking.Has("b"); !ok {
		t.Errorf("expected marking has place 'b'")
	}

	if ok := marking.Has("c"); !ok {
		t.Errorf("expected marking has place 'c'")
	}

	marking, err = workflow.Fire(&subject, "to_a")

	if err != nil {
		t.Errorf(err.Error())
	}

	if ok := marking.Has("a"); !ok {
		t.Errorf("expected marking has place 'a'")
	}

	if ok := marking.Has("b"); ok {
		t.Errorf("expected marking has no place 'b'")
	}

	if ok := marking.Has("c"); ok {
		t.Errorf("expected marking has no place 'c'")
	}

	marking, err = workflow.Fire(&subject, "a_to_bc")
	marking, err = workflow.Fire(&subject, "b_to_c")

	if err != nil {
		t.Errorf(err.Error())
	}

	if ok := marking.Has("a"); ok {
		t.Errorf("expected marking has no place 'a'")
	}

	if ok := marking.Has("b"); ok {
		t.Errorf("expected marking has no place 'b'")
	}

	if ok := marking.Has("c"); !ok {
		t.Errorf("expected marking has place 'c'")
	}

	marking, err = workflow.Fire(&subject, "to_a")

	if err != nil {
		t.Errorf(err.Error())
	}

	if ok := marking.Has("a"); !ok {
		t.Errorf("expected marking has place 'a'")
	}

	if ok := marking.Has("b"); ok {
		t.Errorf("expected marking has no place 'b'")
	}

	if ok := marking.Has("c"); ok {
		t.Errorf("expected marking has no place 'c'")
	}
}

func TestFireWithSameTransitionName2(t *testing.T) {
	places := map[string]string{}
	for i := 'a'; i <= 'd'; i++ {
		places[string(i)] = string(i)
	}
	d, _ := CreateDefinition([]*Transition{
		{
			"t",
			[]string{"a"},
			[]string{"c"},
		},
		{
			"t",
			[]string{"b"},
			[]string{"d"},
		},
	}, nil, places)
	subject := Subject{map[string]bool{"a": true, "b": true}}
	storage := getDefaultMarkingStorage(subject)
	workflow := DefaultWorkflow{
		Definition:     d,
		MarkingStorage: storage,
		Name:           "test",
	}
	marking, err := workflow.Fire(&subject, "t")

	if err != nil {
		t.Errorf(err.Error())
	}

	if ok := marking.Has("a"); ok {
		t.Errorf("expected marking has no place 'a'")
	}

	if ok := marking.Has("b"); ok {
		t.Errorf("expected marking has no place 'b'")
	}

	if ok := marking.Has("c"); !ok {
		t.Errorf("expected marking has no place 'c'")
	}

	if ok := marking.Has("d"); !ok {
		t.Errorf("expected marking has no place 'd'")
	}
}

func TestFireWithSameTransitionName3(t *testing.T) {
	places := map[string]string{}
	for i := 'a'; i <= 'd'; i++ {
		places[string(i)] = string(i)
	}
	d, _ := CreateDefinition([]*Transition{
		{
			"t",
			[]string{"a"},
			[]string{"c"},
		},
		{
			"t",
			[]string{"b"},
			[]string{"d"},
		},
	}, nil, places)
	subject := Subject{map[string]bool{"a": true}}
	storage := getDefaultMarkingStorage(subject)
	workflow := DefaultWorkflow{
		Definition:     d,
		MarkingStorage: storage,
		Name:           "test",
	}
	marking, err := workflow.Fire(&subject, "t")

	if err != nil {
		t.Errorf(err.Error())
	}

	if ok := marking.Has("a"); ok {
		t.Errorf("expected marking has no place 'a'")
	}

	if ok := marking.Has("b"); ok {
		t.Errorf("expected marking has no place 'b'")
	}

	if ok := marking.Has("c"); !ok {
		t.Errorf("expected marking has place 'c'")
	}

	if ok := marking.Has("d"); ok {
		t.Errorf("expected marking has no place 'd'")
	}
}

func TestFireWithSameTransitionName4(t *testing.T) {
	places := map[string]string{}
	for i := 'a'; i <= 'd'; i++ {
		places[string(i)] = string(i)
	}
	d, _ := CreateDefinition([]*Transition{
		{
			"t",
			[]string{"a"},
			[]string{"b"},
		},
		{
			"t",
			[]string{"b"},
			[]string{"c"},
		},
		{
			"t",
			[]string{"c"},
			[]string{"d"},
		},
	}, nil, places)

	subject := Subject{map[string]bool{"a": true}}
	storage := getDefaultMarkingStorage(subject)
	workflow := DefaultWorkflow{
		Definition:     d,
		MarkingStorage: storage,
		Name:           "test",
	}
	marking, err := workflow.Fire(&subject, "t")

	if err != nil {
		t.Errorf(err.Error())
	}

	if ok := marking.Has("b"); !ok {
		t.Errorf("expected marking has place 'b'")
	}

	if ok := marking.Has("d"); ok {
		t.Errorf("expected marking has no place 'd'")
	}
}
