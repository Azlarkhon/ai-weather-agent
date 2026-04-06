package domain

type Tool interface {
	Name() string
	Description() string
	Parameters() map[string]any
	Execute(input map[string]any) (string, error)
}
