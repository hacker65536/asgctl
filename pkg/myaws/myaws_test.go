package myaws

import (
	"encoding/json"
	"reflect"
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
)

func TestDescribeTags(t *testing.T) {
	t.Setenv("AWS_PROFILE", "production")
	f := []types.Filter{
		{
			Name:   aws.String("Key"),
			Values: []string{"Stage"},
		},
		{
			Name:   aws.String("Value"),
			Values: []string{"ojiro"},
		},
	}
	describeTags(f)
}

func TestMakeFilters(t *testing.T) {
	list := make(map[string]string)
	list["project"] = "ddd"
	list["stage"] = "staging"

	f := makeFilters(list)
	j, _ := json.Marshal(f)
	t.Log(string(j))
}

func TestGetAsgs(t *testing.T) {

	list := make(map[string]string)
	list["project"] = "ddd"
	list["stage"] = "staging"

	f := makeFilters(list)
	getAsgs(f)
}

func TestLeftJoin(t *testing.T) {

	a := []string{"same2", "baa", "same", "1", "ddd"}
	b := []string{"a", "3", "ba", "same", "same2"}
	c := []string{"baa", "c", "same", "same2", "e"}

	got := leftJoin(a, b, c)
	sort.SliceStable(got, func(i, j int) bool {
		return got[i] < got[j]
	})

	want := []string{"same2", "same"}
	sort.SliceStable(want, func(i, j int) bool {
		return want[i] < want[j]
	})

	if reflect.DeepEqual(got, want) {
		t.Logf("\ngot = %v\nwant = %v", got, want)
	} else {
		t.Errorf("not equal got = %v, want = %v", got, want)

	}

}
