package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"net/http"
	"net/http/httptest"
)

type MockAwsSession struct {
	Session    *session.Session
	awsServer  *httptest.Server
	awsHandler http.HandlerFunc
}

func NewMockAwsSession() (*MockAwsSession, error) {
	m := &MockAwsSession{
		awsHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	}

	m.awsServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.awsHandler(w, r)
	}))

	c := aws.NewConfig()
	c.WithDisableSSL(true)
	c.WithEndpoint(m.awsServer.URL)
	c.WithRegion("us-east-1")
	c.WithCredentials(credentials.NewStaticCredentials("abc", "123", ""))
	m.Session = session.New(c)

	return m, nil
}

func (m *MockAwsSession) HandleAws(handler http.HandlerFunc) {
	m.awsHandler = handler
}

func (m *MockAwsSession) Close() {
	if m.awsServer != nil {
		m.awsServer.Close()
		m.awsServer = nil
	}
}
