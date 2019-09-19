// https://github.com/gorilla/websocket/blob/master/examples/echo/server.go
// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package agapiserver

import (
	"bytes"
	"net/http"
	// "sort"
	"encoding/json"
	"github.com/xtforgame/cmdraida/crcore"
	// "github.com/xtforgame/cmdraida/t1"
	"os"

	"archive/zip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"strings"
	"time"
)

// https://golangcode.com/unzip-files-in-go/
func unzip(files []*zip.File, dest string) ([]string, error) {
	var filenames []string
	for _, f := range files {

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}
		defer rc.Close()

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {

			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)

		} else {

			// Make File
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return filenames, err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return filenames, err
			}

			_, err = io.Copy(outFile, rc)

			// Close the file without defer to close before next iteration of loop
			outFile.Close()

			if err != nil {
				return filenames, err
			}

		}
	}
	return filenames, nil
}

func Unzip(src string, dest string) ([]string, error) {
	var filenames []string
	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()
	return unzip(r.File, dest)
}

func UnzipFromBytes(b []byte, dest string) ([]string, error) {
	var filenames []string
	r, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return filenames, err
	}
	return unzip(r.File, dest)
}

// func CreateBuildHandler(hs *HttpServer) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		targetFolder := "/usr/azgoapi/runtime/unzipped/wasm_test"
// 		os.RemoveAll(targetFolder)
// 		os.MkdirAll(targetFolder, os.ModePerm)
// 		_, err := Unzip("/usr/azgoapi/examples/wasm_test.zip", targetFolder)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		task := hs.taskManager.RunTask(crcore.CommandType{
// 			Command: "go",
// 			Args:    []string{"build", "-v", "-mod=vendor", "-o", "/usr/azgoapi/runtime/build/main.wasm", "wasm_01.go"},
// 			Timeouts: crcore.TimeoutsType{
// 				Proccess:    50000,
// 				AfterKilled: 10000,
// 			},
// 			Env: append(os.Environ(), "GOOS=js", "GOARCH=wasm"),
// 			Dir: filepath.Join(targetFolder, "wasm_test"),
// 		})
// 		if jsonBytes, err := json.Marshal(task.ResultLog()); err == nil {
// 			w.Write(jsonBytes)
// 			return
// 		}
// 		w.Write([]byte("[]"))
// 	}
// }

// https://zupzup.org/go-http-file-upload-download/
// https://github.com/zupzup/golang-http-file-upload-download/blob/master/main.go
func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func initRand() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type BuildResponse struct {
	Wasm   string            `json:"wasm,"`
	Result *crcore.ResultLog `json:"result,"`
}

const maxUploadSize = 2 * 1024 * 1024 // 2 mb
// const maxUploadSize = 2

func CreateBuildHandler(hs *HttpServer) func(w http.ResponseWriter, r *http.Request) {
	initRand()
	return func(w http.ResponseWriter, r *http.Request) {
		var targetFolder = "/usr/azgoapi/runtime/unzipped/" + RandStringRunes(16)
		_, err := os.Stat(targetFolder)
		for err == nil {
			// file or dir exists
			targetFolder = "/usr/azgoapi/runtime/unzipped/" + RandStringRunes(16)
			_, err = os.Stat(targetFolder)
		}
		targetSrcFolder := filepath.Join(targetFolder, "src")
		buildTarget := filepath.Join(targetFolder, "build", "main.wasm")

		// os.RemoveAll(targetFolder)
		os.MkdirAll(targetSrcFolder, os.ModePerm)
		// defer func () {
		// 	os.RemoveAll(targetFolder)
		// }()

		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
			return
		}

		// parse and validate file and post parameters
		// entryFile := r.PostFormValue("entry")
		// fmt.Println("entryFile :", entryFile)
		file, _, err := r.FormFile("file")
		if err != nil {
			fmt.Println("r.FormFile(file)")
			// fmt.Println("err :", err)
			renderError(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println("ioutil.ReadAll(file)")
			// fmt.Println("err :", err)
			renderError(w, err.Error(), http.StatusBadRequest)
			return
		}

		filenames, err := UnzipFromBytes(fileBytes, targetSrcFolder)
		if err != nil {
			// fmt.Println("fileBytes :", fileBytes)
			fmt.Println("UnzipFromBytes")
			// fmt.Println("err :", err)
			renderError(w, err.Error(), http.StatusBadRequest)
			return
		}
		const goModName = "/go.mod"
		var projectFolder = ""
		var lastSlashIndex = 0
		for _, filename := range filenames {
			startPoint := len(filename) - len(goModName)
			if startPoint < 0 {
				continue
			}
			if filename[startPoint:] == goModName {
				if projectFolder == "" || lastSlashIndex > startPoint {
					projectFolder = filename[:startPoint]
					lastSlashIndex = startPoint
				}
			}
		}

		if projectFolder == "" {
			renderError(w, "'go.mod' not found", http.StatusBadRequest)
			return
		}

		projectFolder, _ = filepath.Abs(projectFolder)

		// _, err = Unzip("/usr/azgoapi/examples/wasm_test.zip", targetSrcFolder)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		task := hs.taskManager.RunTask(crcore.CommandType{
			Command: "go",
			Args: []string{
				"build",
				"-v", "-mod=vendor",
				"-gcflags=-trimpath=" + projectFolder,
				"-asmflags=-trimpath=" + projectFolder,
				"-o", buildTarget /*, entryFile */},
			Timeouts: crcore.TimeoutsType{
				Proccess:    50000,
				AfterKilled: 10000,
			},
			Env: append(os.Environ(), "GOOS=js", "GOARCH=wasm"),
			Dir: projectFolder,
		})
		var wasmBase64 string
		wasmBin, err := ioutil.ReadFile(buildTarget)
		if err == nil {
			wasmBase64 = base64.StdEncoding.EncodeToString(wasmBin)
		}
		if jsonBytes, err := json.Marshal(&BuildResponse{
			Wasm:   wasmBase64,
			Result: task.ResultLog(),
		}); err == nil {
			w.Write(jsonBytes)
			return
		}
		w.Write([]byte("[]"))
	}
}
