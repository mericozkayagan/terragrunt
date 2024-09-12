package terraform

import (
	"regexp"
	"strings"
	"time"

	"github.com/gruntwork-io/go-commons/errors"
	"github.com/gruntwork-io/terragrunt/pkg/log"
	"github.com/gruntwork-io/terragrunt/pkg/log/writer"
)

const (
	// logTimestampFormat is TF_LOG timestamp format
	logTimestampFormat = "2006-01-02T15:04:05.000-0700"
)

var (
	// tfLogTimeLevelMsgReg is a regular expression that matches TF_LOG output, example output:
	//
	// 2024-09-08T13:44:31.229+0300 [DEBUG] using github.com/zclconf/go-cty v1.14.3
	// 2024-09-08T13:44:31.229+0300 [INFO]  Go runtime version: go1.22.1
	tfLogTimeLevelMsgReg = regexp.MustCompile(`(?i)(^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}[-+]\d{4})\s*\[(trace|debug|warn|info|error)\]\s*(.+\S)$`)
)

func ParseLogFunc(msgPrefix string) writer.WriterParseFunc {
	return func(str string) (msg string, ptrTime *time.Time, ptrLevel *log.Level, err error) {
		return ParseLog(msgPrefix, str)
	}
}

func ParseLog(msgPrefix, str string) (msg string, ptrTime *time.Time, ptrLevel *log.Level, err error) {
	const numberOfValues = 4

	if !tfLogTimeLevelMsgReg.MatchString(str) {
		return str, nil, nil, nil
	}

	match := tfLogTimeLevelMsgReg.FindStringSubmatch(str)
	if len(match) != numberOfValues {
		return str, nil, nil, nil
	}

	timeStr, levelStr, msg := match[1], match[2], match[3]

	if timeStr != "" {
		time, err := time.Parse(logTimestampFormat, timeStr)
		if err != nil {
			return "", nil, nil, errors.WithStackTrace(err)
		}

		ptrTime = &time
	}

	if levelStr != "" {
		level, err := log.ParseLevel(strings.ToLower(levelStr))
		if err != nil {
			return "", nil, nil, errors.WithStackTrace(err)
		}

		ptrLevel = &level
	}

	return msgPrefix + msg, ptrTime, ptrLevel, nil
}