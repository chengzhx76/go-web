package task

import (
	"encoding/json"
	"fmt"
	"go-web/cache"
	"go-web/mode"
	"io/ioutil"
	"net/http"
	"time"
)

type wxAccessToken struct {
	mode.CommonError
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// https://github.com/silenceper/wechat/blob/master/context/access_token.go
func requestWxAccessToken(globalCache *cache.Cache) {
	accessToken := func() {
		resp, err := http.Get("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=wx52ddb78878fa6d98&secret=44af2777f136af01accabc96bc78d9cc")
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))

		wxAccessToken := &wxAccessToken{}
		if err = json.Unmarshal(body, wxAccessToken); err != nil {
			accessToken := wxAccessToken.AccessToken
			globalCache.Put("accessToken", accessToken)
			//cache1 := &cache.Cache{}
			//cache2 := new(cache.Cache)
		}
	}
	go accessToken()
}

func StartAccessTokenTask(globalCache *cache.Cache) {
	requestWxAccessToken(globalCache)
	task := time.NewTimer(time.Duration(5) * time.Second)

	/*for range task.C {
		requestWxAccessToken(wxCache)
		task.Reset(5 * time.Second)
	}*/

	for {
		select {
		case <-task.C:
			requestWxAccessToken(globalCache)
			task.Reset(5 * time.Second)
		}
	}
}
