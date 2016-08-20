package services

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Broker struct {
	s *session.Session
}

func NewBroker(s *session.Session) *Broker {
	if s == nil {
		s = session.New()
	}
	return &Broker{
		s: s,
	}
}

func (b *Broker) EC2() *ec2.EC2 {
	return ec2.New(b.s)
}
