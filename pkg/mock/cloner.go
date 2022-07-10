package mock

type Cloner struct{}

func (Cloner) CloneRepository(string) (string, func() error, error) {
	return "", nil, nil
}
