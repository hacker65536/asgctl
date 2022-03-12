package myaws

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDescribeTags(t *testing.T) {
	t.Setenv("AWS_PROFILE", "production")
	describeTags()
}

func TestMakeFilters(t *testing.T) {
	list := make(map[string]string)
	list["project"] = "ddd"
	list["stage"] = "staging"

	f := makeFilters(list)
	j, _ := json.Marshal(f)
	fmt.Println(string(j))
}
