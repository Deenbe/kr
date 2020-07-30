/*
Copyright © 2020 kr contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package lib

import (
	"time"

	"github.com/pkg/errors"
)

type Config struct {
	StreamName   string
	ConsumerName string
	Rewind       string
	Since        string
	Update       bool
}

func (c *Config) CalculatePointInTime(now func() time.Time) (time.Time, error) {
	t, err := c.parseTime(now)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(time.UTC), nil
}

func (c *Config) parseTime(now func() time.Time) (time.Time, error) {
	if c.Rewind != "" {
		d, err := time.ParseDuration(c.Rewind)
		if err != nil {
			return time.Time{}, err
		}
		return now().Add(time.Duration(-1) * d), nil
	} else if c.Since != "" {
		formats := []string{
			"2006-01-02",
			"2006-01-02 15:04",
			"2006-01-02 15:04:05",
			time.RFC3339,
			time.RFC1123,
			time.RFC1123Z,
		}

		for _, f := range formats {
			var t time.Time
			t, err := time.Parse(f, c.Since)
			if err == nil {
				// If tz is not specified in the input assume the user is referring to local time
				_, o := t.Zone()
				if o == 0 {
					t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
				}
				return t, nil
			}
		}

		return time.Time{}, errors.WithStack(errors.New("unable to recognise the specified time"))
	}

	return time.Time{}, errors.WithStack(errors.New("either rewind or since must be specified"))
}
