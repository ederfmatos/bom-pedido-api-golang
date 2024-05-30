package registry

var values map[string]any

func RegisterDependency(name string, value any) {
	if values == nil {
		values = make(map[string]any)
	}
	values[name] = value
}

func GetDependency[T any](name string) *T {
	dependency := values[name]
	if value, ok := dependency.(*T); ok {
		return value
	}
	return nil
}
