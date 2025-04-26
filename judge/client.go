package judge

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
)

func addField(w *multipart.Writer, name string, content string) error {
	fw, err := w.CreateFormField(name)
	if err != nil {
		return err
	}
	fw.Write([]byte(content))
	return nil
}

func addFile(w *multipart.Writer, name string, filename string) error {
	fw, err := w.CreateFormFile(name, filename)
	if err != nil {
		return err
	}
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	io.Copy(fw, file)
	return nil
}

func submitToJudge(source, input, output string, timelimit, memlimit int64, url string) (int64, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// time limit
	s := strconv.FormatInt(timelimit, 10)
	err := addField(w, timename, s)
	if err != nil {return 0, err}
	// memory limit
	s = strconv.FormatInt(memlimit, 10)
	err = addField(w, memname, s)
	if err != nil {return 0, err}
	// source file
	err = addFile(w, srcname, source)
	if err != nil {return 0, err}
	// input/output files
	err = addFile(w, inputname, input)
	if err != nil {return 0, err}
	err = addFile(w, outputname, output)
	if err != nil {return 0, err}
	// I fucking hate this go error handling shit man

	w.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil { return 0, err }

	// we should set the content type because its important and also contains the boundary
	req.Header.Set("Content-Type", w.FormDataContentType())

	res, err := http.DefaultClient.Do(req)
	if err != nil { return 0, err }
	
	d, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return 0, errors.New(string(d))
	}

	id, err := strconv.ParseInt(string(d), 10, 0)
	return id, err
}

func queryState(id int64, url string) (RunResults, error){
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// time limit
	s := strconv.FormatInt(id, 10)
	err := addField(w, idname, s)
	if err != nil {return JudgeError, err}
	w.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {return JudgeError, err}

	// we should set the content type because its important and also contains the boundary
	req.Header.Set("Content-Type", w.FormDataContentType())

	res, err := http.DefaultClient.Do(req)
	if err != nil {return JudgeError, err}
	
	d, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return JudgeError, errors.New(string(d))
	}

	runres, err := resultFromString(string(d))

	return runres, err
}

func dumpRequestOut(r *http.Request) {
	b, _ := httputil.DumpRequestOut(r, true)
	fmt.Println(string(b))
}

