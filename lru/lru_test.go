package lru

import (
	"reflect"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("k1", String("v1"))

	if v, ok := lru.Get("k1"); !ok || string(v.(String)) != "v1" {
		t.Fatal("cache hit k1=v1 failed")
	}

	if _, ok := lru.Get("k2"); ok {
		t.Fatal("cache miss k2 failed")
	}
}

func TestRemoveOldest(t *testing.T) {
	k1, k2, k3 := "k1", "k2", "k3"
	v1, v2, v3 := "v1", "v2", "v3"
	max := len(k1 + k2 + v1 + v2)
	lru := New(int64(max), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, ok := lru.Get("k1"); ok || lru.Len() != 2 {
		t.Fatalf("TestRemoveOldest k1 failed")
	}

	if v3, ok := lru.Get("k3"); ok {
		t.Logf("k3: %v", v3)
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}

	lru := New(int64(10), callback)
	lru.Add("k1", String("123"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))

	expect := []string{"k1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}