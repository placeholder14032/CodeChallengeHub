package main

import (
  "fmt"
  "sync"
  "testing"
  "time"
  "github.com/placeHolder143032/CodeChallengeHub/judge"
)

func oneSubmit(i int) {
  id, err := judge.SubmitToJudge("./testi/source", "./testi/input", "./testi/output", 1_000_000_000, 10_000_000, "http://localhost:8081/submit")
  fmt.Println("got id: ", id, err)
  //
  time.Sleep(10 * time.Second)
  //
  res, err := judge.QueryState(id, "http://localhost:8081/query")
  fmt.Printf("result of %d after 10 seconds: %v %v\n", i, res, err)
}

func TestLoad(t *testing.T) {
  // go judge.Server(8081)
  // time.Sleep(1 * time.Second) // waiting for server to start up properly
  var wg sync.WaitGroup
  for i := 0; i < 1000; i++ {
    wg.Add(1)
    go func() {
      oneSubmit(i)
      wg.Done()
    }()
  }
  wg.Wait()
}