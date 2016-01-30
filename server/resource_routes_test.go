package server

import (
	"bytes"
	simplejson "github.com/bitly/go-simplejson"
	"gopkg.in/check.v1"
	"net/http"
)

func (s *S) TestCreateResourceWithNotSupportedRoot(c *check.C) {
	bufs := []string{
		`[{"name": "fail"}]`,
		`"not-valid"`,
		`10`,
	}
	for _, buf := range bufs {
		r, _ := http.NewRequest("POST", "/api/photos", bytes.NewBufferString(buf))
		response := s.Request(r)
		c.Assert(response.Code, check.Equals, http.StatusBadRequest)

		jsonErr, err := simplejson.NewFromReader(response.Body)
		c.Assert(err, check.IsNil)

		msg := jsonErr.Get("errors").GetIndex(0).Get("_all").GetIndex(0).MustString()
		c.Assert(msg, check.Equals, "Json root not is an object")
	}
}

func (s *S) TestCreateResourceWithInvalidJson(c *check.C) {
	bufs := []string{
		`["name"}`,
		`{1"adf"`,
	}
	for _, buf := range bufs {
		r, _ := http.NewRequest("POST", "/api/photos", bytes.NewBufferString(buf))
		response := s.Request(r)
		c.Assert(response.Code, check.Equals, http.StatusBadRequest)

		jsonErr, err := simplejson.NewFromReader(response.Body)
		c.Assert(err, check.IsNil)

		msg := jsonErr.Get("errors").GetIndex(0).Get("_all").GetIndex(0).MustString()
		c.Assert(msg, check.Matches, "Invalid json: .*")
	}
}

func (s *S) TestCreateResourceWithoutBody(c *check.C) {
	r, _ := http.NewRequest("POST", "/api/photos", bytes.NewBufferString(""))
	response := s.Request(r)
	c.Assert(response.Code, check.Equals, http.StatusBadRequest)

	jsonErr, err := simplejson.NewFromReader(response.Body)
	c.Assert(err, check.IsNil)

	msg := jsonErr.Get("errors").GetIndex(0).Get("_all").GetIndex(0).MustString()
	c.Assert(msg, check.Equals, "Empty resource")
}

func (s *S) TestCreateResource(c *check.C) {
	buf := bytes.NewBufferString(`{"name": "ok"}`)
	r, _ := http.NewRequest("POST", "/api/photos", buf)
	response := s.Request(r)
	c.Assert(response.Code, check.Equals, http.StatusCreated)

	json, err := simplejson.NewFromReader(response.Body)
	c.Assert(err, check.IsNil)

	c.Assert(json.Get("name").MustString(), check.Equals, "ok")
}
