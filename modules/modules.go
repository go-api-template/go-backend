package modules

import (
	// Following modules are called implicitly
	// They must be imported at startup in order to initialize them
	_ "github.com/go-api-template/go-backend/modules/config"
	_ "github.com/go-api-template/go-backend/modules/logger"
)
