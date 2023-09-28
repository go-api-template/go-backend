package logger

import (
	// Following modules are called implicitly
	// They must be imported at startup in order to initialize them
	_ "github.com/go-api-template/go-backend/modules/logger/access"
	_ "github.com/go-api-template/go-backend/modules/logger/database"
	_ "github.com/go-api-template/go-backend/modules/logger/main"
)
