package cli

import "fmt"

type stateMapper struct{}

func NewStateMapper() *stateMapper {
	return &stateMapper{}
}

func (s *stateMapper) Get(key string) (float64, error) {
	fmt.Println("Get", key)
	return 0, nil
}

func (s *stateMapper) Set(key string, value float64) error {
	fmt.Println("Set", key, value)
	return nil
}
