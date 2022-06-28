package utility

//General function
import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"sort"
)

//Check err and print it
func IsErr(err error, msg string) bool {
	if err != nil {
		log.Println("ERROR: "+msg+"\n", err)
		return true
	}
	return false
}

//Sort available servers and return the best
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

//Notice the selected Server to start service
func NoticeServer(ip string, service string) {

	url := "http://" + ip + ":8000/start"

	//json序列化
	post := "{\"service\":\"" + service +
		"\"}"

	var jsonStr = []byte(post)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println(err)
		return
	}
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

//Data struct for sort
type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
