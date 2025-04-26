package judge

import (
	"context"
	"math/rand"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type Problem struct {
	ID int64
	State RunResults
	source []byte
	input []byte
	output []byte
	timeLimit int64
	memoryLimit int64
}

const maxconc = 10

var problems map[int64]*Problem = make(map[int64]*Problem, 0)
var mu sync.Mutex // for protecting the problem list
var sem = semaphore.NewWeighted(maxconc)

var sourcedir = "sources"
var builddir = "compiled"

func makedirs() {
	os.MkdirAll(sourcedir, 0755)
	os.MkdirAll(builddir, 0755)
}

func writeSource(data []byte, id int64) (string, error) {
	name := strconv.FormatInt(id, 10) + ".go"
	pt := path.Join(sourcedir, name)
	return pt, os.WriteFile(pt, data, 0755)
}

func makeSubmissionFile(id int64) (string, error) {
	name := strconv.FormatInt(id, 10) + ".out"
	pt := path.Join(builddir, name)
	return pt, os.WriteFile(pt, []byte{}, 0755)
}

func markProblemAs(id int64, res RunResults) {
	mu.Lock()
	problems[id].State = res
	mu.Unlock()
}

func prepareAndCompile(ctx context.Context, source []byte, id int64) (string, error) {
	sdir, err := writeSource(source, id)
	if err != nil {
		return "", err
	}
	cdir, err := makeSubmissionFile(id)
	if err != nil {
		return "", err
	}
	err = compile(ctx, sdir, cdir)
	if err != nil {
		return "", err
	}
	return cdir, nil
}

func runProblem(source, input, output []byte, timel, mem int64, id int64) {
	ctx := context.Background()

	sem.Acquire(ctx, 1) // making sure we dont die under load
	defer sem.Release(1)

	cdir, err := prepareAndCompile(ctx, source, id)
	if err != nil {
		markProblemAs(id, CompileError)
		return
	}

	res, err := runCode(ctx, cdir, time.Duration(timel),mem, input, output)
	// even if it errors we mark as judge error so its ok
	markProblemAs(id, res)
}

func addProblem(source, input, output []byte, time, mem int64) (int64) {
	id := rand.Int63() // probability of collision is 1:2^63 which is good enough i guess
	problem := Problem{
		ID: id,
		State: Pending,
		source: source,
		input: input,
		output: output,
		timeLimit: time,
		memoryLimit: mem,
	}
	mu.Lock()
	problems[id] = &problem
	mu.Unlock()
	go runProblem(source, input, output, time, mem, id)
	return id
}

func getState(id int64) RunResults {
	mu.Lock()
	defer mu.Unlock()
	prob, ok := problems[id]
	if !ok {
		return JudgeError // just dont want to deal with errors
	}
	return prob.State
}

