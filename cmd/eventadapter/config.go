package main

import (
	"fmt"
	"regexp"

	cloudevents "github.com/cloudevents/sdk-go/v2"

	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
)

type Filter struct {
	Key      string   `yaml:"key"`
	Value    string   `yaml:"value"`
	Negative bool     `yaml:"negative"`
	ValueIn  []string `yaml:"valueIn"`
}

type FilterExp struct {
	Filter   `yaml:",inline"`
	ValueExp *regexp.Regexp `yaml:"-"`
}

func (t *FilterExp) UnmarshalYAML(unmarshal func(interface{}) error) error {
	err := unmarshal(&t.Filter)
	if err != nil {
		return err
	}

	if t.Value == "" && len(t.ValueIn) == 0 {
		return fmt.Errorf("value and valueIn can't be empty. ")
	}

	if t.Value != "" {
		t.ValueExp, err = regexp.Compile(t.Value)
		if err != nil {
			return fmt.Errorf("get invalid value %s. ", t.Value)
		}
	}
	return nil
}

func (t FilterExp) IsMatch(val string) bool {
	if t.ValueExp != nil {
		ok := t.ValueExp.MatchString(val)
		if t.Negative {
			return !ok
		}
		return ok
	}
	for _, s := range t.ValueIn {
		if s == val {
			return true
		}
	}
	return false
}

type AdapterConfig struct {
	Filters       []FilterExp    `yaml:"filters"`
	DestListeners []ListenerInfo `yaml:"destListeners"`
}

type ListenerInfo struct {
	EventListenerName string             `yaml:"eventListenerName"`
	EventListenerNs   string             `yaml:"eventListenerNamespace"`
	Bindings          []triggersv1.Param `yaml:"params"`
	ReqHeadFields     map[string]string  `yaml:"reqHeadFields"`
	ReqBodyTemplate   string             `yaml:"reqBodyTemplate"`
}

const ConfigKey = "adapters"

var CurrentConfig = make([]AdapterConfig, 0)

// 用于jsonPath 获取值
type DataEvent struct {
	Context    cloudevents.EventContext `json:"context"`
	Data       interface{}              `json:"data"`
	Extensions interface{}              `json:"extensions"`
}

func GetListenerKey(name, namespace string) string {
	return fmt.Sprintf("%s/%s", name, namespace)
}
