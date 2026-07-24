package sandboxprotocol

import "testing"

func validRequest() Request {
	return Request{
		Source:           "package main\n",
		Harness:          "package main\n",
		Tests:            []ExpectedTest{{ID: "test-1", Expected: "ok"}},
		CompileTimeoutMS: 20_000,
		RuntimeTimeoutMS: 4_000,
		OutputLimitBytes: 256 * 1024,
	}
}

func TestRequestValidation(t *testing.T) {
	if err := validRequest().Validate(192 * 1024); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		edit func(*Request)
	}{
		{"missing source", func(request *Request) { request.Source = "" }},
		{"oversized harness", func(request *Request) { request.Harness = string(make([]byte, 192*1024+1)) }},
		{"duplicate tests", func(request *Request) { request.Tests = append(request.Tests, request.Tests[0]) }},
		{"compile timeout", func(request *Request) { request.CompileTimeoutMS = 31_000 }},
		{"runtime timeout", func(request *Request) { request.RuntimeTimeoutMS = 11_000 }},
		{"output limit", func(request *Request) { request.OutputLimitBytes = 256*1024 + 1 }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := validRequest()
			test.edit(&request)
			if err := request.Validate(192 * 1024); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestResponseValidationRejectsUnknownStatusesAndTests(t *testing.T) {
	expected := validRequest().Tests
	if err := (Response{Status: StatusSuccess}).Validate(expected); err != nil {
		t.Fatal(err)
	}
	if err := (Response{Status: "future"}).Validate(expected); err == nil {
		t.Fatal("expected unknown status error")
	}
	if err := (Response{
		Status:  StatusSuccess,
		Results: []TestResult{{ID: "not-requested"}},
	}).Validate(expected); err == nil {
		t.Fatal("expected unknown test error")
	}
}
