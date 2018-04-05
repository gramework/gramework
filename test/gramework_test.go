package test

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/gramework/gramework"
)

func TestGrameworkHTTP(t *testing.T) {
	app := gramework.New()
	const text = "test one two three"
	var preCalled, mwCalled, postCalled, jsonOK bool
	app.GET("/", text)
	app.GET("/bytes", []byte(text))
	app.GET("/float32", float32(len(text)))
	app.GET("/dumb", func() {})
	app.GET("/dumbWithErr", func() error { return nil })
	app.GET("/json", func(ctx *gramework.Context) {
		m := map[string]map[string]map[string]map[string]int{
			"abc": {
				"def": {
					"ghk": {
						"wtf": 42,
					},
				},
			},
		}
		jsonOK = true

		if err := ctx.JSON(m); err != nil {
			jsonOK = false
			ctx.Logger.Errorf("can't JSON(): %s", err)
		}

		if b, err := ctx.ToJSON(m); err == nil {
			var m2 map[string]map[string]map[string]map[string]int
			if _, err := ctx.UnJSONBytes(b, &m2); err != nil {
				ctx.Logger.Errorf("can't unjson: %s", err)
				jsonOK = false
				return
			}
			b2, err := ctx.ToJSON(m2)
			if err != nil {
				ctx.Logger.Errorf("ToJSON returns error: %s", err)
				jsonOK = false
				return
			}
			if len(b2) != len(b) {
				ctx.Logger.Errorf("len is not equals, got len(b2) = [%v], len(b) = [%v]", len(b2), len(b))
				jsonOK = false
				return
			}
			if !reflect.DeepEqual(b, b2) {
				jsonOK = false
				return
			}
		}
	})

	app.ServeFile("/sf", "./nanotime.s")
	app.SPAIndex("./nanotime.s")
	app.GET("/sdnc_static/dist/*static", app.ServeDirNoCache("./"))
	app.GET("/sdncc_static/dist/*static", app.ServeDirNoCacheCustom("./", 0, false, false, []string{}))
	app.MethodNotAllowed(func(ctx *gramework.Context) {
		if _, err := ctx.WriteString("GTFO"); err != nil {
			t.Error(err.Error())
		}
	})

	var err error
	if err = app.UsePre(func() {
		preCalled = true
	}); err != nil {
		t.Error(err.Error())
	}

	if err = app.UsePre(func(ctx *gramework.Context) {
		ctx.CORS()
	}); err != nil {
		t.Error(err.Error())
	}

	if err = app.Use(func() {
		mwCalled = true
	}); err != nil {
		t.Error(err.Error())
	}

	if err = app.UseAfterRequest(func() {
		postCalled = true
	}); err != nil {
		t.Error(err.Error())
	}

	go func() {
		if err := app.ListenAndServe(":9977"); err != nil {
			t.Error(err.Error())
		}
	}()

	time.Sleep(2 * time.Second)

	resp, err := http.Get("http://127.0.0.1:9977")
	if err != nil {
		t.Fatalf("Gramework isn't working! Got error: %s", err)
		t.FailNow()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Gramework isn't working! Can't read body: %s", err)
		t.FailNow()
	}

	if err = resp.Body.Close(); err != nil {
		t.Error(err.Error())
	}

	if string(body) != text {
		t.Fatalf(
			"Gramework returned unexpected body! Got %q, expected %q",
			string(body),
			text,
		)
		t.FailNow()
	}

	resp, err = http.Get("http://127.0.0.1:9977/json")
	if err != nil {
		t.Fatalf("Gramework isn't working! Got error: %s", err)
		t.FailNow()
	}

	if _, err = ioutil.ReadAll(resp.Body); err != nil {
		t.Error(err.Error())
	}

	if err = resp.Body.Close(); err != nil {
		t.Error(err.Error())
	}

	switch {
	case !preCalled:
		t.Fatalf("pre wasn't called")
		t.FailNow()
	case !mwCalled:
		t.Fatalf("middleware wasn't called")
		t.FailNow()
	case !postCalled:
		t.Fatalf("post middleware wasn't called")
		t.FailNow()
	case !jsonOK:
		t.Fatalf("json response isn't OK")
		t.FailNow()
	}
}

func TestGrameworkDomainHTTP(t *testing.T) {
	app := gramework.New()
	const text = "test one two three"
	app.Domain("127.0.0.1:9978").GET("/", text)
	var preCalled, mwCalled, postCalled bool
	var err error

	if err = app.UsePre(func() {
		preCalled = true
	}); err != nil {
		t.Error(err.Error())
	}

	if err = app.Use(func() {
		mwCalled = true
	}); err != nil {
		t.Error(err.Error())
	}

	if err = app.UseAfterRequest(func() {
		postCalled = true
	}); err != nil {
		t.Error(err.Error())
	}

	go func() {
		if err := app.ListenAndServe(":9978"); err != nil {
			t.Error(err.Error())
		}
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

	if err = resp.Body.Close(); err != nil {
		t.Error(err.Error())
	}

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
	app := gramework.New()
	const text = "test one two three"
	app.GET("/", text)
	app.TLSEmails = []string{"k@guava.by"}

	go func() {
		if err := app.ListenAndServeAutoTLSDev(":9443"); err != nil {
			t.Error(err.Error())
		}
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

	if err = resp.Body.Close(); err != nil {
		t.Error(err.Error())
	}

	if string(body) != text {
		t.Fatalf(
			"Gramework returned unexpected body! Got %q, expected %q",
			string(body),
			text,
		)
	}
}

func TestGrameworkListenAll(t *testing.T) {
	app := gramework.New()
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

	if err = resp.Body.Close(); err != nil {
		t.Error(err.Error())
	}

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

	if err = resp.Body.Close(); err != nil {
		t.Error(err.Error())
	}

	if string(body) != text {
		t.Fatalf(
			"Gramework returned unexpected body! Got %q, expected %q",
			string(body),
			text,
		)
	}
}
