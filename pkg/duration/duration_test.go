package duration

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDuration(t *testing.T) {
	data, err := json.Marshal(Duration(time.Minute + time.Second*15))
	assert.NoError(t, err)
	assert.Equal(t, `"1m15s"`, string(data))

	var duration Duration
	err = json.Unmarshal([]byte(`"1m15s"`), &duration)
	assert.NoError(t, err)
	assert.Equal(t, Duration(time.Minute+time.Second*15), duration)
}
