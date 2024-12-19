package mock

type MockLogger struct{}

func NewLogger() *MockLogger {
	return &MockLogger{}
}

func (m *MockLogger) Log(v ...interface{}) {}
