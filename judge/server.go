package judge

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
)

const (
	default_port = 8081
	max_upload_size = 1024 * 1024 * 2 // 2 Mib
	srcname = "source"
	inputname = "input"
	outputname = "output"
	timename = "timelimit"
	memname = "memorylimit"
	idname = "ID"
)

func sendError(w http.ResponseWriter, msg string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(msg))
}

func getFile(w http.ResponseWriter, r *http.Request, name string) ([]byte, error) {
	// parse and validate file and post parameters
	file, fileHeader, err := r.FormFile(name)
	if err != nil {
		sendError(w, "INVALID_FILE", http.StatusBadRequest)
		return nil, fmt.Errorf("invalid file")
	}
	defer file.Close()
	// Get and print out file size
	fileSize := fileHeader.Size
	fmt.Printf("File size (bytes): %v\n", fileSize)
	// validate file size
	if fileSize > max_upload_size {
		sendError(w, "FILE_TOO_BIG", http.StatusBadRequest)
		return nil, fmt.Errorf("file too big")
	}
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		sendError(w, "INVALID_FILE", http.StatusBadRequest)
		return nil, fmt.Errorf("cant read file")
	}
	return fileBytes, nil
}

func getValAsInt(w http.ResponseWriter, r *http.Request, name string) (int64, error) {
	val := r.PostFormValue(name)
	i, err := strconv.ParseInt(val, 10, 0)
	if err != nil {
		sendError(w, "bad request", http.StatusBadRequest)
		return 0, fmt.Errorf("bad request")
	}
	return i, nil
}

func writeAnswer(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _  := template.ParseFiles("submit.html")
		t.Execute(w, nil)
		return
	}
	if err := r.ParseMultipartForm(max_upload_size); err != nil {
		fmt.Println("could not parse form data")
		sendError(w, "cant parse form", http.StatusInternalServerError)
		return
	}

	source, err := getFile(w, r, srcname)
	if err != nil { return }
	input, err := getFile(w, r, inputname)
	if err != nil { return }
	output, err := getFile(w, r, outputname)
	if err != nil { return }
	time, err := getValAsInt(w, r, timename)
	if err != nil { return }
	mem, err := getValAsInt(w, r, memname)
	if err != nil { return }

	fmt.Println(time, mem)
	id := addProblem(source, input, output, time, mem)
	
	strid := strconv.FormatInt(id, 10)
	writeAnswer(w, strid)
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("query.html")
		t.Execute(w, nil)
		return
	}

	if err := r.ParseMultipartForm(max_upload_size); err != nil {
		fmt.Println("could not parse form data")
		sendError(w, "cant parse form", http.StatusInternalServerError)
		return
	}

	id, err := getValAsInt(w, r, idname)
	if err != nil { return }

	res := getState(id)
	writeAnswer(w, res.String())
}

func Server(port int, ) {
	makedirs()

	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/query", queryHandler)

	log.Printf("Server started on localhost:%d, use /submit", port)
	addr := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func dumpRequest(r *http.Request) {
	b, _ := httputil.DumpRequest(r, true)
	fmt.Println(string(b))
}

