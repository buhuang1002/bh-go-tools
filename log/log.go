package bh_log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func Warn(a ...interface{}) {
	fmt.Print(warnColor(logger(a...)))
}

func Info(a ...interface{}) {
	fmt.Print(infoColor(logger(a...)))
}

func Important(a ...interface{}) {
	fmt.Print(importanceColor(logger(a...)))
}

func EchoWarn(a ...interface{}) {
	fmt.Print(warnColor(echo(a...)))
}

func EchoInfo(a ...interface{}) {
	fmt.Print(infoColor(echo(a...)))
}

func EchoImportant(a ...interface{}) {
	fmt.Print(importanceColor(echo(a...)))
}

func EchoWarnf(format string, args ...interface{}) {
	fmt.Print(warnColor(echo(fmt.Sprintf(format, args...))))
}

func EchoInfof(format string, args ...interface{}) {
	fmt.Print(infoColor(echo(fmt.Sprintf(format, args...))))
}

func EchoImportantf(format string, args ...interface{}) {
	fmt.Print(importanceColor(echo(fmt.Sprintf(format, args...))))
}

func Warnf(format string, args ...interface{}) {
	fmt.Print(warnColor(logger(fmt.Sprintf(format, args...))))
}

func Infof(format string, args ...interface{}) {
	fmt.Print(infoColor(logger(fmt.Sprintf(format, args...))))
}

func Importantf(format string, args ...interface{}) {
	fmt.Print(importanceColor(logger(fmt.Sprintf(format, args...))))
}

func PrintFuncSet(level int) {
	if level == 0 {
		level = 99999
	}

	funcInfoList := []string{}
	defer func() {
		fmt.Println(commonColor("\n ------------------------begin------------------------"))
		for i := len(funcInfoList) - 1; i >= 0; i-- {
			fmt.Print(commonColor(funcInfoList[i]))
		}
		fmt.Println(commonColor("-------------------------end-------------------------\n"))
	}()

	for i := 1; i < level+1; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			return
		}
		fn := runtime.FuncForPC(pc)
		switch {
		case strings.HasSuffix(fn.Name(), "main") ||
			strings.HasSuffix(fn.Name(), "goexit"):
			return
		}
		funcInfoList = append(funcInfoList, fmt.Sprintf("%s:%d [%s]\n", simplePath(file), line, simpleFunc(fn.Name())))
	}
}

func PrintFunc(level int, args ...interface{}) {
	pc, _, _, _ := runtime.Caller(level + 1)
	fn := runtime.FuncForPC(pc)
	fmt.Print(commonColor(fmt.Sprintf("%s [%s] ", time.Now().Format("15:04:05.999"), simpleFunc(fn.Name())) + fmt.Sprintln(args...)))
}

//utils

func logger(a ...interface{}) string {
	return caller(3) + ": " + fmt.Sprintln(a...)
}

func echo(a ...interface{}) string {
	return time.Now().Format("15:04:05.999") + " " + fmt.Sprintln(a...)
}

var goRoot = filepath.Join(runtime.GOROOT(), "src")

var goPath = os.Getenv("GOPATH")

var home = os.Getenv("HOME")

var (
	longTime bool = false
	longPath bool = false
)

func currentTime() string {
	if longTime {
		return time.Now().Format("2006-01-02 15:04:05.999")
	}
	return time.Now().Format("15:04:05.999")
}

func simplePath(path string) string {
	if longPath {
		return path
	}
	switch {
	case len(goRoot) > 0 && strings.HasPrefix(path, goRoot):
		return path[len(goRoot)+1:]
	case len(goPath) > 0 && strings.HasPrefix(path, goPath):
		return path[len(goPath)+1:]
	case len(home) > 0 && strings.HasPrefix(path, home):
		return path[len(home)+1:]
	default:
		return path
	}
}

func simpleFunc(name string) string {
	if name == "" {
		return name
	}
	paths := strings.Split(name, "/")
	//record := strings.Split(paths[len(paths)-1], ".")
	return paths[len(paths)-1]
}

func caller(skip int) string {
	pc, file, line, _ := runtime.Caller(skip)
	fn := runtime.FuncForPC(pc)
	return fmt.Sprintf("%s %s:%d [%s]", currentTime(), simplePath(file), line, simpleFunc(fn.Name()))
}

//color

func warnColor(s string) string {
	return printColor(s, 48, 33, 1)
}

func infoColor(s string) string {
	return printColor(s, 48, 32, 1)
}

func importanceColor(s string) string {
	return printColor(s, 48, 31, 1)
}

func commonColor(s string) string {
	return printColor(s, 48, 36, 1)
}

func printColor(s string, bColor, fontColor, mode int) string {
	return fmt.Sprintf(" %c[%d;%d;%dm%s%c[0m", 0x1B, mode, bColor, fontColor, s, 0x1B)
}
