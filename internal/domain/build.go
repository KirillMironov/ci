package domain

import "encoding/json"

type Build struct {
	Commit Commit `json:"commit"`
	Status Status `json:"status"`
	LogId  int    `json:"log_id"`
}

func (b Build) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Commit Commit `json:"commit"`
		Status string `json:"status"`
		LogId  int    `json:"log_id"`
	}{
		Commit: b.Commit,
		Status: b.Status.String(),
		LogId:  b.LogId,
	})
}
