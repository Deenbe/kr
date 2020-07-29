package lib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var TZ *time.Location = time.FixedZone("", 60*60*10)

func Now() time.Time {
	return time.Date(2020, 1, 2, 8, 0, 0, 0, TZ)
}

func TestParsingDuration(t *testing.T) {
	c := &Config{Rewind: "1m"}

	r, err := c.CalculatePointInTime(Now)
	assert.NoError(t, err)
	assert.Equal(t, r, time.Date(2020, 1, 2, 7, 59, 0, 0, TZ))
}

func TestParsingInvalidDuration(t *testing.T) {
	c := &Config{Rewind: "1msfd"}
	_, err := c.CalculatePointInTime(Now)
	assert.Error(t, err, "time: unknown unit msfd in duration 1msfd")
}

func TestParsingSince(t *testing.T) {
	c := &Config{Since: "2020-01-02T07:00:00+10:00"}
	r, err := c.CalculatePointInTime(Now)
	assert.NoError(t, err)
	assert.Equal(t, r, time.Date(2020, 1, 2, 7, 0, 0, 0, TZ))
}

func TestParsingInvalidSince(t *testing.T) {
	c := &Config{Since: "2020-01-02B"}
	_, err := c.CalculatePointInTime(Now)
	assert.EqualError(t, err, "unable to recognise the specified time")
}
