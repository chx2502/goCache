package lfu

import (
	"fmt"
	"gocache/model"
	"reflect"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestLfu_Get(t *testing.T) {
	lfu := New(int64(0), nil)
	lfu.Add("testKey1", String("testValue1"))
	want := "testValue1"
	if v, ok := lfu.Get("testKey1"); !ok || string(v.(String)) != "testValue1" {
		t.Fatal(fmt.Sprintf("want %s but got %s", want, v))
	}
	if _, ok := lfu.Get("testKey2"); ok {
		t.Fatal("cache miss testKey2 failed")
	}
}

func TestLfu_Remove(t *testing.T) {
	keys := []string { "key1", "key2", "key3" }
	values := []string { "value1", "value2", "value3" }
	capacity := len(keys[0] + keys[1] + values[0] + values[1])
	lfu := New(int64(capacity), nil)
	for i := 0; i < len(keys); i++ {
		lfu.Add(keys[i], String(values[i]))
	}

	testKey := keys[0]
	if _, ok := lfu.Get(testKey); ok || lfu.Len() != 2 {
		t.Fatal(fmt.Sprintf("Remove key: %s failed", testKey))
	}
	lfu.Add(keys[0], String(values[0]))
	if v, ok := lfu.Get(keys[0]); !ok || string(v.(String)) != values[0] {
		want := values[0]
		t.Fatal(fmt.Sprintf("want %s but got %s", want, v))
	}
	lfu.Add(keys[1], String(values[1]))
	if _, ok := lfu.Get(keys[2]); ok || lfu.Len() != 2 {
		t.Fatal(fmt.Sprintf("Remove key: %s failed", keys[2]))
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value model.Value) {
		keys = append(keys, key)
	}
	lfu := New(int64(10), callback)
	lfu.Add("key1", String("value1"))
	lfu.Add("k2", String("k2"))
	lfu.Add("k3", String("k3"))
	lfu.Add("k4", String("k4"))

	expect := []string{ "key1", "k2" }
	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
