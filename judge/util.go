package judge

import "errors"

type RunResults int64

const (
	Accepted RunResults = iota
	Wrong
	CompileError
	MemoryLimit
	TimeLimit
	RuntimeError
	JudgeError
	Pending
)

var (
	list_run_results = [...]string{
		"Accepted",
		"Wrong",
		"CompileError",
		"MemoryLimit",
		"TimeLimit",
		"RuntimeError",
		"JudgeError",
		"Pending",
	}
)

func (r RunResults) String() string {
	return list_run_results[r]
}

func resultFromString(s string) (RunResults, error) {
	for i, x := range list_run_results {
		if s == x {
			return RunResults(i), nil
		}
	}
	return JudgeError, errors.New("bad name")
}

