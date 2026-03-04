package testutil

import (
	"crypto/rand"

	"github.com/HorizenOfficial/vela/pkg/common"
)

type MockFunctions struct {
	MockedFunctions map[string]interface{}
}

func (c *MockFunctions) AddMockedFunc(key string, mockedFunc interface{}) {
	c.MockedFunctions[key] = mockedFunc
}

func (c *MockFunctions) RemoveMockedFunc(key string) {
	delete(c.MockedFunctions, key)
}

func (c *MockFunctions) GetMockedFunc(key string) (interface{}, bool) {
	f, ok := c.MockedFunctions[key]
	return f, ok
}

func NewMockFunctions() *MockFunctions {
	return &MockFunctions{
		MockedFunctions: make(map[string]interface{}),
	}
}

// GenerateRandomRequestID generates a random ID
func GenerateRandomRequestID() common.RequestIdType {
	b := make([]byte, 32)
	_, _ = rand.Read(b)

	return common.RequestIdType(b)
}
