package mock

type Cloner struct{}

func (Cloner) Clone(string) (string, error) {
	return "", nil
}
