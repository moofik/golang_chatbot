package runtime

type MaintenanceHandler interface {
	Handle(cmd Command, currentState *State, token TokenProxy) bool
}
