package utils

import (
	"bufio"
	"fmt"
	"chat/library/config"
	"io"
	"os"
	"strings"
)

// @Description 输出字符串到控制台，并接收输入值到value
func DumpSAndScanVar(str string, value *string) {
	fmt.Println(str)
	inputStr := ""
	for {
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadBytes('\n')
		if err != nil || err == io.EOF {
			fmt.Println("老铁，你的输入有问题哦！")
			continue
		}
		inputStr = string(input)
		if CheckSubstrings(inputStr, config.GetConfig().Verify.IllegalStr...) != 0 {
			fmt.Println("老铁，请不要输入非法字符哦")
			continue
		}
		break
	}
	inputStr = strings.Replace(inputStr, "\n", "", -1)
	*value = strings.Replace(inputStr, "\r", "", -1)
}

// @Description debug开启时打印(客户端)
func CDD(str string) {
	if config.GetConfig().App.ClientDebug {
		fmt.Println(str)
	}
}

// @Description debug开启时打印(服务端)
func SDD(str string) {
	if config.GetConfig().App.ServerDebug {
		fmt.Println(str)
	}
}

//@Description debug开启时打印(公共)
func DD(str string) {
	if config.GetConfig().App.ServerDebug || config.GetConfig().App.ClientDebug {
		fmt.Println(str)
	}
}

// @Description 判断str中是否包含指定字符串，返回出现的次数，重复出现只算一次，subs为空串也会成立哦
func CheckSubstrings(str string, subs ...string) int {
	matches := 0
	for _, sub := range subs {
		if strings.Contains(str, sub) {
			matches += 1
		}
	}
	return matches
}

// @Description 删除字符串切片指定元素
func StrSliceDelete(target string, search []string) (result []string) {
	for _, v := range search {
		if v == target {
			continue
		}
		result = append(result, v)
	}
	return result
}

// @Description 空接口切片转字符串切片
func InterfaceSliceToStrSlice(data []interface{}) (res []string) {
	for _, v := range data {
		if str, ok := v.(string); ok {
			res = append(res, str)
		}
	}
	return res
}

//字符串切片去重
func RemoveDuplicateElement(languages []string) []string {
	result := make([]string, 0, len(languages))
	temp := map[string]struct{}{}
	for _, item := range languages {
		if _, ok := temp[item]; !ok { //如果字典中找不到元素，ok=false，!ok为true，就往切片中append元素。
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}