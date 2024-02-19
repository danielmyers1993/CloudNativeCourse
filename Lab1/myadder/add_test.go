package myadder

import "testing"

func TestAdd(t *testing.T) {
    want := 7
    got := Add(3, 4)
    if want != got {
        t.Errorf("Error in myadder.Add; Want 7, Got %d", got)
    }
}