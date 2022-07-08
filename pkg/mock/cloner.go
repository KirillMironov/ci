package mock

import "io"

type Cloner struct{}

func (Cloner) CloneRepository(string) (io.ReadCloser, error) {
	return nil, nil
}
