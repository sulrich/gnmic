package event_to_tag

import (
	"reflect"
	"testing"

	"github.com/karimra/gnmic/formatters"
)

type item struct {
	input  *formatters.EventMsg
	output *formatters.EventMsg
}

var testset = map[string]struct {
	processorType string
	processor     map[string]interface{}
	tests         []item
}{
	"1_value_match": {
		processorType: processorType,
		processor: map[string]interface{}{
			"value-names": []string{".*name$"},
		},
		tests: []item{
			{
				input:  nil,
				output: nil,
			},
			{
				input: &formatters.EventMsg{
					Values: map[string]interface{}{}},
				output: &formatters.EventMsg{
					Values: map[string]interface{}{}},
			},
			{
				input: &formatters.EventMsg{
					Values: map[string]interface{}{"name": "dummy"}},
				output: &formatters.EventMsg{
					Tags:   map[string]string{"name": "dummy"},
					Values: map[string]interface{}{}},
			},
		},
	},
	"1_value_match_with_keep": {
		processorType: processorType,
		processor: map[string]interface{}{
			"value-names": []string{".*name$"},
			"keep":        true,
		},
		tests: []item{
			{
				input:  nil,
				output: nil,
			},
			{
				input: &formatters.EventMsg{
					Values: map[string]interface{}{}},
				output: &formatters.EventMsg{
					Values: map[string]interface{}{}},
			},
			{
				input: &formatters.EventMsg{
					Values: map[string]interface{}{"name": "dummy"}},
				output: &formatters.EventMsg{
					Tags:   map[string]string{"name": "dummy"},
					Values: map[string]interface{}{"name": "dummy"}},
			},
		},
	},
	"2_value_match": {
		processorType: processorType,
		processor: map[string]interface{}{
			"value-names": []string{".*name$"},
		},
		tests: []item{
			{
				input:  nil,
				output: nil,
			},
			{
				input: &formatters.EventMsg{
					Values: map[string]interface{}{}},
				output: &formatters.EventMsg{
					Values: map[string]interface{}{}},
			},
			{
				input: &formatters.EventMsg{
					Values: map[string]interface{}{
						"name":        "dummy",
						"second_name": "dummy2"},
				},
				output: &formatters.EventMsg{
					Tags: map[string]string{
						"name":        "dummy",
						"second_name": "dummy2"},
					Values: map[string]interface{}{}},
			},
		},
	},
	"2_value_match_with_keep": {
		processorType: processorType,
		processor: map[string]interface{}{
			"value-names": []string{".*name$"},
			"keep":        true,
		},
		tests: []item{
			{
				input:  nil,
				output: nil,
			},
			{
				input: &formatters.EventMsg{
					Values: map[string]interface{}{}},
				output: &formatters.EventMsg{
					Values: map[string]interface{}{}},
			},
			{
				input: &formatters.EventMsg{
					Values: map[string]interface{}{
						"name":        "dummy",
						"second_name": "dummy2"},
				},
				output: &formatters.EventMsg{
					Tags: map[string]string{
						"name":        "dummy",
						"second_name": "dummy2"},
					Values: map[string]interface{}{
						"name":        "dummy",
						"second_name": "dummy2"}},
			},
		},
	},
}

func TestEventToTag(t *testing.T) {
	for name, ts := range testset {
		if pi, ok := formatters.EventProcessors[ts.processorType]; ok {
			t.Log("found processor")
			p := pi()
			err := p.Init(ts.processor, nil)
			if err != nil {
				t.Errorf("failed to initialize processors: %v", err)
				return
			}
			t.Logf("processor: %+v", p)
			for i, item := range ts.tests {
				t.Run("uint_convert", func(t *testing.T) {
					t.Logf("running test item %d", i)
					var inputMsg *formatters.EventMsg
					if item.input != nil {
						inputMsg = &formatters.EventMsg{
							Name:      item.input.Name,
							Timestamp: item.input.Timestamp,
							Tags:      make(map[string]string),
							Values:    make(map[string]interface{}),
							Deletes:   item.input.Deletes,
						}
						for k, v := range item.input.Tags {
							inputMsg.Tags[k] = v
						}
						for k, v := range item.input.Values {
							inputMsg.Values[k] = v
						}
					}
					p.Apply(item.input)
					t.Logf("input: %+v, changed: %+v", inputMsg, item.input)
					if !reflect.DeepEqual(item.input, item.output) {
						t.Errorf("failed at %s item %d, expected %+v, got: %+v", name, i, item.output, item.input)
					}
				})
			}
		} else {
			t.Errorf("event processor %s not found", ts.processorType)
		}
	}
}
