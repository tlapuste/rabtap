// Copyright (C) 2017 Jan Delgado

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testMessage used troughout tests
var testMessage = &amqp.Delivery{
	Exchange:        "exchange",
	RoutingKey:      "routingkey",
	Priority:        99,
	Expiration:      "2017-05-22 17:00:00",
	ContentType:     "plain/text",
	ContentEncoding: "utf-8",
	MessageId:       "4711",
	Timestamp:       time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	Type:            "some type",
	CorrelationId:   "4712",
	Headers:         amqp.Table{"header": "value"},
	AppId:           "123",
	UserId:          "456",
	Body:            []byte("simple test message."),
}

// TestSaveMessageToFiles tests the SaveMessagesToFiles() function by
// writing to and reading from temporary files.
func TestSaveMessageToRawFile(t *testing.T) {
	testdir, err := ioutil.TempDir("", "")
	require.Nil(t, err)
	defer os.RemoveAll(testdir)

	// SaveMessagesToFiles() will create files "test.dat" and "test.json" in
	// testdir.
	basename := filepath.Join(testdir, "test")
	err = SaveMessageToRawFile(basename, testMessage)
	assert.Nil(t, err)

	// check contents of message body .dat file
	datFilename := basename + ".dat"
	contentsBody, err := ioutil.ReadFile(datFilename)
	assert.Nil(t, err)
	assert.Equal(t, []byte("simple test message."), contentsBody)

	// check contents of metadata file
	metaFilename := basename + ".json"
	contentsMeta, err := ioutil.ReadFile(metaFilename)
	assert.Nil(t, err)
	// deserialize from .json file
	var jsonMetaActual RabtapPersistentMessage
	err = json.Unmarshal(contentsMeta, &jsonMetaActual)
	assert.Nil(t, err)

	// test some of the attributes
	assert.Equal(t, testMessage.AppId, jsonMetaActual.AppID)
	assert.Equal(t, len(testMessage.Headers), len(jsonMetaActual.Headers))
	assert.Equal(t, testMessage.Headers["header"], jsonMetaActual.Headers["header"])
	assert.Equal(t, testMessage.Timestamp, jsonMetaActual.Timestamp)
}

func TestSaveMessageToFilesToInvalidDir(t *testing.T) {
	// use nonexisting path
	filename := filepath.Join("/thispathshouldnotexist", "test")
	err := SaveMessageToRawFile(filename, testMessage)
	assert.NotNil(t, err)
}

// TestSaveMessageToFile tests the SaveMessagesToFile() function by
// writing to and reading a temporary files.
func TestSaveMessageToJSONFile(t *testing.T) {
	testdir, err := ioutil.TempDir("", "")
	assert.Nil(t, err)
	defer os.RemoveAll(testdir)

	filename := filepath.Join(testdir, "test")
	err = SaveMessageToJSONFile(filename, testMessage)
	assert.Nil(t, err)

	contents, err := ioutil.ReadFile(filename)
	assert.Nil(t, err)
	// deserialize from .json file
	var jsonActual RabtapPersistentMessage
	err = json.Unmarshal(contents, &jsonActual)
	assert.Nil(t, err)

	assert.Equal(t, testMessage.AppId, jsonActual.AppID)
	assert.Equal(t, len(testMessage.Headers), len(jsonActual.Headers))
	assert.Equal(t, testMessage.Headers["header"], jsonActual.Headers["header"])
	assert.Equal(t, testMessage.Timestamp, jsonActual.Timestamp)
	assert.Equal(t, []byte("simple test message."), jsonActual.Body)
}

func TestSaveMessageToFileToInvalidDir(t *testing.T) {
	// use nonexisting path
	filename := filepath.Join("/thispathshouldnotexist", "test")
	err := SaveMessageToJSONFile(filename, testMessage)
	assert.NotNil(t, err)
}

func TestCreateTimestampFilename(t *testing.T) {
	tm := time.Date(2009, time.November, 10, 23, 1, 2, 3, time.UTC)
	filename := CreateTimestampFilename(tm)
	assert.Equal(t, "2009-11-10T23_01_02.000000003Z", filename)
}

func ExampleWriteMessageBodyBlob() {
	var testMessage = &amqp.Delivery{
		Body: []byte("simple test message."),
	}
	err := WriteMessageBodyBlob(os.Stdout, testMessage)
	if err != nil {
		log.Fatal(err)
	}

	// Output:
	// simple test message.

}

func ExampleWriteMessageJSON_withBody() {

	// serialize with message body, Body will be base64 encoded.
	err := WriteMessageJSON(os.Stdout, true /* w/ body*/, testMessage)
	if err != nil {
		log.Fatal(err)
	}

	// Output:
	// {
	//   "Headers": {
	//     "header": "value"
	//   },
	//   "ContentType": "plain/text",
	//   "ContentEncoding": "utf-8",
	//   "DeliveryMode": 0,
	//   "Priority": 99,
	//   "CorrelationID": "4712",
	//   "ReplyTo": "",
	//   "Expiration": "2017-05-22 17:00:00",
	//   "MessageID": "4711",
	//   "Timestamp": "2009-11-10T23:00:00Z",
	//   "Type": "some type",
	//   "UserID": "456",
	//   "AppID": "123",
	//   "DeliveryTag": 0,
	//   "Redelivered": false,
	//   "Exchange": "exchange",
	//   "RoutingKey": "routingkey",
	//   "Body": "c2ltcGxlIHRlc3QgbWVzc2FnZS4="
	// }
}

func ExampleWriteMessageJSON_withoutBody() {
	err := WriteMessageJSON(os.Stdout, false /*w/o body*/, testMessage)
	if err != nil {
		log.Fatal(err)
	}

	// Output:
	// {
	//   "Headers": {
	//     "header": "value"
	//   },
	//   "ContentType": "plain/text",
	//   "ContentEncoding": "utf-8",
	//   "DeliveryMode": 0,
	//   "Priority": 99,
	//   "CorrelationID": "4712",
	//   "ReplyTo": "",
	//   "Expiration": "2017-05-22 17:00:00",
	//   "MessageID": "4711",
	//   "Timestamp": "2009-11-10T23:00:00Z",
	//   "Type": "some type",
	//   "UserID": "456",
	//   "AppID": "123",
	//   "DeliveryTag": 0,
	//   "Redelivered": false,
	//   "Exchange": "exchange",
	//   "RoutingKey": "routingkey",
	//   "Body": ""
	// }
}
