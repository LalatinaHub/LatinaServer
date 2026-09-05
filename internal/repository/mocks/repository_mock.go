package mocks

import (
	"github.com/LalatinaHub/common/repository/mocks"
)

// Re-export mock repositories from github.com/LalatinaHub/common/repository/mocks.
type (
	MockUserRepository   = mocks.MockUserRepository
	MockServerRepository = mocks.MockServerRepository
	MockKVRepository     = mocks.MockKVRepository
	MockProxyRepository  = mocks.MockProxyRepository
)
