package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

type CommandOptions struct {
	pathname    string
	accept      []string
	recursion   bool
	withDataURI bool
}

func containsString(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

func timeNow() string {
	return time.Now().Format("2006-01-01 02:02:02")
}

func printInfo() {
	fmt.Printf("[INFO %v]", timeNow())
}

func printError(errors []error) {
	fmt.Printf("[ERR %v]", timeNow())
	for _, err := range errors {
		fmt.Println(err.Error())
	}
}

/**
 * 生成base64文件
 */
func makeFile(filename string, prefix bool) error {
	data, err := os.ReadFile(filename)

	if err != nil {
		return err
	}

	encoded := base64.StdEncoding.EncodeToString(data)

	// 添加DataURI前缀
	if prefix {
		mimes := map[string]string{
			".png":  "image/png",
			".jpg":  "image/jpeg",
			".jpeg": "image/jpeg",
			".ico":  "image/x-icon",
			".ttf":  "font/truetype;charset=utf-8",
		}
		for ext, mime := range mimes {
			if strings.HasSuffix(filename, ext) {
				encoded = "data:" + mime + ";base64," + encoded
			}
		}
	}
	return os.WriteFile(filename+".b64.txt", []byte(encoded), 0644)
}

/**
 * 遍历目录生成base64文件
 */
func eachDir(options CommandOptions) []error {
	var errors []error

	entries, readDirError := os.ReadDir(options.pathname)

	if readDirError != nil {
		return []error{
			readDirError,
		}
	}

	for _, entry := range entries {
		info, entryInfoError := entry.Info()
		if entryInfoError != nil {
			errors = append(errors, entryInfoError)
			continue
		}

		name := info.Name()

		if info.IsDir() {
			// 跳过隐藏目录
			if strings.HasPrefix(name, ".") && !strings.HasPrefix(name, "./") {
				continue
			}
			if options.recursion {
				eachDirError := eachDir(CommandOptions{
					path.Join(options.pathname, name),
					options.accept,
					options.recursion,
					options.withDataURI,
				})
				if eachDirError != nil {
					errors = append(errors, eachDirError...)
				}
			}
			continue
		}

		// 跳过 .b64.txt
		if strings.HasSuffix(name, ".b64.txt") || strings.HasPrefix(name, ".") {
			continue
		}
		// 未在设定文件类型中
		if !containsString(options.accept, path.Ext(name)) {
			continue
		}

		filename := path.Join(options.pathname, name)

		makeFileError := makeFile(filename, options.withDataURI)
		if makeFileError != nil {
			errors = append(errors, makeFileError)
		}
	}
	return errors
}

func main() {
	args := os.Args

	options := CommandOptions{
		pathname:    "",
		accept:      []string{".png", ".jpg", ".ico", ".ttf", ".txt"},
		withDataURI: false,
		recursion:   false,
	}

	var arguments []string
	// 用户设置 --accept 时过滤文件扩展名称
	for _, argv := range args[1:] {
		if strings.HasPrefix(argv, "-V") {
			script := args[0:1][0]
			scriptInfo, err := os.Stat(script)
			if err == nil {
				fmt.Println(scriptInfo.ModTime().Format("2006-01-02 15:04:05"))
			}
			return
		}
		if strings.HasPrefix(argv, "--accept=") {
			pieces := strings.Split(argv, "=")
			for _, ext := range strings.Split(pieces[1], ",") {
				if strings.HasPrefix(ext, ".") {
					options.accept = append(options.accept, ext)
				} else {
					options.accept = append(options.accept, "."+ext)
				}
			}
		} else if strings.HasPrefix(argv, "-R") {
			options.recursion = true
		} else if strings.HasPrefix(argv, "-P") {
			options.withDataURI = true
		} else if !strings.HasPrefix(argv, "-") {
			arguments = append(arguments, argv)
		}
	}

	// 自定义目录或文件
	if 1 == len(arguments) {
		options.pathname = strings.TrimRight(strings.TrimSpace(arguments[0]), "/")

		info, statError := os.Stat(options.pathname)

		if statError != nil {
			printError([]error{statError})
			return
		}

		// 目录
		if info.IsDir() {
			errors := eachDir(options)
			if len(errors) > 0 {
				printError(errors)
			}
			return
		}
		makeFileError := makeFile(options.pathname, options.withDataURI)
		if makeFileError == nil {
			printInfo()
		} else {
			printError([]error{makeFileError})
		}
		return
	}

	pathname, err := os.Getwd()
	if err == nil {
		errors := eachDir(CommandOptions{
			pathname,
			options.accept,
			options.recursion,
			options.withDataURI,
		})
		if len(errors) > 0 {
			printError(errors)
		} else {
			printInfo()
		}
		return
	}
	printError([]error{err})
}
