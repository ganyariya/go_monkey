package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello"}
	hello2 := &String{Value: "Hello"}
	other1 := &String{Value: "Hello, Other"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Fatalf("same object does not have same hash key")
	}
	if hello1.HashKey() == other1.HashKey() {
		t.Fatalf("different object has same hash key ")
	}
}
