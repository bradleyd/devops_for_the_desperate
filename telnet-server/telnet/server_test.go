package telnet

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
)

type tcpServerMock struct {
	mock.Mock
}

func (t *tcpServerMock) Run() {
	fmt.Println("Mocked charge notification function")
	t.Called()
}

func TestServerRun(t *testing.T) {
	tcpServer := new(tcpServerMock)
	tcpServer.On("Run").Once()
	tcpServer.Run()
	tcpServer.AssertExpectations(t)
}
