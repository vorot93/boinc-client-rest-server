package main

import (
	"encoding/xml"
	"errors"
	"testing"

	"github.com/vorot93/goutil"
)

func TestMessageXMLUnmarshal(t *testing.T) {
	var err error
	var fixture = Message{ProjectName: "ABC", Priority: 1, SeqNum: 1, Body: "Hello world! <hi>", Timestamp: int64(55555)}

	var input = "<project>ABC</project><pri>1</pri><seqno>1</seqno><body>Hello world! <hi></body><time>55555</time>"
	var result Message
	resultErr := xml.Unmarshal([]byte(input), &result)

	if resultErr != nil {
		t.Errorf(resultErr.Error())
	}

	if result != fixture {
		err = errors.New("Mismatch")
	}

	if err != nil {
		t.Errorf(goutil.ErrorOut(goutil.ErrMismatch, fixture, result))
	}
}
