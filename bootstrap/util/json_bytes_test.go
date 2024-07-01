package util

import (
	"reflect"
	"testing"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/suite"
)

type UtilSuite struct {
	suite.Suite
}

func TestUtilSuite(t *testing.T) {
	suite.Run(t, new(UtilSuite))
}

type Message struct {
	Message string `json:"message,omitempty"`
}

func (s *UtilSuite) TestJSONBytes() {

	event := v2.NewEvent()
	event.SetID("changeme")
	event.SetSubject("changeme")
	event.SetSource("changeme")
	event.SetType("changeme")
	event.SetExtension("partitionkey", "changeme")
	event.SetData("", Message{
		Message: "changeme",
	})

	mockJson := `{"specversion":"1.0","id":"changeme","source":"changeme","type":"changeme","subject":"changeme","data":{"message":"changeme"},"partitionkey":"changeme"}`
	wantBytes := []byte(mockJson)
	gotBytes, err := JSONBytes(event)
	s.Assert().True(err == nil, "Error on call JSONBytes: %s", err)
	s.Assert().True(reflect.DeepEqual(gotBytes, wantBytes), "Error got %s want %s", gotBytes, wantBytes)
}
