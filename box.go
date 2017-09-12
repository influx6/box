package box

// CancelContext defines a type which provides Done signal for cancelling operations.
type CancelContext interface {
	Done() <-chan struct{}
}

// Spell defines an interface which expose an exec method.
type Spell interface {
	Exec(CancelContext) error
}
