package keyval

import "testing"

func TestDeepCopy(t *testing.T) {
	source := []byte(`{"hello": 1, "world": {"something": 2}}`)
	kv, err := NewFromJson(source)
	if err != nil {
		t.Error(err)
	}
	copy := kv.Copy()
	copy
}
