package domain

type Log struct {
	Data string `json:"data"`
}

type LogsStorage interface {
	GetByBuildId(buildId string) (Log, error)
}
