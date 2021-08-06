package petrinet

type Transition struct {
	Name string
	From []string
	To   []string
}
