package svc

import (
	"sshx/internal/types"
)

type ServiceContext struct {
	Config types.Config
}

func NewServiceContext(c types.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
