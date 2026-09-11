package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LogLevel string

const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
	LevelFatal LogLevel = "FATAL"
)

type Log struct {
	ID          primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	Service     string                 `json:"service" bson:"service"`
	Environment string                 `json:"environment,omitempty" bson:"environment,omitempty"`
	Level       LogLevel               `json:"level" bson:"level"`
	Message     string                 `json:"message" bson:"message"`
	TraceID     string                 `json:"trace_id,omitempty" bson:"trace_id,omitempty"`
	SpanID      string                 `json:"span_id,omitempty" bson:"span_id,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Timestamp   time.Time              `json:"timestamp" bson:"timestamp"`
}
