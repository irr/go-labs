package engines

// Engine ...
type Engine interface {
	GetMessage() (string, error)
	Process(msg *string) error
	Validate(msg *string) bool
}
