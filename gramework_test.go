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
	app.EnableFirewall = true
	const text = "test one two three"
	var preCalled, mwCalled, postCalled, jsonOK bool
	app.GET("/", text)
	app.GET("/json", func(ctx *Context) {
		m := map[string]map[string]map[string]map[string]int{
			"abc": {
				"def": {
					"ghk": {
						"wtf": 42,
					},
				},
			},
		}

		if err := ctx.JSON(m); err == nil {
			jsonOK = true
		}

		if b, err := ctx.ToJSON(m); err == nil {
			var m2 map[string]map[string]map[string]map[string]int
			if _, err := ctx.UnJSONBytes(b, &m2); err != nil {
				jsonOK = false
				return
			}
			b2, err := ctx.ToJSON(m2)
			if err != nil || len(b2) != len(b) {
				jsonOK = false
				return
			}
			for k := range b2 {
				if v := b[k]; v != b2[k] {
					jsonOK = false
					return
				}
			}
		}
	})
	app.UsePre(func() {
		preCalled = true
	})
	app.UsePre(func(ctx *Context) {
		ctx.CORS()
	})
	app.Use(func() {
		mwCalled = true
	})
	app.UseAfterRequest(func() {
		postCalled = true
	})

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

	resp, err = http.Get("http://127.0.0.1:9977/json")
	if err != nil {
		t.Fatalf("Gramework isn't working! Got error: %s", err)
	}
	ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if !preCalled {
		t.Fatalf("pre wasn't called")
	}
	if !mwCalled {
		t.Fatalf("middleware wasn't called")
	}
	if !postCalled {
		t.Fatalf("post middleware wasn't called")
	}
	if !jsonOK {
		t.Fatalf("json response isn't OK")
	}
}

func TestGrameworkDomainHTTP(t *testing.T) {
	app := New()
	const text = "test one two three"
	app.Domain("127.0.0.1:9978").GET("/", text)
	var preCalled, mwCalled, postCalled bool
	app.UsePre(func() {
		preCalled = true
	})
	app.Use(func() {
		mwCalled = true
	})
	app.UseAfterRequest(func() {
		postCalled = true
	})

	go func() {
		app.ListenAndServe(":9978")
	}()

	time.Sleep(1 * time.Second)

	resp, err := http.Get("http://127.0.0.1:9978")
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

	if !preCalled {
		t.Fatalf("pre wasn't called")
	}
	if !mwCalled {
		t.Fatalf("middleware wasn't called")
	}
	if !postCalled {
		t.Fatalf("post middleware wasn't called")
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

func TestGrameworkListenAll(t *testing.T) {
	app := New()
	const text = "test one two three"
	app.GET("/", text)
	app.TLSEmails = []string{"k@guava.by"}

	go func() {
		app.ListenAndServeAllDev(":9449")
	}()

	time.Sleep(3 * time.Second)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get("http://127.0.0.1:9449")
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

	resp, err = client.Get("https://127.0.0.1:443")
	if err != nil {
		t.Fatalf("Gramework isn't working! Got error: %s", err)
	}
	body, err = ioutil.ReadAll(resp.Body)
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
