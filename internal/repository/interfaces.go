package repository

import (
	"github.com/LalatinaHub/common/repository"
)

// Re-export repository interfaces from github.com/LalatinaHub/common/repository.
type (
	UserRepository   = repository.UserRepository
	ServerRepository = repository.ServerRepository
	KVRepository     = repository.KVRepository
	ProxyRepository  = repository.ProxyRepository
)

// Forward repository constructors to github.com/LalatinaHub/common/repository.
var (
	NewUserRepository   = repository.NewUserRepository
	NewServerRepository = repository.NewServerRepository
	NewKVRepository     = repository.NewKVRepository
	NewProxyRepository  = repository.NewProxyRepository
)
