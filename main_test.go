package main

import "testing"

func TestGenerate(t *testing.T) {

	cfg := Config{
		file:       "./testdata/user.go",
		structName: "UserLink",
	}

	err := cfg.GenerateGormCode()
	if err != nil {
		t.Fatal(err)
	}

}
