package passwords

import (
	"testing"
	"fmt"
)

func TestCompute(t * testing.T) {
	password := ComputeWithIteration("garrydanger", 3)

	a := ComputeWithSalt("garrydanger", 3, password.Salt)
	if a.Hash != password.Hash {
		t.Fatal("nonderministic hash detected")
	}

  fmt.Println(a)
}
