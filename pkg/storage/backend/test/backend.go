package test

import (
	"github.com/lerenn/chonkfs/pkg/storage/backend"
	"github.com/stretchr/testify/suite"
)

type BackEndSuite struct {
	BackEnd backend.BackEnd
	suite.Suite
}
