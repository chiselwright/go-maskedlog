package maskedlog

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// MaskStrings is a convenience type for the list of sensitive values to hide away
type MaskStrings []string

// MaskLog is a structure to hold "useful" state
type MaskLog struct {
	SensitiveStrings *MaskStrings
	Opts             interface{}
}

var once sync.Once

var (
	maskStrings MaskStrings
)

// GetSingleton ...
func GetSingleton() MaskLog {
	// via: https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
	once.Do(func() { // <-- atomic, does not allow repeating
		maskStrings = make(MaskStrings, 0) // <-- thread safe
	})

	ml := MaskLog{
		SensitiveStrings: &maskStrings,
		Opts:             nil,
	}

	return ml
}

// Reset ...
func (ml *MaskLog) Reset() {
	maskStrings = make(MaskStrings, 0)
}

// AddSensitiveValue adds a new token/password value to mask in any log output
func (ml *MaskLog) AddSensitiveValue(s string) {
	maskStrings = append(maskStrings, s)
}

// SafeString ...
func SafeString(s string) string {
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
func (ml MaskLog) SanitizeInterfaceValues(z interface{}) {
	for i := range z.([]interface{}) {
		q := z.([]interface{})[i]

		// only replace (in) strings
		_, ok := q.(string)
		if ok {
			for _, ms := range *ml.SensitiveStrings {
				q = strings.ReplaceAll(q.(string), ms, SafeString(ms))
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

// PrepareMessage ...
func (ml MaskLog) PrepareMessage(v interface{}) string {
	setLogFormat()
	ml.SanitizeInterfaceValues(v)
	msg := Stringify(v)
	return msg
}

// LogWarn does, uhm, warnings?
func (ml MaskLog) LogWarn(v ...interface{}) {
	log.Warn().Msg(ml.PrepareMessage(v))
}

// LogVerbose does other things
func (ml MaskLog) LogVerbose(v ...interface{}) {
	/*
		if *vars.Verbose {
			log.Trace().Msg(ml.SensitiveStrings.PrepareMessage(v))
		}
	*/
}

// LogFatal does cool things
func (ml MaskLog) LogFatal(v ...interface{}) {
	log.Fatal().Msg(ml.PrepareMessage(v))
}
