package logging

import (
	"fmt"
	"io"
	"strings"
)

type Logger struct {
	builder strings.Builder
}

func (l *Logger) Log(sender string, msg string, a ...any) {
	l.builder.WriteString(fmt.Sprintf("%s: %s\n", strings.ToUpper(sender), fmt.Sprintf(msg, a...)))
}

func (l *Logger) Dump(writer io.Writer) {
	writer.Write([]byte(l.builder.String()))
}

var AppLogger Logger
