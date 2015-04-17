package shopline

import (
	"testing"
)

func Test0(t *testing.T) {
	t.Log("0\n")
	r := algoritmo("Now in the street, there is violence And-and a lots of work to be done No place to hang out our washing And-and I can't blame all on the sun", "1234567890123456")
	t.Log(r + "\n")
}

func Test1(t *testing.T) {
	r0 := rr(5, 7)
	r1 := rr(5, 7)
	r2 := rr(5, 7)
	if r0 < 5 || r0 > 6 {
		t.Fatalf("r0 is %v\n", r0)
	}
	if r1 < 5 || r1 > 6 {
		t.Fatalf("r1 is %v\n", r1)
	}
	if r2 < 5 || r2 > 6 {
		t.Fatalf("r2 is %v\n", r2)
	}
}
