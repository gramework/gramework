package gramework

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestGrameworkHTTP(t *testing.T) {
	app := New()
	const text = "test one two three"
	app.GET("/", text)

	go func() {
		app.ListenAndServe(":9977")
	}()

	time.Sleep(1 * time.Second)

	resp, err := http.Get("http://127.0.0.1:9977")
	if err != nil {
		t.Fatalf("Gramework isn't working! Got error: %s", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Gramework isn't working! Can't read body: %s", err)
	}
	resp.Body.Close()
	if string(body) != text {
		t.Fatalf(
			"Gramework returned unexpected body! Got %q, expected %q",
			string(body),
			text,
		)
	}
}

func TestGrameworkHTTPS(t *testing.T) {
	app := New()
	const text = "test one two three"
	app.GET("/", text)
	app.TLSEmails = []string{"k@guava.by"}

	go func() {
		app.ListenAndServeAutoTLSDev(":9443")
	}()

	time.Sleep(3 * time.Second)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://127.0.0.1:9443")
	if err != nil {
		t.Fatalf("Gramework isn't working! Got error: %s", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Gramework isn't working! Can't read body: %s", err)
	}
	resp.Body.Close()
	if string(body) != text {
		t.Fatalf(
			"Gramework returned unexpected body! Got %q, expected %q",
			string(body),
			text,
		)
	}
}
