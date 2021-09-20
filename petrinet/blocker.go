package petrinet

const CODE_NOT_ENABLED = "not-enabled"
const CODE_UNKNOWN = "unknown"

type Blocker struct {
	message string
	code    string
}

func createNotEnabledBlocker() *Blocker {
	return &Blocker{message: "Transition is prohibited by marking", code: CODE_NOT_ENABLED}
}

func createUnknownBlocker() *Blocker {
	return &Blocker{message: "Transition is prohibited by unknown reason", code: CODE_UNKNOWN}
}
