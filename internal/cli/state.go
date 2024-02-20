package cli

type stateMapper struct {
	cache map[string]float64
}

func NewStateMapper() *stateMapper {
	return &stateMapper{
		cache: map[string]float64{},
	}
}

func (s *stateMapper) Get(key string) (float64, error) {
	return s.cache[key], nil
}

func (s *stateMapper) Set(key string, value float64) error {
	s.cache[key] = value
	return nil
}
