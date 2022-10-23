package keyval

import (
	"encoding/json"
	"testing"
)

func TestDeepCopy(t *testing.T) {
	source := []byte(`{"hello": 1, "world": {"something": 2}}`)
	kv, err := NewFromJson(source)
	if err != nil {
		t.Error(err)
		return
	}
	copy := kv.Copy()
	err = copy.SetValue(3, "world", "something")
	if err != nil {
		t.Errorf("Unable to set value: %v", err)
		return
	}

	num, err := kv.Number("world", "something")
	if err != nil {
		t.Errorf("Unable to retrieve value: %v", err)
		return
	}
	if num != 2 {
		t.Errorf("Expected 2, got %v", num)
		return
	}
}

func TestStacking(t *testing.T) {
	layerA := []byte(`{"hello": 1, "world": {"something": 2}, "wilbur": "razzle"}`)
	layerB := []byte(`{"hello": 3, "yellow": 56, "world": {"another": 32}}`)
	kvA, err := NewFromJson(layerA)
	if err != nil {
		t.Error(err)
		return
	}
	kvB, err := NewFromJson(layerB)
	if err != nil {
		t.Error(err)
		return
	}

	final := kvA.Stack(kvB)
	data, err := json.Marshal(final.root)
	if err != nil {
		t.Error(err)
		return
	}
	strData := string(data)
	if strData != `{"hello":3,"wilbur":"razzle","world":{"another":32,"something":2},"yellow":56}` {
		t.Error("Got unexpected value")
		return
	}
}
