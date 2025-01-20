package middlewares

import "github.com/vapusdata-oss/aistudio/aistudio/pkgs"

// Write all middleware init here

var (
	logger = pkgs.GetSubDMLogger("Middleware", "Base")
)
