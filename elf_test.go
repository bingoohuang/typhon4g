package typhon4g_test

import (
	"testing"

	"github.com/bingoohuang/typhon4g"
	"github.com/stretchr/testify/assert"
)

func TestDiffProperties(t *testing.T) {
	events := make([]typhon4g.DiffPropertiesChangeEvent, 0)
	err := typhon4g.DiffProperties("k1=v1\nk2=v2\n", "k1=v10\nk3=v3", func(event typhon4g.DiffPropertiesChangeEvent) {
		events = append(events, event)
	})

	assert.Nil(t, err)
	assert.Equal(t, []typhon4g.DiffPropertiesChangeEvent{
		{ChangeType: typhon4g.ChangeTypeModify, Key: "k1", LeftValue: "v1", RightValue: "v10"},
		{ChangeType: typhon4g.ChangeTypeDeleted, Key: "k2", LeftValue: "v2", RightValue: ""},
		{ChangeType: typhon4g.ChangeTypeNew, Key: "k3", LeftValue: "", RightValue: "v3"},
	}, events)
}
