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
	layerB := []byte(`{"hello": 3, "yellow": 56, "world": {"another": 32, "yetanother": 33}}`)
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
	expected := `{"hello":3,"wilbur":"razzle","world":{"another":32,"something":2,"yetanother":33},"yellow":56}`
	strData := string(data)
	if strData != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s\n", expected, strData)
		return
	}
}

func TestGetString(t *testing.T) {
	source := []byte(`{"hello": 1, "world": {"something": "doggy"}}`)
	kv, err := NewFromJson(source)
	if err != nil {
		t.Error(err)
		return
	}
	v, err := kv.String("world", "something")
	if err != nil {
		t.Error(err)
		return
	}

	if v != "doggy" {
		t.Errorf("Expected doggy, got %s", v)
		return
	}
}

func TestSetValue(t *testing.T) {
	source := []byte(`{"hello": 1, "world": {"something": "doggy"}}`)
	kv, err := NewFromJson(source)
	if err != nil {
		t.Error(err)
		return
	}
	err = kv.SetValue(33, "world", "something")
	if err != nil {
		t.Error(err)
		return
	}

	v, err := kv.Number("world", "something")
	if v != 33 {
		t.Errorf("Expected 33, got %v", v)
		return
	}
}

func TestCreateValue(t *testing.T) {
	source := []byte(`{"hello": 1, "world": {"something": "doggy"}}`)
	kv, err := NewFromJson(source)
	if err != nil {
		t.Error(err)
		return
	}
	err = kv.CreateValue(33, "world", "something", "snoop")
	if err != nil {
		t.Error(err)
		return
	}

	v, err := kv.Number("world", "something", "snoop")
	if v != 33 {
		t.Errorf("Expected 33, got %v", v)
		return
	}
}

func TestEmptyKeyval(t *testing.T) {
	var source []byte
	_, err := NewFromJson(source)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSplitKey(t *testing.T) {
	keys := SplitKey("hello.world.now")
	if len(keys) != 3 {
		t.Errorf("Expected 3 key components")
	}
}

func TestSplitKeySpecialDelim(t *testing.T) {
	keys := SplitKey("hello:world:now", ":")
	if len(keys) != 3 {
		t.Errorf("Expected 3 key components")
	}
}

func TestGetKv(t *testing.T) {
	data := []byte(`{"users": {"maxine": {"age": 38, "name": "maxine"}, "grigori": {"age": 22, "name": "grigori"}}}`)
	kv, err := NewFromJson(data)
	if err != nil {
		t.Error(err)
		return
	}

	mkv, err := kv.GetKeyVal("users", "maxine")
	if err != nil {
		t.Error(err)
		return
	}

	age, err := mkv.Value("age")
	if err != nil {
		t.Error(err)
		return
	}

	if age.(float64) != 38 {
		t.Error(err)
		return
	}
}

func TestGetValueBadKey(t *testing.T) {
	kv := New()
	kv.CreateValue("new york city", "city", "ny", "name")
	kv.CreateValue(32423534, "city", "ny", "size")
	_, err := kv.Value("city.ny.nam")
	if err == nil {
		t.Error(err)
		return
	}
}
