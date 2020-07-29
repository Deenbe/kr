package lib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Now() time.Time {
	return time.Date(2020, 1, 2, 8, 0, 0, 0, time.Local)
}

func TestParsingDuration(t *testing.T) {
	c := &Config{"a", "a", "1m", ""}
	r, err := c.CalculatePointInTime(Now)
	assert.NoError(t, err)
	assert.Equal(t, r, time.Date(2020, 1, 2, 7, 59, 0, 0, time.Local))
}

func TestParsingInvalidDuration(t *testing.T) {
	c := &Config{"a", "a", "1msfd", ""}
	_, err := c.CalculatePointInTime(Now)
	assert.Error(t, err, "time: unknown unit msfd in duration 1msfd")
}

func TestParsingSince(t *testing.T) {
	c := &Config{"a", "a", "", "2020-01-02T07:00:00+11:00"}
	r, err := c.CalculatePointInTime(Now)
	assert.NoError(t, err)
	assert.Equal(t, r, time.Date(2020, 1, 2, 7, 0, 0, 0, time.Local))
}

func TestParsingInvalidSince(t *testing.T) {
	c := &Config{"a", "a", "", "2020-01-02"}
	_, err := c.CalculatePointInTime(Now)
	assert.EqualError(t, err, "unable to recognise the specified time")
}
