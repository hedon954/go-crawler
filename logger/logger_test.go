package logger

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestNewFilePlugin(t *testing.T) {
	filePrefix := "test"
	fileSuffix := ".log"
	gzipSuffix := ".gz"

	// Creates a logger
	plugin, closer := NewFilePlugin(filePrefix+fileSuffix, zapcore.DebugLevel)
	logger := NewLogger(plugin)

	// Simulates to write content to log
	bs := make([]byte, 10000)
	for count := 10000; count > 0; count-- {
		logger.Info(string(bs))
	}

	// Waits for writing completely and
	// lumberjac to compress the log files
	err := closer.Close()
	require.NoError(t, err)
	time.Sleep(1 * time.Second)

	// Check result
	dir, err := ioutil.ReadDir(".")
	require.NoError(t, err)
	var gzCount, logCount int
	for _, f := range dir {
		fn := f.Name()
		if strings.HasSuffix(fn, fileSuffix) {
			logCount++
			assert.NoError(t, os.Remove(f.Name()))
			continue
		}
		if strings.HasSuffix(fn, gzipSuffix) {
			gzCount++
			logCount++
			assert.NoError(t, os.Remove(f.Name()))
			continue
		}
	}

	require.Equal(t, 3, logCount)
	require.Equal(t, 2, gzCount)
}
