package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

const (
	appId     = `57c9695026a14a0781342987c15c5d1b`
	secretKey = `h[Mb,mKY0*qRwfNPq%LWD.%NjGTPX^cUdFtBx9*7l3i]iL!%qABobCLBEQh#ESi1aaqadwZ3V6Qqx4jOtRmgk-9Hbs$I[iQpTC,nx7=zun)nSH-i@a*aTkhdqX$8I-C+`
)

func TestSecret(t *testing.T) {
	resp, err := http.PostForm("http://localhost:8000/application/login.cgi",
		url.Values{"id": {appId}, "secret": {secretKey}})
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}

func TestService(t *testing.T) {
	req, _ := http.NewRequest("POST", "http://localhost:9066/search.cgi", strings.NewReader(""))
	req.Header.Set("authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJEYXRhIjp7IkFwcElkIjoiNTdjOTY5NTAyNmExNGEwNzgxMzQyOTg3YzE1YzVkMWIiLCJTZWNyZXRLZXkiOiJoW01iLG1LWTAqcVJ3Zk5QcSVMV0QuJU5qR1RQWF5jVWRGdEJ4OSo3bDNpXWlMISVxQUJvYkNMQkVRaCNFU2kxYWFxYWR3WjNWNlFxeDRqT3RSbWdrLTlIYnMkSVtpUXBUQyxueDc9enVuKW5TSC1pQGEqYVRraGRxWCQ4SS1DKyJ9LCJleHAiOjE1NzIxMDQ3MDJ9.7jC0FOPiL4G3P79wErIJ8_jK9VOeXR_c4hVpbR6UoB4")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		t.Log("save topic failed", err.Error())
		t.FailNow()
	}
	defer resp.Body.Close()

	respByte, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(respByte))
}
