package lru

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

func TestCache_Get(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("testKey1", String("testValue1"))
	want := "testValue1"
	if v, ok := lru.Get("testKey1"); !ok || string(v.(String)) != "testValue1" {
		t.Fatal(fmt.Sprintf("want %s but got %s", want, v))
	}
	if _, ok := lru.Get("testKey2"); ok {
		t.Fatal("cache miss testKey2 failed")
	}
}

func TestCache_RemoveOldest(t *testing.T) {
	keys := []string { "key1", "key2", "key3" }
	values := []string { "value1", "value2", "value3" }
	capacity := len(keys[0] + keys[1] + values[0] + values[1])
	lru := New(int64(capacity), nil)
	for i := 0; i < len(keys); i++ {
		lru.Add(keys[i], String(values[i]))
	}

	testKey := "key1"
	if _, ok := lru.Get(testKey); ok || lru.Len() != 2 {
		t.Fatal(fmt.Sprintf("Remove key: %s failed", testKey))
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value model.Value) {
		keys = append(keys, key)
	}
	lru := New(int64(10), callback)
	lru.Add("key1", String("value1"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))

	expect := []string{ "key1", "k2" }
	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
