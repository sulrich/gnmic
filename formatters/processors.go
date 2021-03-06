package formatters

import (
	"log"

	"github.com/mitchellh/mapstructure"
)

var EventProcessors = map[string]Initializer{}

var EventProcessorTypes = []string{
	"event-add-tag",
	"event-convert",
	"event-date-string",
	"event-delete",
	"event-drop",
	"event-override-ts",
	"event-strings",
	"event-to-tag",
	"event-write",
}

type Initializer func() EventProcessor

func Register(name string, initFn Initializer) {
	EventProcessors[name] = initFn
}

type EventProcessor interface {
	Init(interface{}, *log.Logger) error
	Apply(*EventMsg)
}

func DecodeConfig(src, dst interface{}) error {
	decoder, err := mapstructure.NewDecoder(
		&mapstructure.DecoderConfig{
			DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
			Result:     dst,
		},
	)
	if err != nil {
		return err
	}
	return decoder.Decode(src)
}
