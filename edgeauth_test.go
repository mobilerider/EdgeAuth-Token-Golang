package edgeauth

import (
	"crypto"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

const sampleKey = "52a152a152a152a152a152a152a1"
const samplePath = "/this/is/a/test"

func TestGenerateToken(t *testing.T) {
	config := &Config{
		Algo:           crypto.SHA256,
		Key:            sampleKey,
		StartTime:      time.Now(),
		DurationWindow: 300 * time.Second,
	}

	client := NewClient(config)

	token, err := client.GenerateToken(samplePath, false)

	if err != nil {
		t.Error(err.Error())
	}

	fields := strings.Split(token, config.FieldDelimiter)

	if len(fields) != 4 {
		t.Error("there should be 4 fields in the token")
	}

	expected := "st=" + strconv.FormatInt(config.StartTime.Unix(), 10)

	if expected != fields[0] {
		t.Errorf("first field must be equal to `%s`", expected)
	}

	endDate := config.StartTime.Add(config.DurationWindow)
	expected = "exp=" + strconv.FormatInt(endDate.Unix(), 10)

	if expected != fields[1] {
		t.Errorf("first field must be equal to `%s`", expected)
	}

	matched, _ := regexp.MatchString(`acl=.+`, fields[2])

	if !matched {
		t.Error("second field must in the form `acl=<path>`")
	}

	matched, _ = regexp.MatchString(`hmac=\w{64}`, fields[3])

	if !matched {
		t.Error("third field must in the form `hmac=<hash>`")
	}

	t.Log("Token: " + token)
}

func TestGenerateTokenWithInvalidStartAndEndDate(t *testing.T) {
	expected := "end time must be greater than start time"

	config := &Config{
		Algo:      crypto.SHA256,
		Key:       sampleKey,
		StartTime: time.Now().Add(300 * time.Second),
		EndTime:   time.Now(),
	}

	client := NewClient(config)

	_, err := client.GenerateToken(samplePath, false)

	if err != nil {
		if err.Error() != expected {
			t.Error("error must be " + expected)
		}
	} else {
		t.Error("error must be returned")
	}
}
