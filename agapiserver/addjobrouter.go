// https://github.com/gorilla/websocket/blob/master/examples/echo/server.go
// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package agapiserver

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/xtforgame/agak/requestsender"
	"github.com/xtforgame/agak/scheduler"
	"github.com/xtforgame/agak/utils"
	"io/ioutil"
	"net/http"
	"time"
)

type JobSave struct {
	Name     string `json:"name,"`
	Owner    string `json:"owner,"`
	DateTime string `json:"dateTime,"`
	CronExp  string `json:"cronExp,"`
	RunOnce  bool   `json:"runOnce,"`
	Disabled bool   `json:"disabled,"`
}

func newSchedule(
	mainScheduler *scheduler.Scheduler,
	jobName string,
	spec string,
	runFunc func(*scheduler.Job, *scheduler.Entry),
) (*scheduler.Entry, error) {
	job := &scheduler.Job{Name: jobName}
	entry, err := mainScheduler.AddEntry(spec, job, func(ent *scheduler.Entry) {
		job.RunFunc = func() {
			e := ent
			runFunc(e.GetJob(), e)
		}
	})

	return entry, err
}

func AddJobRouter(
	router *chi.Mux,
	reqSender *requestsender.RequestSender,
	mainScheduler *scheduler.Scheduler,
) {
	router.Get("/jobs", func(w http.ResponseWriter, r *http.Request) {
		if jsonBytes, err := json.Marshal(map[string]string{}); err == nil {
			w.Write(jsonBytes)
			return
		}
		w.Write([]byte("{}"))
	})

	router.Post("/jobs", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte("{}"))
			return
		}
		js := &JobSave{}
		if err := json.Unmarshal(body, js); err != nil {
			w.Write([]byte("{}"))
			return
		}
		fmt.Println("js :", js)

		t2, err := time.Parse("2006-01-02T15:04:05.000Z", js.DateTime)
		if err != nil {
			return
		}
		fmt.Println("t2 :", t2)
		t1 := time.Now()
		t1Ms := utils.ToMillisecond(t1)

		t2Ms := utils.ToMillisecond(t2)

		utils.SetTimeout(func() {
			fmt.Println("utils.SetTimeout(func() {")
			res, err := reqSender.SendRequest(
				nil,
				&requestsender.RequestConfig{
					Method: "POST",
					Url:    "https://httpbin.org/post",
					Header: map[string]string{
						"Content-Type": "application/json",
					},
					Body:      []byte(`["3338"]`),
					Validator: requestsender.DefaultResponseValidator,
				},
			)
			if err == nil {
				fmt.Println("res", string(res.Body))
			}
		}, int(t2Ms-t1Ms))

		// fmt.Println("body :", string(body))
		// entry, err := newSchedule(
		// 	mainScheduler,
		// 	js.Name,
		// 	"@every 2s",
		// 	func(job *scheduler.Job, entry *scheduler.Entry) {
		// 		mainScheduler.RemoveEntry(entry)
		// 		fmt.Println("Job Name :", js.Name)
		// 		res, err := reqSender.SendRequest(
		// 			nil,
		// 			&requestsender.RequestConfig{
		// 				Method: "POST",
		// 				Url:    "https://httpbin.org/post",
		// 				Header: map[string]string{
		// 					"Content-Type": "application/json",
		// 				},
		// 				Body:      []byte(`["3338"]`),
		// 				Validator: requestsender.DefaultResponseValidator,
		// 			},
		// 		)
		// 		if err == nil {
		// 			fmt.Println("res", string(res.Body))
		// 		}
		// 	},
		// )
		if err != nil {
			w.Write([]byte("{}"))
			return
		}
		// idStr := fmt.Sprintf("%d", entry.GetEntryID())
		// if jsonBytes, err := json.Marshal(map[string]interface{}{"id": entry.GetEntryID()}); err == nil {
		// 	w.Write(jsonBytes)
		// 	return
		// }
		w.Write([]byte("{}"))
	})
}
