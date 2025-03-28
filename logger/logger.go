package logger

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"strconv"
)

var Logg *logrus.Logger

type ColorLogger struct {
	logrus.TextFormatter
	ForceColors   bool
	ColorInfo     *color.Color
	ColorWarning  *color.Color
	ColorError    *color.Color
	ColorCritical *color.Color
}

func (c *ColorLogger) Format(entry *logrus.Entry) ([]byte, error) {
	if c.ForceColors {
		switch entry.Level {
		case logrus.InfoLevel:
			c.ColorInfo.Println(entry.Message)
		case logrus.WarnLevel:
			c.ColorWarning.Println(entry.Message)
		case logrus.ErrorLevel:
			c.ColorError.Println(entry.Message)
		case logrus.FatalLevel, logrus.PanicLevel:
			c.ColorCritical.Println(entry.Message)
		default:
			c.PrintColored(entry)
		}
		return nil, nil
	} else {
		return c.TextFormatter.Format(entry)
	}
}

func (c *ColorLogger) PrintColored(entry *logrus.Entry) {
	levelColor := color.New(color.FgCyan, color.Bold)             // 定义蓝色和粗体样式
	levelText := levelColor.Sprintf("%-6s", entry.Level.String()) // 格式化日志级别文本

	msg := levelText + " " + entry.Message
	if entry.HasCaller() {
		msg += " (" + entry.Caller.File + ":" + strconv.Itoa(entry.Caller.Line) + ")" // 添加调用者信息
	}

	fmt.Fprintln(color.Output, msg) // 使用有颜色的方式打印消息到终端
}

func NewLogger() {
	Logg = logrus.New()
	Logg.Formatter = &ColorLogger{
		ForceColors:   true,
		ColorInfo:     color.New(color.FgBlue),
		ColorWarning:  color.New(color.FgYellow),
		ColorError:    color.New(color.FgRed),
		ColorCritical: color.New(color.BgRed, color.FgWhite),
	}
}
