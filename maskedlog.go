package maskedlog

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// MaskStrings is a convenience type for the list of sensitive values to hide away
type MaskStrings []string

var maskStrings = MaskStrings{}

func safeString(s string) string {
	rs := []rune(s)
	const replaceChar = 'x'
	const preserveChar = '-'
	var start, end int

	// mask "short things" slightly differently
	if len(rs) <= 10 {
		start = 2
		end = len(rs) - 1
	} else {
		start = 4
		end = len(rs) - 4
	}

	for i := start; i < end; i++ {
		if rs[i] != preserveChar {
			rs[i] = replaceChar
		}
	}

	return string(rs)
}

// SanitizeInterfaceValues will mask sensitive values in output
func SanitizeInterfaceValues(z interface{}) {
	for i := range z.([]interface{}) {
		q := z.([]interface{})[i]

		// only replace (in) strings
		_, ok := q.(string)
		if ok {
			for _, ms := range maskStrings {
				q = strings.ReplaceAll(q.(string), ms, safeString(ms))
			}
			z.([]interface{})[i] = q
		}
	}
}

// Stringify ... turns things into strings
func Stringify(z interface{}) string {
	var data []string

	for i := range z.([]interface{}) {
		q := z.([]interface{})[i]
		data = append(data, fmt.Sprintf("%+v", q))
	}

	return strings.Join(data, " ")
}

func setLogFormat() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
}

func prepareMessage(v interface{}) string {
	setLogFormat()
	SanitizeInterfaceValues(v)
	msg := Stringify(v)
	return msg
}

// LogWarn does, uhm, warnings?
func LogWarn(v ...interface{}) {
	log.Warn().Msg(prepareMessage(v))
}

// LogVerbose does other things
func LogVerbose(v ...interface{}) {
	/*
		if *vars.Verbose {
			log.Trace().Msg(prepareMessage(v))
		}
	*/
}

// LogFatal does cool things
func LogFatal(v ...interface{}) {
	log.Fatal().Msg(prepareMessage(v))
}
