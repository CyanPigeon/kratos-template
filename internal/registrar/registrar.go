package registrar

import (
	"github.com/CyanPigeon/toktik/middleware/discovery"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	"github.com/hashicorp/consul/api"
)

var ProviderSet = wire.NewSet(NewRegistry)

func NewRegistry() (registry.Registrar, error) {
	return discovery.New(api.DefaultConfig())
}
