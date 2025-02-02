package pkg_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"turboenigma/pkg"
)

func init() {
	pkg.Client = &pkg.MockClient{}
}

func TestPostOnSlack(t *testing.T) {
	var url = "http://turboenigma.com"

	pkg.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(`ok`))),
		}, nil
	}

	dat, err := ioutil.ReadFile("./payload/merge_request-open.json")
	pkg.Assert(err)

	os.Setenv("SLACK_WEBHOOK_URL", url)
	os.Setenv("MERGE_REQUEST_LABEL", "just-testing")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPost,
		url,
		strings.NewReader(string(dat)),
	)

	if err != nil {
		t.Fatalf("could not create request %v", err)
	}

	pkg.PostOnSlack(rec, req)

	if string(rec.Body.Bytes()) != "OK" {
		t.Errorf("rec.Body should be 'OK'; '%v' given", string(rec.Body.Bytes()))
	}
}

func TestPostOnSlackWithEmptyBody(t *testing.T) {
	var url = "http://turboenigma.com"

	os.Setenv("SLACK_WEBHOOK_URL", url)
	os.Setenv("MERGE_REQUEST_LABEL", "just-testing")

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(
		http.MethodPost,
		url,
		strings.NewReader(""),
	)

	pkg.PostOnSlack(rec, req)

	if string(rec.Body.Bytes()) != "Body is missing\n" {
		t.Errorf("Emtpy body should be validated; '%v' given on body", string(rec.Body.Bytes()))
	}
	if rec.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("Expected bad request status code; '%v' given", rec.Result().StatusCode)
	}
}

func TestPostOnSlackWithIgnoredLabel(t *testing.T) {
	var url = "http://turboenigma.com"

	dat, err := ioutil.ReadFile("./payload/merge_request-open.json")
	pkg.Assert(err)

	os.Setenv("SLACK_WEBHOOK_URL", url)
	os.Setenv("MERGE_REQUEST_LABEL", "invalid-label")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPost,
		url,
		strings.NewReader(string(dat)),
	)

	if err != nil {
		t.Fatalf("could not create request %v", err)
	}

	pkg.PostOnSlack(rec, req)

	if string(rec.Body.Bytes()) != "We didn't find the right label" {
		t.Errorf("Label should be ignored; '%v' given on body", string(rec.Body.Bytes()))
	}
}

func TestPostOnSlackWithNewIssue(t *testing.T) {
	var url = "http://turboenigma.com"

	dat, err := ioutil.ReadFile("./payload/issue-open.json")
	pkg.Assert(err)

	os.Setenv("SLACK_WEBHOOK_URL", url)
	os.Setenv("MERGE_REQUEST_LABEL", "just-testing")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPost,
		url,
		strings.NewReader(string(dat)),
	)

	if err != nil {
		t.Fatalf("could not create request %v", err)
	}

	pkg.PostOnSlack(rec, req)

	if string(rec.Body.Bytes()) != "We just care about new merge_requests" {
		t.Errorf("Event should be ignored; '%v' given on body", string(rec.Body.Bytes()))
	}
}

func TestPostOnSlackWithMergeRequestApproved(t *testing.T) {
	var url = "http://turboenigma.com"

	dat, err := ioutil.ReadFile("./payload/merge_request-approved.json")
	pkg.Assert(err)

	os.Setenv("SLACK_WEBHOOK_URL", url)
	os.Setenv("MERGE_REQUEST_LABEL", "just-testing")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPost,
		url,
		strings.NewReader(string(dat)),
	)

	if err != nil {
		t.Fatalf("could not create request %v", err)
	}

	pkg.PostOnSlack(rec, req)

	if string(rec.Body.Bytes()) != "We just care about new merge_requests" {
		t.Errorf("Event should be ignored; '%v' given on body", string(rec.Body.Bytes()))
	}
}
