package service

import (
	"bytes"
	"testing"

	"github.com/Jeffail/benthos/v3/lib/log"
	"github.com/stretchr/testify/assert"
)

func TestReverseAirGapLogger(t *testing.T) {
	lConf := log.NewConfig()
	lConf.AddTimeStamp = false

	var buf bytes.Buffer
	logger := log.New(&buf, lConf)

	agLogger := newReverseAirGapLogger(logger)
	agLogger2 := agLogger.With("field1", "value1", "field2", "value2")

	agLogger.Debugf("foo: %v", "bar1")
	agLogger.Infof("foo: %v", "bar2")

	agLogger2.Debugf("foo2: %v", "bar1")
	agLogger2.Infof("foo2: %v", "bar2")

	agLogger.Warnf("foo: %v", "bar3")
	agLogger.Errorf("foo: %v", "bar4")

	agLogger2.Warnf("foo2: %v", "bar3")
	agLogger2.Errorf("foo2: %v", "bar4")

	assert.Equal(t, `{"@service":"benthos","component":"benthos","level":"INFO","message":"foo: bar2"}
{"@service":"benthos","component":"benthos","field1":"value1","field2":"value2","level":"INFO","message":"foo2: bar2"}
{"@service":"benthos","component":"benthos","level":"WARN","message":"foo: bar3"}
{"@service":"benthos","component":"benthos","level":"ERROR","message":"foo: bar4"}
{"@service":"benthos","component":"benthos","field1":"value1","field2":"value2","level":"WARN","message":"foo2: bar3"}
{"@service":"benthos","component":"benthos","field1":"value1","field2":"value2","level":"ERROR","message":"foo2: bar4"}
`, buf.String())
}

func TestReverseAirGapLoggerDodgyFields(t *testing.T) {
	lConf := log.NewConfig()
	lConf.AddTimeStamp = false

	var buf bytes.Buffer
	logger := log.New(&buf, lConf)

	agLogger := newReverseAirGapLogger(logger)

	agLogger.With("field1", "value1", "field2").Infof("foo1")
	agLogger.With(10, 20).Infof("foo2")
	agLogger.With("field3", 30).Infof("foo3")
	agLogger.With("field4", "value4").With("field5", "value5").Infof("foo4")

	assert.Equal(t, `{"@service":"benthos","component":"benthos","field1":"value1","level":"INFO","message":"foo1"}
{"10":"20","@service":"benthos","component":"benthos","level":"INFO","message":"foo2"}
{"@service":"benthos","component":"benthos","field3":"30","level":"INFO","message":"foo3"}
{"@service":"benthos","component":"benthos","field4":"value4","field5":"value5","level":"INFO","message":"foo4"}
`, buf.String())
}
