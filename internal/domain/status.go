package domain

const (
	Success Status = iota
	Failure
	Skipped
)

type Status uint8

func (s Status) String() string {
	return [...]string{"success", "failure", "skipped"}[s]
}
