package edgeauth

import (
	"crypto"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"hash"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Algo           crypto.Hash
	Key            string
	Salt           string
	FieldDelimiter string
	ACLDelimiter   string
	StartTime      time.Time
	EndTime        time.Time
	DurationWindow time.Duration
	IP             string
	SessionID      string
	Payload        string
	Verbose        bool
}

type Client struct {
	Config *Config
}

func NewClient(config *Config) *Client {
	if config.Algo == 0 {
		config.Algo = crypto.SHA256
	}

	if config.FieldDelimiter == "" {
		config.FieldDelimiter = "~"
	}

	if config.ACLDelimiter == "" {
		config.ACLDelimiter = "!"
	}

	return &Client{config}
}

func createSignature(hasher func() hash.Hash, value string, key []byte) string {
	hm := hmac.New(hasher, key)
	hm.Write([]byte(value))

	return hex.EncodeToString(hm.Sum(nil))
}

func (c *Client) GenerateToken(path string, isUrl bool) (string, error) {
	var hasher func() hash.Hash

	switch c.Config.Algo {
	case crypto.SHA256:
		hasher = sha256.New
	case crypto.SHA1:
		hasher = sha1.New
	case crypto.MD5:
		hasher = md5.New
	default:
		return "", errors.New("altorithm should be sha256 or sha1 or md5")
	}

	now := time.Now()
	startTime := c.Config.StartTime
	endTime := c.Config.EndTime

	if startTime.IsZero() {
		startTime = now
	}

	if endTime.IsZero() {
		if c.Config.DurationWindow == 0 {
			return "", errors.New("you must provide end time or duration window")
		}

		endTime = startTime.Add(c.Config.DurationWindow)
	}

	if startTime.Equal(endTime) {
		return "", errors.New("start and end time cannot be the same")
	}

	if endTime.Before(startTime) {
		return "", errors.New("end time must be greater than start time")
	}

	if endTime.Before(now) {
		return "", errors.New("end time must be in the future")
	}

	query := []string{}

	if c.Config.IP != "" {
		query = append(query, c.Config.IP)
	}

	// Include StartTime only if explicitly given
	if !c.Config.StartTime.IsZero() {
		query = append(query, "st="+strconv.FormatInt(c.Config.StartTime.Unix(), 10))
	}

	query = append(query, "exp="+strconv.FormatInt(endTime.Unix(), 10))

	if !isUrl {
		query = append(query, "acl="+path)
	}

	if c.Config.SessionID != "" {
		query = append(query, "id="+c.Config.SessionID)
	}

	if c.Config.Payload != "" {
		query = append(query, "data="+c.Config.Payload)
	}

	hashSource := make([]string, len(query))
	copy(hashSource, query)

	if isUrl {
		hashSource = append(hashSource, "url="+path)
	}

	if c.Config.Salt != "" {
		hashSource = append(hashSource, "salt="+c.Config.Salt)
	}

	key, err := hex.DecodeString(c.Config.Key)

	if err != nil {
		return "", err
	}

	token := createSignature(
		hasher,
		strings.Join(hashSource, c.Config.FieldDelimiter),
		key,
	)

	query = append(query, "hmac="+token)

	return strings.Join(query, c.Config.FieldDelimiter), nil
}
