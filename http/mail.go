package http

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/open-falcon/mail-provider/config"
	"github.com/toolkits/http/httpclient"
	"github.com/toolkits/smtp"
	"github.com/toolkits/web/param"
)

type ApiResponse struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

func geneResp(w http.ResponseWriter, status int, msg string) {
	u := &ApiResponse{Status: status, Msg: msg}
	js, err := json.Marshal(u)
	log.Println(string(js))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func configProcRoutes() {

	http.HandleFunc("/api/sender/mail", func(w http.ResponseWriter, r *http.Request) {
		cfg := config.Config()
		token := param.String(r, "token", "")
		if cfg.Http.Token != token {
			geneResp(w, http.StatusForbidden, "no privilege")
			return
		}

		tos := param.MustString(r, "tos")
		subject := param.MustString(r, "subject")
		content := param.MustString(r, "content")
		tos = strings.Replace(tos, ",", ";", -1)

		s := smtp.New(cfg.Smtp.Addr, cfg.Smtp.Username, cfg.Smtp.Password)
		err := s.SendMail(cfg.Smtp.From, tos, subject, content)
		if err != nil {
			geneResp(w, http.StatusInternalServerError, err.Error())
		} else {
			geneResp(w, 0, "success")
		}
	})

	http.HandleFunc("/api/sender/sms", func(w http.ResponseWriter, r *http.Request) {
		cfg := config.Config()
		token := param.String(r, "token", "")
		if cfg.Http.Token != token {
			geneResp(w, http.StatusForbidden, "no privilege")
			return
		}
		tos := param.MustString(r, "tos")
		content := param.MustString(r, "content")

		transport := &httpclient.Transport{
			ConnectTimeout:   1 * time.Second,
			RequestTimeout:   5 * time.Second,
			ReadWriteTimeout: 3 * time.Second,
		}
		client := &http.Client{Transport: transport}
		defer transport.Close()
		req, _ := http.NewRequest("GET", cfg.Sms.Addr, nil)
		t := time.Now()
		q := req.URL.Query()
		q.Add("mobile", tos)
		q.Add("content", content)
		q.Add("rrid", "s"+t.Format("20060102150405"))
		q.Add("sn", cfg.Sms.Sn)
		q.Add("pwd", cfg.Sms.Pwd)
		q.Add("ext", cfg.Sms.Ext)
		q.Add("stime", cfg.Sms.Stime)
		req.URL.RawQuery = q.Encode()
		resp, err := client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			geneResp(w, http.StatusInternalServerError, "1st request failed - "+err.Error())
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		log.Println(body)
		if err != nil {
			log.Println(err.Error())
			geneResp(w, http.StatusInternalServerError, err.Error())
			return
		}
		geneResp(w, 0, "success")
	})

}
