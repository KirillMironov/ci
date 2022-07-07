package mock

type Cloner struct{}

func (Cloner) CloneRepository(string) (string, error) {
	return "", nil
}
