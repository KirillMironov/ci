package domain

type Log struct {
	Data []byte
}

type LogsUsecase interface {
	GetByBuildId(buildId string) (Log, error)
}
