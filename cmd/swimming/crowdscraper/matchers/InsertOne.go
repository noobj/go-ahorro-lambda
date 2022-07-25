package matchers

import (
	"fmt"
	"reflect"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
)

type ArgumentType struct {
	Amount int
	Time   string
}

type InsertOneMatcher struct {
	pattern *regexp.Regexp
}

func Regexp(pattern string) *InsertOneMatcher {
	return &InsertOneMatcher{
		pattern: regexp.MustCompile(pattern),
	}
}

func (m *InsertOneMatcher) String() string {
	return fmt.Sprintf("matches pattern /%v/", m.pattern)
}

func (m *InsertOneMatcher) Matches(x interface{}) bool {
	doc, _ := bson.Marshal(x)
	var argument ArgumentType
	err := bson.Unmarshal(doc, &argument)

	if err != nil {
		return false
	}

	if reflect.TypeOf(argument.Amount).String() != "int" {
		return false
	}

	return m.pattern.MatchString(argument.Time)
}
