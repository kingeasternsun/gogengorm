package main

import "testing"

func TestGenerateUniqueIndex(t *testing.T) {

	cfg := Config{
		file:         "./testdata/user.go",
		structName:   "UserLink",
		templateFile: "./gorm.template",
	}

	err := cfg.GenerateGormCode()
	if err != nil {
		t.Fatal(err)
	}

}

func TestGenerateIndex(t *testing.T) {

	cfg := Config{
		file:       "./testdata/user.go",
		structName: "UserSell",
		// all:          true,
		templateFile: "./gorm.template",
	}

	err := cfg.GenerateGormCode()
	if err != nil {
		t.Fatal(err)
	}

}
