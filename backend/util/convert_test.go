package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Fullname string `json:"full_name"`
}

func TestConvertStructToMap(t *testing.T) {
	test := testStruct{
		Username: "test111",
		Password: "test123",
		Fullname: "George",
	}

	result := StructToMap(test)
	for fieldName, fieldValue := range result {
		switch fieldName {
		case "username":
			assert.Equal(t, "test111", fieldValue)
		case "password":
			assert.Equal(t, "test123", fieldValue)
		case "full_name":
			assert.Equal(t, "George", fieldValue)
		default:
			t.Errorf("Unexpected field name %s, and value %v", fieldName, fieldValue)

		}
	}
}
