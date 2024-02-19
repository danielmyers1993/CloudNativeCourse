package cache

import "testing"

/*
func TestReadWrite(t *testing.T) {
	testlru := NewCacher[string, string](3)
	testkey := "key"
	testval := "val"
	if err := testlru.Put(testkey, testval); err != nil {
		t.Error("Write test failed", err)
	}
	if val, err := testlru.Get(testkey); err != nil {
		t.Error("Read test failed", err)
	} else if val != "val" {
		t.Error("Read test failed, incorrect value")
	}
}
*/

func TestWriteWithEviction(t *testing.T) {
	testlru := NewCacher[string, string](3)
	testlru.Put("key1", "val1")
	testlru.Put("key2", "val2")
	testlru.Put("key3", "val3")
	testlru.Put("key4", "val4") // Triggers eviction of "key1:val1"
	if val, err := testlru.Get("key1"); err == nil {
		t.Error("LRU replacement policy test failed, got val", val)
	}
}
