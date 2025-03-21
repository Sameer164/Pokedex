package main

import (
	"Pokedex/internal/pokecache"
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := pokecache.NewCache(interval)
			cache.Set(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5
	const waitTime = baseTime + 5
	cache := pokecache.NewCache(baseTime)
	cache.Set("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime * time.Second)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}

func TestCleanInput(t *testing.T) {
	type Input struct {
		input    string
		expected []string
	}

	var tests []Input
	tests = append(tests, Input{input: "  Hello  world  ", expected: []string{"Hello", "world"}})

	for _, test := range tests {
		actual := cleanInput(test.input)
		if len(actual) != len(test.expected) {
			t.Errorf("\nLengths not equal. Expected: %v, Actual: %v", test.expected, actual)
			return
		}
		for i := range actual {
			word := actual[i]
			expectedWord := test.expected[i]
			if word != expectedWord {
				t.Errorf("\nExpected: %v\nActual: %v", expectedWord, word)
			}
		}
	}

}
