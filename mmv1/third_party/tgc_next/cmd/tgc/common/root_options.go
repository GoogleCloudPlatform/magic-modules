package common

import "go.uber.org/zap"

type RootOptions struct {
	Verbosity            string
	ErrorLogger          *zap.Logger
	OutputLogger         *zap.Logger
	UseStructuredLogging bool
}
