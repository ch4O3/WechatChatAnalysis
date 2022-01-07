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
	return ms[i].Val > ms[j].Val // 按值排序
	//return ms[i].Key < ms[j].Key // 按键排序
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

	for scanner.Scan() {
		patternStr := `^.{1,10}\s[(][2][0][2][0-2]-[0-9][0-9]-[0-9][0-9]\s[0-9][0-9]:[0-9][0-9]:[0-9][0-9][)]:.*`
		if gregex.IsMatchString(patternStr, scanner.Text()) {
			a := strings.Split(scanner.Text(), " (")
			b := strings.Split(a[1], "):")
			t := User{}
			t.username = a[0]
			t.time = b[0]
			t.content = b[1]
			//群聊用户出现次数
			v, ok := userCount[t.username]
			if ok {
				userCount[t.username] = v + 1
			} else {
				userCount[t.username] = 1
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
}
