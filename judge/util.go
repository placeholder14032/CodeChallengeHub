package judge

import "errors"

type RunResults int64

const (
	Pending RunResults = iota
	Accepted
	CompileError
	Wrong
	MemoryLimit
	TimeLimit
	RuntimeError
	JudgeError
)

var (
	list_run_results = [...]string{
		"Pending", 
		"Accepted", 
		"CompileError", 
		"Wrong", 
		"MemoryLimit", 
		"TimeLimit", 
		"RuntimeError", 
		"JudgeError", 
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

