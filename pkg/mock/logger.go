package mock

type Logger struct{}

func (Logger) Info(...interface{}) {}

func (Logger) Infof(string, ...interface{}) {}

func (Logger) Error(...interface{}) {}

func (Logger) Errorf(string, ...interface{}) {}
