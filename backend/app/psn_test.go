package app

import "testing"

func TestSearch(t *testing.T) {
	searched, err := Search("Человек паук")
	if err != nil {
		panic(err)
	}

	for _, s := range searched {
		t.Log(s)
	}
}
