package testutil

import "runtime"

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

func FnName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
