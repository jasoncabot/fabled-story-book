package cli

type stateMapper struct {
	cache map[string]any
}

func NewStateMapper() *stateMapper {
	return &stateMapper{
		cache: map[string]any{},
	}
}

func (s *stateMapper) Get(key string) (any, error) {
	return s.cache[key], nil
}

func (s *stateMapper) Set(key string, value any) error {
	s.cache[key] = value
	return nil
}
