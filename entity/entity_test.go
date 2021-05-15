package entity

import (
	"testing"

	"github.com/bilginyuksel/gorify"
)

type TestI struct {
	Username string `between:"5-15"`
	Password string `default:"123456"`
	Age      int    `between:"20-80"`
}

func TestValidatation(t *testing.T) {
	testInt := &TestI{Username: "bilgili", Age: 25}
	err := gorify.Validate(testInt)
	if err != nil {
		t.Errorf("err: %v\n", err)
	}
	t.Logf("final testI: %v\n", testInt)
}
