package log

import (
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func Test_NewFilePlugin(t *testing.T) {
	const filePrefix = "test"
	const fileSuffix = ".log"
	const gzipSuffix = ".gz"

	fp, closer := NewFilePlugin(filePrefix+fileSuffix, zapcore.DebugLevel)
	logger := NewLogger(fp)
	buf := make([]byte, 10000)
	count := 10000
	for count > 0 {
		count--
		logger.Info(string(buf))
	}
	err := closer.Close()
	require.NoError(t, err)

	// NOTE: With the current implementation of Lumberjack,
	// close does not stop the compression coroutine.
	// So here we need to wait for Lumberjack to compress the log file.
	time.Sleep(1 * time.Second)

	fs, err := ioutil.ReadDir(".")
	require.NoError(t, err)
	gzCount := 0
	logCount := 0
	for _, f := range fs {
		fn := f.Name()
		if strings.HasPrefix(fn, filePrefix) {
			if strings.HasSuffix(fn, fileSuffix) {
				logCount++
				//assert.NoError(t, os.Remove(f.Name()))
				continue
			}
			if strings.HasSuffix(fn, fileSuffix+gzipSuffix) {
				gzCount++
				logCount++
				//assert.NoError(t, os.Remove(f.Name()))
				continue
			}
		}
	}

	require.Equal(t, 3, logCount)
	require.Equal(t, 2, gzCount)
}
