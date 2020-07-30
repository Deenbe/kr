package lib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var TZ *time.Location = time.FixedZone("", 60*60*1)

func Now() time.Time {
	return time.Date(2020, 1, 2, 8, 0, 0, 0, TZ)
}

func TestParsingDuration(t *testing.T) {
	c := &Config{Rewind: "1m"}
	r, err := c.CalculatePointInTime(Now)
	assert.NoError(t, err)
	assert.Equal(t, time.Date(2020, 1, 2, 6, 59, 0, 0, time.UTC), r)
}

func TestParsingInvalidDuration(t *testing.T) {
	c := &Config{Rewind: "1msfd"}
	_, err := c.CalculatePointInTime(Now)
	assert.Error(t, err, "time: unknown unit msfd in duration 1msfd")
}

func TestParsingSince(t *testing.T) {
	c := &Config{Since: "2020-01-02T07:00:00+01:00"}
	r, err := c.CalculatePointInTime(Now)
	assert.NoError(t, err)
	assert.Equal(t, time.Date(2020, 1, 2, 6, 0, 0, 0, time.UTC), r)
}

func TestParsingInvalidSince(t *testing.T) {
	c := &Config{Since: "2020-01-02B"}
	_, err := c.CalculatePointInTime(Now)
	assert.EqualError(t, err, "unable to recognise the specified time")
}

func TestParsingUTCTime(t *testing.T) {
	c := &Config{Rewind: "1m"}
	r, err := c.CalculatePointInTime(func() time.Time { return time.Date(2020, 1, 2, 8, 0, 0, 0, time.UTC) })
	assert.NoError(t, err)
	assert.Equal(t, time.Date(2020, 1, 2, 7, 59, 0, 0, time.UTC), r)
}
