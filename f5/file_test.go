package f5

import (
	"reflect"
	"testing"
)

func TestMerge(t *testing.T) {
	c1 := NewF5Config()
	c1.LtmNode["Foo"] = &LtmNode{Name: "Foo"}
	c1.LtmNode["Bar"] = &LtmNode{Name: "Bar"}
	c2 := NewF5Config()
	c2.LtmRule["Foo"] = &LtmRule{Name: "Foo"}
	c2.LtmRule["Bar"] = &LtmRule{Name: "Bar"}

	c3 := NewF5Config()
	c3.LtmNode["Foo"] = &LtmNode{Name: "Foo"}
	c3.LtmNode["Bar"] = &LtmNode{Name: "Bar"}
	c3.LtmRule["Foo"] = &LtmRule{Name: "Foo"}
	c3.LtmRule["Bar"] = &LtmRule{Name: "Bar"}

	err := c1.Merge(c2)
	if err != nil {
		t.Errorf("Error while merging F5 configuration: %s", err)
	}
	if !reflect.DeepEqual(c1, c3) {
		t.Errorf("Failed to merge F5 configuration")
	}
}
