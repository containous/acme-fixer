package check

import (
	"fmt"
	"testing"
)

func TestIsSafe(t *testing.T) {
	// acura.roimotors.com acuraofoakville.com
	// "*.letsencrypt.org"
	// automotriztauro.com automotriztauro.com.mx
	// *.miusina.net *.tupesca.net miusina.net tupesca.net
	// hostalsanjorgecafayate.com.ar

	safe, err := IsSafe("*.miusina.net", true)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(safe)
}
