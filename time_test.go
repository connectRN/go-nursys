package nursys

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testTime struct {
	DateTime Time `json:"date_time"`
}

func Test_MarshalJSON(t *testing.T) {

	EST, err := time.LoadLocation("America/New_York")
	assert.NoError(t, err)

	t.Run("MarshalJSON", func(t *testing.T) {
		dateTime := time.Date(2021, 1, 2, 3, 4, 5, 1000, EST)

		t1 := testTime{DateTime: Time(dateTime)}
		s1, err := json.Marshal(t1)
		assert.NoError(t, err)
		assert.Equal(t, `{"date_time":"2021-01-02T03:04:05.000001-05:00"}`, string(s1))
	})

}

func Test_UnmarshalJSON(t *testing.T) {

	EST, err := time.LoadLocation("America/New_York")
	assert.NoError(t, err)
	dateTime := time.Date(2021, 1, 2, 3, 4, 5, 1000, EST)

	t.Run("RFC3339", func(t *testing.T) {
		s := []byte(`{"date_time":"2021-01-02T03:04:05.000001-05:00"}`)
		var t1 testTime
		err := json.Unmarshal(s, &t1)
		assert.NoError(t, err)
		assert.Equal(t, dateTime.UTC(), time.Time(t1.DateTime).UTC())
	})

	t.Run("Non-strict RFC3339", func(t *testing.T) {
		s := []byte(`{"date_time":"2021-01-02 03:04:05.000001-05:00"}`)
		var t1 testTime
		err := json.Unmarshal(s, &t1)
		assert.NoError(t, err)
		assert.Equal(t, dateTime.UTC(), time.Time(t1.DateTime).UTC())
	})

	t.Run("Custom format without time zone", func(t *testing.T) {
		s := []byte(`{"date_time":"2021-07-08T11:34:55"}`)
		var t1 testTime
		err := json.Unmarshal(s, &t1)
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2021, 7, 8, 11, 34, 55, 0, time.UTC), time.Time(t1.DateTime).UTC())
	})

	t.Run("Custom format with time zone", func(t *testing.T) {
		s := []byte(`{"date_time":"2021-07-08T11:34:55-07:00"}`)
		var t1 testTime
		err := json.Unmarshal(s, &t1)
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2021, 7, 8, 11, 34, 55, 0, time.FixedZone("", -7*60*60)), time.Time(t1.DateTime))
	})
}
