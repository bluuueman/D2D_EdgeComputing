package utility

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"sort"
)

func IsErr(err error, msg string) bool {
	if err != nil {
		log.Println("ERROR: "+msg+"\n", err)
		return true
	}
	return false
}

func Rank(wordFrequencies map[string]int) []string {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i += 1
	}
	//从小到大排序
	//sort.Sort(pl)
	//从大到小排序
	sort.Sort(sort.Reverse(pl))
	result := make([]string, i)
	for idx, _ := range pl {
		result[idx] = pl[idx].Key
	}
	return result
}

func NoticeServer(ip string, service string, port string) {
        //time.Sleep(time.Duration(2)*time.Second)
	url := "http://" + ip + ":8081/start"

	//json序列化
	post := "{\"service\":\"" + service +
		"\",\"port\":\"" + port +
		"\"}"

	var jsonStr = []byte(post)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
        if err!=nil{
                fmt.Println(err)
                return
        }
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
                return
	}
	defer resp.Body.Close()
	fmt.Println("status", resp.Status)
	return
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
