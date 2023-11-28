package helpers

import (
	"errors"
	"fmt"
	"testing"
)

// go test -run TestErrPanicIfErr
func TestErrPanicIfErr(t *testing.T) {
	for i := 0; i < 99999; i++ {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()
		PanicIfErr(errors.New("test error"))
		t.Errorf("The code did not panic")
	}
}

// go test -run TestNilPanicIfErr
func TestNilPanicIfErr(t *testing.T) {
	for i := 0; i < 99999; i++ {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Recovered in f", r)
			}
		}()
		PanicIfErr(nil)
	}
}
