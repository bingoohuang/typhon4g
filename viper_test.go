package typhon4g

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/bingoohuang/typhon4g/base"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestUnmarshalToml(t *testing.T) {
	a, b, ok := unmarshal(toml.Unmarshal, base.ConfFileChangeEvent{
		Old:     `foo="bar"`,
		Current: `foo="bar"`,
	})

	assert.True(t, ok)
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, a)
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, b)
}

func TestUnmarshalYaml(t *testing.T) {
	a, b, ok := unmarshal(yaml.Unmarshal, base.ConfFileChangeEvent{
		Old:     `foo: bar`,
		Current: `foo: bar`,
	})

	assert.True(t, ok)
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, a)
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, b)
}

func TestDiffMap(t *testing.T) {
	{
		m1 := map[string]string{"foo": "bar"}
		m2 := map[string]string{"foo": "bar"}

		diffItems := DiffMap(m1, m2)
		assert.Equal(t, []DiffItem{
			{
				ChangeType: Same,
				Key:        "foo",
				LeftValue:  "bar",
				RightValue: "bar",
			},
		}, diffItems)
	}

	{
		m1 := map[string]string{"foo": "bar1"}
		m2 := map[string]string{"foo": "bar2"}

		diffItems := DiffMap(m1, m2)
		assert.Equal(t, []DiffItem{
			{
				ChangeType: Modified,
				Key:        "foo",
				LeftValue:  "bar1",
				RightValue: "bar2",
			},
		}, diffItems)
	}

	{
		m1 := map[string]string{"foo": "bar"}
		m2 := map[string]string{}

		diffItems := DiffMap(m1, m2)
		assert.Equal(t, []DiffItem{
			{
				ChangeType: Removed,
				Key:        "foo",
				LeftValue:  "bar",
				RightValue: nil,
			},
		}, diffItems)
	}

	{
		m1 := map[string]string{}
		m2 := map[string]string{"foo": "bar"}

		diffItems := DiffMap(m1, m2)
		assert.Equal(t, []DiffItem{
			{
				ChangeType: Added,
				Key:        "foo",
				LeftValue:  nil,
				RightValue: "bar",
			},
		}, diffItems)
	}
}
