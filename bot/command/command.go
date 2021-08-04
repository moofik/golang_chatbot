package command

type Command interface {

}

type UserInputCommand struct {
	Text string
}

type MockCommand struct {

}