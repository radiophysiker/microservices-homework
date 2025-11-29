package logger

import (
	"context"
	"fmt"
	"math"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
	"go.uber.org/zap/zapcore"
)

// newOTLPCore создает zapcore.Core для отправки логов в OTEL Collector через OTLP
// и возвращает Core вместе с функцией корректного завершения провайдера.
func newOTLPCore(ctx context.Context, cfg Config) (zapcore.Core, func(context.Context) error, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName()),
		),
		resource.WithHost(),
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		return nil, nil, err
	}

	exporter, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(cfg.OTELCollectorEndpoint()),
		otlploggrpc.WithInsecure(),
	)
	if err != nil {
		return nil, nil, err
	}

	processor := sdklog.NewBatchProcessor(exporter)

	provider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(processor),
	)

	otlpLogger := provider.Logger(
		cfg.ServiceName(),
		log.WithInstrumentationVersion("1.0.0"),
	)

	writeSyncer := &otlpWriteSyncer{
		logger: otlpLogger,
	}

	core := &otlpCore{
		encoder: zapcore.NewJSONEncoder(buildProductionEncoderConfig()),
		writer:  writeSyncer,
		level:   dynamicLevel,
	}

	shutdown := func(ctx context.Context) error {
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		return provider.Shutdown(shutdownCtx)
	}

	return core, shutdown, nil
}

// otlpWriteSyncer реализует zapcore.WriteSyncer для отправки логов в OTLP
type otlpWriteSyncer struct {
	logger log.Logger
}

// Write реализует io.Writer, но не используется напрямую
// Вместо этого используется метод Emit() через кастомный Core
func (w *otlpWriteSyncer) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// Sync сбрасывает буферы
func (w *otlpWriteSyncer) Sync() error {
	return nil
}

// otlpCore реализует zapcore.Core для конвертации zap записей в OTLP
type otlpCore struct {
	encoder zapcore.Encoder
	writer  *otlpWriteSyncer
	level   zapcore.LevelEnabler
	fields  []zapcore.Field
}

// Enabled проверяет, должен ли уровень логирования быть обработан
func (c *otlpCore) Enabled(level zapcore.Level) bool {
	return c.level.Enabled(level)
}

// With добавляет поля к core
func (c *otlpCore) With(fields []zapcore.Field) zapcore.Core {
	cloned := c.encoder.Clone()
	for _, field := range fields {
		field.AddTo(cloned)
	}

	newFields := make([]zapcore.Field, len(c.fields), len(c.fields)+len(fields))
	copy(newFields, c.fields)
	newFields = append(newFields, fields...)

	return &otlpCore{
		encoder: cloned,
		writer:  c.writer,
		level:   c.level,
		fields:  newFields,
	}
}

// Check проверяет, нужно ли логировать запись, и возвращает Entry
func (c *otlpCore) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return checked.AddCore(entry, c)
	}

	return checked
}

// Write конвертирует zap.Entry в OTLP log.Record и отправляет через logger
func (c *otlpCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	allFields := fields
	if len(c.fields) > 0 {
		allFields = make([]zapcore.Field, 0, len(c.fields)+len(fields))
		allFields = append(allFields, c.fields...)
		allFields = append(allFields, fields...)
	}

	record := convertZapEntryToOTLPRecord(entry, allFields)

	c.writer.logger.Emit(context.Background(), record)

	return nil
}

// Sync сбрасывает буферы
func (c *otlpCore) Sync() error {
	return c.writer.Sync()
}

// convertZapEntryToOTLPRecord конвертирует zap.Entry в OTLP log.Record
func convertZapEntryToOTLPRecord(entry zapcore.Entry, fields []zapcore.Field) log.Record {
	record := log.Record{}
	record.SetTimestamp(entry.Time)

	record.SetSeverity(convertZapLevelToOTLPSeverity(entry.Level))
	record.SetSeverityText(entry.Level.String())

	record.SetBody(log.StringValue(entry.Message))

	attrs := make([]log.KeyValue, 0, len(fields)+3)

	if entry.Caller.Defined {
		attrs = append(attrs,
			log.String("log.file.name", entry.Caller.File),
			log.Int("log.file.line", entry.Caller.Line),
			log.String("log.file.function", entry.Caller.Function),
		)
	}

	for _, field := range fields {
		attrs = append(attrs, convertZapFieldToLogKeyValue(field))
	}

	record.AddAttributes(attrs...)

	return record
}

// convertZapLevelToOTLPSeverity конвертирует zap.Level в OTLP severity
func convertZapLevelToOTLPSeverity(level zapcore.Level) log.Severity {
	switch level {
	case zapcore.DebugLevel:
		return log.SeverityDebug
	case zapcore.InfoLevel:
		return log.SeverityInfo
	case zapcore.WarnLevel:
		return log.SeverityWarn
	case zapcore.ErrorLevel:
		return log.SeverityError
	case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		return log.SeverityFatal
	default:
		return log.SeverityInfo
	}
}

// convertZapFieldToLogKeyValue конвертирует zap.Field в log.KeyValue
//
//nolint:cyclop
func convertZapFieldToLogKeyValue(field zapcore.Field) log.KeyValue {
	key := field.Key

	switch field.Type {
	case zapcore.StringType:
		return log.String(key, field.String)
	case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type:
		return log.Int64(key, field.Integer)
	case zapcore.Uint64Type, zapcore.Uint32Type, zapcore.Uint16Type, zapcore.Uint8Type, zapcore.UintptrType:
		return log.Int64(key, field.Integer)
	case zapcore.Float64Type:
		//nolint:gosec
		return log.Float64(key, math.Float64frombits(uint64(field.Integer)))
	case zapcore.Float32Type:
		//nolint:gosec
		return log.Float64(key, float64(math.Float32frombits(uint32(field.Integer))))
	case zapcore.BoolType:
		return log.Bool(key, field.Integer == 1)
	case zapcore.DurationType:
		return log.String(key, time.Duration(field.Integer).String())
	case zapcore.TimeType:
		t := time.Unix(0, field.Integer)

		if field.Interface != nil {
			if loc, ok := field.Interface.(*time.Location); ok && loc != nil {
				t = t.In(loc)
			}
		}

		return log.String(key, t.Format(time.RFC3339Nano))
	case zapcore.TimeFullType:
		if t, ok := field.Interface.(time.Time); ok {
			return log.String(key, t.Format(time.RFC3339Nano))
		}

		return log.String(key, fmt.Sprintf("%v", field.Interface))
	case zapcore.ErrorType:
		if err, ok := field.Interface.(error); ok && err != nil {
			return log.String(key, err.Error())
		}

		return log.String(key, "<nil>")
	case zapcore.StringerType:
		if stringer, ok := field.Interface.(fmt.Stringer); ok && stringer != nil {
			return log.String(key, stringer.String())
		}

		return log.String(key, "<nil>")
	case zapcore.BinaryType:
		if b, ok := field.Interface.([]byte); ok {
			return log.String(key, string(b))
		}

		return log.String(key, fmt.Sprintf("%v", field.Interface))
	case zapcore.ByteStringType:
		return log.String(key, field.String)
	case zapcore.Complex128Type, zapcore.Complex64Type:
		return log.String(key, fmt.Sprintf("%v", field.Interface))
	case zapcore.ReflectType:
		return log.String(key, fmt.Sprintf("%v", field.Interface))
	case zapcore.NamespaceType:
		return log.String(key, "<namespace>")
	case zapcore.SkipType:
		return log.String(key, "<skipped>")
	default:
		return log.String(key, fmt.Sprintf("%v", field.Interface))
	}
}
