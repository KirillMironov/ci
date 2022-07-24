package domain

const (
	Success = iota
	Failure
	Skipped
)

type (
	Build struct {
		Commit Commit
		Status Status
		LogId  int
	}

	Status uint8
)

func (s Status) String() string {
	return [...]string{"success", "failure", "skipped"}[s]
}
