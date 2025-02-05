package mem

import (
	"testing"

	"github.com/lerenn/chonkfs/pkg/storage/backend/test"
	"github.com/stretchr/testify/suite"
)

func TestBackEndSuite(t *testing.T) {
	suite.Run(t, new(BackEndSuite))
}

type BackEndSuite struct {
	test.BackEndSuite
}

func (suite *BackEndSuite) SetupTest() {
	suite.BackEnd = NewBackEnd()
}
