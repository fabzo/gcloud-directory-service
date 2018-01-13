package sync

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJsonEncodeStatus(t *testing.T) {
	a := assert.New(t)

	now, err := time.Parse("2006-01-02 15:04:05 -0700 UTC", "2018-01-10 20:21:05 +0000 UTC")
	a.Nil(err)

	status := Status{}
	status.LastSync = now
	status.SyncInProgress = false

	status.LastSyncDuration = Duration{15 * time.Second}
	status.NextSync = now.Add(time.Duration(30) * time.Minute)

	b, err := json.Marshal(status)
	a.Nil(err)

	a.EqualValues(`{"last_sync":"2018-01-10T20:21:05Z","last_sync_duration":"15s","next_sync":"2018-01-10T20:51:05Z","known_groups":0,"known_users":0,"sync_in_progress":false}`, string(b))

}
