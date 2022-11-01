// Copyright (c) 2021 Jing-Ying Chen. Subject to the MIT License.

package goutil

import (
	"encoding/json"
	"reflect"
)

func JsonEqual(s1, s2 []byte) bool {
	var o1, o2 interface{}
	return json.Unmarshal(s1, &o1) == nil &&
		json.Unmarshal(s2, &o2) == nil &&
		reflect.DeepEqual(o1, o2)
}

func JsonContains(s1, s2 []byte) bool {
	var o1, o2 interface{}
	return json.Unmarshal(s1, &o1) == nil &&
		json.Unmarshal(s2, &o2) == nil &&
		ContainsRecursive(o1, o2)
}

func JsonStrEqual(s1, s2 string) bool {
	return JsonEqual([]byte(s1), []byte(s2))
}

func JsonStrContains(s1, s2 string) bool {
	return JsonContains([]byte(s1), []byte(s2))
}

func JsonMarshalContains(i1, i2 interface{}) bool {
	buf1, err := json.Marshal(i1)
	if err != nil {
		return false
	}
	buf2, err := json.Marshal(i2)
	if err != nil {
		return false
	}
	return JsonContains(buf1, buf2)
}

// REF: https://google.github.io/styleguide/jsoncstyleguide.xml

type JsonMsg struct {
	ApiVersion string `json:"apiVersion,omitEmpty"`
	Data       `json:"data,omitempty"`
	Error      `json:"error,omitempty"`
}
type Data struct {
	Kind    string            `json:"kind,omitempty"`
	Payload string            `json:"payload,omitempty"` // make it simple - not using Fields
	Values  map[string]string `json:"values,omitempty"`  // make it simple - not using Fields
	Items   []json.RawMessage `json:"items,omitempty"`
}

func (d *Data) AddItem(item interface{}) error {
	js, err := json.Marshal(item)
	if err != nil {
		return err
	}
	d.Items = append(d.Items, js)
	return nil
}

func SimpleJsonError(msg string, codes ...int) []byte {
	err := JsonMsg{
		Error: Error{
			Message: msg,
		},
	}
	if len(codes) > 0 {
		err.Error.Code = codes[0]
	}
	ret, _ := json.MarshalIndent(err, "", "  ")
	return ret
}

func SimpleJsonData(values map[string]string, kinds ...string) []byte {
	msg := JsonMsg{
		Data: Data{},
	}
	if len(values) > 0 {
		msg.Data.Values = values
	}
	if len(kinds) > 0 {
		msg.Data.Kind = kinds[0]
	}
	ret, _ := json.MarshalIndent(msg, "", "  ")
	return ret
}
