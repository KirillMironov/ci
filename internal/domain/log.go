package domain

type Log struct {
	Data []byte
}

type LogsStorage interface {
	GetByBuildId(buildId string) (Log, error)
}
