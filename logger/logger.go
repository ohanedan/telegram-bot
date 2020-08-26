package logger

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	Yellow = color.New(color.FgYellow).SprintFunc()
	Red    = color.New(color.FgRed).SprintFunc()
	Green  = color.New(color.FgGreen).SprintFunc()
	Blue   = color.New(color.FgBlue).SprintFunc()
	Cyan   = color.New(color.FgCyan).SprintFunc()
)

type Logger struct {
	Disabled bool
}

func (this *Logger) Printf(sender string, format string, a ...interface{}) (int, error) {
	if this.Disabled {
		return 0, nil
	}
	format = parseColors(format)
	time := time.Now()
	format = color.CyanString("[%v]", time.Format("2 Jan 2006 15:04:05")) +
		color.YellowString("[%v] ", sender) + format

	n, err := fmt.Printf(format, a...)
	return n, err
}

func (this *Logger) Sprintf(format string, a ...interface{}) string {
	format = parseColors(format)

	result := fmt.Sprintf(format, a...)
	return result
}

func (this *Logger) Println(sender string, format string, a ...interface{}) (int, error) {
	if this.Disabled {
		return 0, nil
	}
	format = parseColors(format)
	time := time.Now()
	format = fmt.Sprintf("%v%v %v\n",
		color.CyanString("[%v]", time.Format("2 Jan 2006 15:04:05")),
		color.YellowString("[%v]", sender), format)

	n, err := fmt.Printf(format, a...)
	return n, err
}

func parseColors(format string) string {
	allowedColors := []string{"yellow", "red", "green", "blue", "cyan"}

	for _, _color := range allowedColors {
		colorString := _color + "{"
		index := strings.Index(format, colorString)
		for index != -1 {
			text := ""
			lastIndex := -1
			for i := index + len(colorString); i < len(format); i++ {
				if format[i] == '}' {
					lastIndex = i
					break
				}
				text = text + string(format[i])
			}
			start := format[:index]
			end := format[lastIndex+1:]
			switch _color {
			case "yellow":
				text = color.YellowString(text)
			case "red":
				text = color.RedString(text)
			case "green":
				text = color.GreenString(text)
			case "blue":
				text = color.BlueString(text)
			case "cyan":
				text = color.CyanString(text)
			}
			format = start + text + end

			index = strings.Index(format, colorString)
		}
	}

	return format
}
