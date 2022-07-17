package main

import (
	"context"
	"reflect"
	"regexp"
	"testing"

	container "github.com/golobby/container/v3"
	"github.com/noobj/swim-crowd-lambda-go/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
)

type MockEntryModel struct {
	insertOneIsCalled bool
	insertOneArgument interface{}
}

func (m *MockEntryModel) InsertOne(doc bson.D) {
	m.insertOneIsCalled = true
	m.insertOneArgument = doc
}

func (m *MockEntryModel) Aggregate(stages []bson.D) []any {
	mockData := []any{}
	return mockData
}

func (m *MockEntryModel) Disconnect() func() {
	return func() {}
}

type ArgumentType struct {
	Amount int
	Time   string
}

func TestHandler(t *testing.T) {
	mockEntryModel := MockEntryModel{}

	container.Singleton(func() repositories.Repository {
		return &mockEntryModel
	})
	if _, err := Handler(context.TODO()); err != nil {
		t.Errorf("error %s", err)
	}

	if mockEntryModel.insertOneIsCalled != true {
		t.Error("InsertOne has not been called ")
	}

	doc, _ := bson.Marshal(mockEntryModel.insertOneArgument)
	var argument ArgumentType
	err := bson.Unmarshal(doc, &argument)

	dateTimeReg := regexp.MustCompile(`^20\d\d-[0-1][0-9]-[0-3][0-9] \d{2}:\d{2}`)

	if reflect.TypeOf(argument.Amount).String() != "int" {
		t.Error("wrong argument Amount Type")
	}

	if !dateTimeReg.Match([]byte(argument.Time)) {
		t.Error("wrong argument Time format")
	}

	if err != nil {
		t.Errorf("error %s", err)
	}
}
