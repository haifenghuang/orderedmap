package orderedmap

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"fmt"
)

var errEOA = errors.New("End of Array")

// An OrderedMap is a map where the keys keep the order that they're added.
// It is similar to indexed arrays. You can get the strings by key or by position.
// The OrderedMap is not safe for concurrent use.
type OrderedMap struct {
	// for preserving the order of key
	keys   []string
	values map[string]interface{}
}

// New create a new OrderedMap.
func New() *OrderedMap {
	return &OrderedMap{keys:[]string{}, values:make(map[string]interface{})}
}

// Get returns the value of the map based on its key.
// It will return nil if it doesn't exist.
func (om *OrderedMap) Get(key string) (interface{}, bool) {
	val, exists := om.values[key]
	return val, exists
}

// GetAt returns the value based on the provided pos.
func (om *OrderedMap) GetAt(pos int) (val interface{}, ok bool) {
	if om.values == nil {
		return nil, false
	}
	if pos >= len(om.keys) {
		val = nil
	} else {
		val, ok = om.values[om.keys[pos]]
		return val, ok
	}
	return
}

// Set sets the key/value of the map based on key and value.
// If the value already exists, the value will be replaced.
func (om *OrderedMap) Set(key string, value interface{}) {
	_, exists := om.values[key]
	if !exists { //not exists
		//add it to the keys array
		om.keys = append(om.keys, key)
	}
	om.values[key] = value
}

// SetAt sets the given key to the given value at the specified index.
func (om *OrderedMap) SetAt(index int, key string, val interface{}) {
	valLen := len(om.values)
	if index == -1 || index >= valLen {
		om.Set(key, val)
	}

	if om.values == nil {
		om.values = make(map[string]interface{})
	}

	if _, ok := om.values[key]; !ok { //if key not exists
		if index < -valLen {
			index = 0 // set at the begining
		}
		if index < 0 {
			index = valLen + index + 1
		}

		om.keys = append(om.keys, "") //assume the key is empty
		copy(om.keys[index+1:], om.keys[index:]) //shift the keys array
		om.keys[index] = key // reassign the key.
	}
	om.values[key] = val
}

// Delete remove an item from the map by the supplied key.
// If the key does not exist, nothing happens.
func (om *OrderedMap) Delete(key string) {
	_, ok := om.values[key]
	if !ok { // key not exists, do nothing.
		return
	}

	// remove from keys
	for i, k := range om.keys {
		if k == key {
			om.keys = append(om.keys[:i], om.keys[i+1:]...)
			break
		}
	}
	// remove from values
	delete(om.values, key)
}

// DeleteAt delete the key/value pair from the map by the supplied offset.
// If the offset is outside the range of the ordered map, nothing happens.
func (om *OrderedMap) DeleteAt(offset int) {
	if offset < 0 || offset >= len(om.keys) {
		return
	}
	om.Delete(om.keys[offset])
}

// Keys return the keys of the map in the order they were added.
func (om *OrderedMap) Keys() []string {
	return om.keys
}

// Values returns a slice of the values in the order they were added.
func (om *OrderedMap) Values() []interface{} {
	vals := make([]interface{}, len(om.keys))

	for idx, v := range om.keys {
		vals[idx] = om.values[v]
	}

	return vals
}

// Exists test whether the key exists or not.
func (om *OrderedMap) Exists(key string) (ok bool) {
	if om.values == nil {
		return false
	}

	_, ok = om.values[key]
	return
}

// Index returns the offset of the key in the ordered map.
// If the key could not be found, -1 is returned.
func (om *OrderedMap) Index(key string) int {
	for idx, k := range om.keys {
		if k == key {
			return idx
		}
	}
	return -1
}

func (om *OrderedMap) Len() int {
	return len(om.keys)
}

// String returns the JSON serialized string representation.
func (om *OrderedMap) String() string {
	json, _ := om.MarshalJSON()
	return string(json)
}

// MarshalJSON implements the json.Marshaller interface, so it could be serialized.
// When serializing, the keys of the map will keep the order they are added.
func (om OrderedMap) MarshalJSON() ([]byte, error) {
	var out bytes.Buffer

	out.WriteString("{")

	for idx, key:= range om.keys {
		if idx > 0 {
			out.WriteString(",")
		}

		esc := strings.Replace(key, `"`, `\"`, -1)
		out.WriteString(`"` + esc + `"`)

		out.WriteString(":")

		//marshal the value
		b, err := json.Marshal(om.values[key])
		if err != nil {
			return []byte{}, err
		}
		out.WriteString(string(b))
	} //end for

	out.WriteString("}")
	return out.Bytes(), nil
}

// UnmarshalJSON implements the json.Unmarshaller interface.
// so it could be use like below:
//      o := New()
//      err := json.Unmarshal([]byte(jsonString), &o)
func (om *OrderedMap) UnmarshalJSON(b []byte) error {
	//Using Decoder to parse the bytes.
	in := bytes.TrimSpace(b)
	dec := json.NewDecoder(bytes.NewReader(in))

	t, err := dec.Token()
	if err != nil {
		return err
	}

	// must open with a delim token '{'
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("expect JSON object open with '{'")
	}

	om.unmarshalJSON(dec)

	t, err = dec.Token() //'}'
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '}' {
		return fmt.Errorf("expect JSON object close with '}'")
	}

	return nil
}

func (om *OrderedMap) unmarshalJSON(dec *json.Decoder) error {
	for dec.More() { // Loop until it has no more tokens
		t, err := dec.Token()
		if err != nil {
			return err
		}

		key, ok := t.(string)
		if !ok {
			return fmt.Errorf("key must be a string, got %T\n", t)
		}

		val, err := parseObject(dec)
		if err != nil {
			return err
		}
		om.Set(key, val)
	}
	return nil
}

func parseObject(dec *json.Decoder) (interface{}, error) {
	t, err := dec.Token()
	if err != nil {
		return nil, err
	}

	switch tok := t.(type) {
	case json.Delim:
		switch tok {
		case '[': // If it's an array
			return parseArray(dec)
		case '{': // If it's a map
			om := New()
			err := om.unmarshalJSON(dec)
			if err != nil {
				return nil, err
			}
			_, err = dec.Token() // }
			if err != nil {
				return nil, err
			}
			return om, nil
		case ']':
			return nil, errEOA
		case '}':
			return nil, errors.New("unexpected '}'")
		default:
			return nil, fmt.Errorf("Unexpected delimiter: %q", tok)
		}
	default:
		return tok, nil
	}
}

func parseArray(dec *json.Decoder) ([]interface{}, error) {
	ret :=[]interface{}{}
	for {
		v, err := parseObject(dec)
		if err == errEOA {
			return ret, nil
		}
		if err != nil {
			return nil, err
		}
		ret = append(ret, v)
	}
}
