package httpmock

import (
	"bytes"
	"errors"
	checker "gopkg.in/check.v1"
	"io/ioutil"
	"net/http"
	"testing"
)

// two helper funcs here
func GetResponseWithBody(body string) *http.Response {
	response := &http.Response{}
	response.Body = ioutil.NopCloser(bytes.NewBufferString(body))
	response.ContentLength = int64(len(body))
	return response
}

func BodyToString(response *http.Response) string {
	defer response.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	return string(bodyBytes)
}

func Test(t *testing.T) { checker.TestingT(t) }

type HttpmockSuite struct{}

var _ = checker.Suite(&HttpmockSuite{})

func (s *HttpmockSuite) SetUpTest(c *checker.C) {
	Activate()
}

func (s *HttpmockSuite) TearDownTest(c *checker.C) {
	Deactivate()
}

func (suite *HttpmockSuite) TestActivate_NoResponder_Error(c *checker.C) {
	response, err := http.Get("http://example.com/")
	c.Assert(err.Error(), checker.Equals, "Get http://example.com/: no responder found")
	c.Assert(response, checker.Equals, (*http.Response)(nil))
}

func (suite *HttpmockSuite) TestActivate_ResponderOk_Same(c *checker.C) {
	RegisterResponder("GET", "http://example.com/", func(*http.Request) (*http.Response, error) {
		return GetResponseWithBody("xin chao"), nil
	})

	response, _ := http.Get("http://example.com/")

	c.Assert(BodyToString(response), checker.Equals, "xin chao")
}

func (suite *HttpmockSuite) TestActivate_ResponderError_Same(c *checker.C) {
	RegisterResponder("GET", "http://example.com/", func(*http.Request) (*http.Response, error) {
		return nil, errors.New("what")
	})

	_, err := http.Get("http://example.com/")

	c.Assert(err.Error(), checker.Equals, "Get http://example.com/: what")
}

func (suite *HttpmockSuite) TestActivate_ResponderDifferentGetNoFail_Error(c *checker.C) {
	RegisterResponder("GET", "http://another.com/", func(*http.Request) (*http.Response, error) {
		return GetResponseWithBody("xin chao"), nil
	})

	suite.TestActivate_NoResponder_Error(c)
}

var impossibleUrl = "localhost.bla.bla.bla"

func (suite *HttpmockSuite) TestActivate_ResponderImpossibleUrl_Same(c *checker.C) {
	RegisterResponder("GET", impossibleUrl, func(*http.Request) (*http.Response, error) {
		return GetResponseWithBody("xin chao"), nil
	})

	response, _ := http.Get(impossibleUrl)

	c.Assert(BodyToString(response), checker.Equals, "xin chao")
}

func (suite *HttpmockSuite) TestDeactivate_ResponderImpossibleUrl_Same(c *checker.C) {
	Deactivate()
	_, err := http.Get(impossibleUrl)
	c.Assert(err.Error(), checker.Equals, "Get localhost.bla.bla.bla: unsupported protocol scheme \"\"")
}
