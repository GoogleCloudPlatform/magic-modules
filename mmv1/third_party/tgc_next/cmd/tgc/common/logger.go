// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"bytes"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type errorDetails struct {
	// error message
	error string
	// stacktrace or additional context
	context string
}

func (ed errorDetails) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("error", ed.error)
	enc.AddString("context", ed.context)
	return nil
}

type errorEncoder struct {
	zapcore.Encoder
}

func (enc errorEncoder) Clone() zapcore.Encoder {
	return errorEncoder{
		Encoder: enc.Encoder.Clone(),
	}
}

func (enc errorEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	ed := errorDetails{
		error:   ent.Message,
		context: ent.Stack,
	}
	fields = append([]zapcore.Field{
		zap.Object("error_details", ed),
	}, fields...)
	return enc.Encoder.EncodeEntry(ent, fields)
}

func newJSONEncoder(cfg zapcore.EncoderConfig) errorEncoder {
	return errorEncoder{
		Encoder: zapcore.NewJSONEncoder(cfg),
	}
}

func newConsoleEncoder(cfg zapcore.EncoderConfig) errorEncoder {
	return errorEncoder{
		Encoder: zapcore.NewConsoleEncoder(cfg),
	}
}

func NewErrorLogger(verbosity string, useStructuredLogging bool, writeSyncer zapcore.WriteSyncer) *zap.Logger {
	// Return a logger that produces expected structured output format for errors
	var level zapcore.LevelEnabler
	options := []zap.Option{
		zap.Fields(
			// Message format version
			zap.String("version", "v1.0.0"),
		),
	}

	switch verbosity {
	case "debug":
		level = zap.DebugLevel
		options = append(options, zap.AddStacktrace(zap.WarnLevel))
	// case "info" is handled by default
	case "warning":
		level = zap.WarnLevel
		options = append(options, zap.AddStacktrace(zap.ErrorLevel))
	case "error":
		level = zap.ErrorLevel
		options = append(options, zap.AddStacktrace(zap.ErrorLevel))
	case "critical":
		level = zap.PanicLevel
		options = append(options, zap.AddStacktrace(zap.PanicLevel))
	case "none":
		return zap.NewNop()
	default:
		level = zap.InfoLevel
		options = append(options, zap.AddStacktrace(zap.WarnLevel))
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "level",
		MessageKey:    "",
		StacktraceKey: "",
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime:    zapcore.RFC3339NanoTimeEncoder,
	}
	var encoder errorEncoder
	if useStructuredLogging {
		encoder = newJSONEncoder(encoderConfig)
	} else {
		encoder = newConsoleEncoder(encoderConfig)
	}
	core := zapcore.NewCore(encoder, writeSyncer, level)
	return zap.New(core, options...)
}

func NewOutputLogger(writeSyncer zapcore.WriteSyncer) *zap.Logger {
	// Return a logger that produces expected structured output format for output
	options := []zap.Option{
		zap.Fields(
			// Message format version
			zap.String("version", "v1.0.0"),
		),
	}

	level := zap.NewAtomicLevelAt(zap.DebugLevel)
	options = append(options, zap.AddStacktrace(zap.WarnLevel))

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "",
		MessageKey:    "body",
		StacktraceKey: "",
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime:    zapcore.RFC3339NanoTimeEncoder,
	}
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, writeSyncer, level)
	return zap.New(core, options...)
}

type bufferWriteSyncer struct {
	*bytes.Buffer
}

func (bws bufferWriteSyncer) Sync() error {
	return nil
}

func NewTestErrorLogger(verbosity string, useStructuredLogging bool) (*zap.Logger, *bytes.Buffer) {
	buf := new(bytes.Buffer)
	syncer := bufferWriteSyncer{Buffer: buf}
	logger := NewErrorLogger(verbosity, useStructuredLogging, syncer)
	return logger, syncer.Buffer
}

func NewTestOutputLogger() (*zap.Logger, *bytes.Buffer) {
	buf := new(bytes.Buffer)
	syncer := bufferWriteSyncer{Buffer: buf}
	logger := NewOutputLogger(syncer)
	return logger, syncer.Buffer
}
