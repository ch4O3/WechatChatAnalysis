package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/gogf/gf/text/gregex"
)

type User struct {
	username string
	time     string
	content  string
}

type MapSorter []Item

type Item struct {
	Key string
	Val int
}

func NewMapSorter(m map[string]int) MapSorter {
	ms := make(MapSorter, 0, len(m))
	for k, v := range m {
		ms = append(ms, Item{k, v})
	}
	return ms
}

func (ms MapSorter) Len() int {
	return len(ms)
}

func (ms MapSorter) Less(i, j int) bool {
	return ms[i].Val > ms[j].Val
}

func (ms MapSorter) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

func main() {
	var file string

	flag.StringVar(&file, "i", "", "")

	flag.CommandLine.Usage = func() {
		fmt.Println("使用说明：./main -i 1.txt")
	}
	flag.Parse()

	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	//初始化userCount，用于存放群聊的用户名及发言次数。
	userCount := make(map[string]int)
	//初始化content集合
	contentCi := 0

	//图片数量
	Piccount := 0

	//开发语言
	gocount := 0
	javacount := 0
	netcount := 0
	phpcoumt := 0
	jsconut := 0
	pythoncount := 0

	for scanner.Scan() {
		patternStr := `^.{1,10}\s[(][2][0][2][0-2]-[0-9][0-9]-[0-9][0-9]\s[0-9][0-9]:[0-9][0-9]:[0-9][0-9][)]:.*`
		if gregex.IsMatchString(patternStr, scanner.Text()) {
			a := strings.Split(scanner.Text(), " (")
			b := strings.Split(a[1], "):")
			t := User{}
			t.username = a[0]
			t.time = b[0]
			t.content = b[1]
			contentCi++
			//群聊用户出现次数
			v, ok := userCount[t.username]
			if ok {
				userCount[t.username] = v + 1
			} else {
				userCount[t.username] = 1
			}
			//群聊聊过的内容,分词
			if strings.Contains(t.content, "[图片]") {
				Piccount++
			}
			if strings.Contains(t.content, "go") {
				gocount++
			}
			if strings.Contains(t.content, "java") {
				javacount++
			}
			if strings.Contains(t.content, ".net") {
				netcount++
			}
			if strings.Contains(t.content, "php") {
				phpcoumt++
			}
			if strings.Contains(t.content, "js") {
				jsconut++
			}
			if strings.Contains(t.content, "python") {
				pythoncount++
			}

		}
	}
	//统计群聊中发言Top10并输出
	ms := NewMapSorter(userCount)
	sort.Sort(ms)
	count := 0
	for _, item := range ms {
		count++
		fmt.Printf("发言次数 Top%v 用户是： %s ", count, item.Key)
		fmt.Printf("\n发言次数：%d\n\n", item.Val)
		if count == 10 {
			break
		}
	}
	//图片统计
	fmt.Printf("聊天记录存在 %v张 图片\n\n", Piccount)

	//发言统计
	fmt.Printf("摸鱼摸了%v条\n\n", contentCi)

	//开发语言讨论分析
	language := map[string]int{
		"Go":         gocount,
		"Java":       javacount,
		".Net":       netcount,
		"Php":        phpcoumt,
		"Javascript": jsconut,
		"Python":     pythoncount,
	}
	ab := NewMapSorter(language)
	sort.Sort(ab)
	for _, item := range ab {
		fmt.Printf("%v语言讨论%d次\n", item.Key, item.Val)
	}
}
