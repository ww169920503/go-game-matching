package matcher_test

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/ww169920503/go-game-matching/matcher"
)
/*
func benchmarkMatcher_Match(b *testing.B, concurrent int) {
	const maxScore = 300
	var m = matcher.NewMatcher(120, maxScore, 10)
	var mu sync.Mutex
	isRun := true

	for i := 0; i < concurrent; i++ {
		go func(base int) {
			for id := base; isRun; id++ {
				score := rand.Uint32() % maxScore
				mu.Lock()
				err := m.JoinQueue(matcher.PlayerId(strconv.Itoa(id)), matcher.Time(time.Now().Unix()), matcher.PlayerScore(score))
				mu.Unlock()
				if err != nil {
					panic(err)
				}
			}
		}(i * 1000000)
	}

	for {
		mu.Lock()
		m.Match(matcher.Time(time.Now().Unix()), 25)
		mu.Unlock()
		if m.PlayerCount()-m.PlayerInQueueCount() > b.N {
			isRun = false
			break
		}
		time.Sleep(time.Microsecond * 1000)
	}
}

// -test.benchtime 10s
func BenchmarkMatcher_Match(b *testing.B) {
	benchmarkMatcher_Match(b, 1)
}

// -test.benchtime 10s
func BenchmarkMatcher_Match_1000(b *testing.B) {
	benchmarkMatcher_Match(b, 1000)
}*/


var m = matcher.NewMatcher(150, 40, 5)
var mu sync.Mutex

func TestNewMatcher(t *testing.T) {
	go func() {
		for {
			fmt.Println("Matching service is running ...", m.PlayerInQueueCount())
			mu.Lock()
			m.Match(matcher.Time(time.Now().Unix()), 2)
			mu.Unlock()
			time.Sleep(time.Second)
		}
	}()
	go func() {
		for {
			mu.Lock()
			groups := m.GroupsPlayerIds()
			fmt.Println("groups",groups)
			mu.Unlock()
			time.Sleep(time.Second)
		}
	}()
	http.HandleFunc("/join",JoinHandler)
	err := http.ListenAndServe("0.0.0.0:8880",nil)
	if err!=nil{
		panic(err)
	}
	time.Sleep(1*time.Hour)
}

func JoinHandler(w http.ResponseWriter,r *http.Request)  {
	query := r.URL.Query()
	playerId,_ := strconv.ParseInt(query["player_id"][0],10,64)
	duanId,_ := strconv.ParseInt(query["duan_id"][0],10,64)
	matcherTime := matcher.Time(time.Now().Unix())

	mu.Lock()

	//  分数/5必须总分数/5 不然会报错
	err := m.JoinQueue(matcher.PlayerId(strconv.Itoa(int(playerId))), matcherTime, matcher.PlayerScore(duanId))
	if err!=nil{
		fmt.Println(err.Error())
	}
	mu.Unlock()

}