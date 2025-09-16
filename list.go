package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/network"
)

type VideoList struct {
	ObjectId    string
	CreateTime  string
	CoverUrl    string
	ShortTitle  string
	Description string
}

func getVideoList(cookies []*network.Cookie) (*[]VideoList, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	url := "/micro/interaction/cgi-bin/mmfinderassistant-bin/post/post_list"
	payload := `{"pageSize":50,"onlyUnread":false,"needAllCommentCount":true,"userpageType":13,"forMcn":false,"lastBuff":"","timestamp":"%s","_log_finder_uin":"%s","_log_finder_id":"","rawKeyBuff":null,"pluginSessionId":null,"scene":7,"reqScene":7}`
	payload = fmt.Sprintf(payload, timestamp, finderUsername)

	bodyString, err := requestUrl(url, payload, cookies)
	if err != nil {
		return nil, fmt.Errorf("获取视频列表失败: %v", err)
	}

	var v struct {
		Data struct {
			List []struct {
				ObjectId   string `json:"objectId"`
				CreateTime int64  `json:"createTime"`
				Desc       struct {
					Media []struct {
						CoverUrl string `json:"coverUrl"`
					} `json:"media"`
					ShortTitle []struct {
						Title string `json:"shortTitle"`
					} `json:"shortTitle"`
					Description string `json:"description"`
				} `json:"desc"`
			} `json:"list"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(*bodyString), &v); err != nil {
		return nil, fmt.Errorf("解析 VideoList JSON失败: %v", err)
	}

	var videoList []VideoList
	for _, list := range v.Data.List {
		videoList = append(videoList, VideoList{
			ObjectId:    list.ObjectId,
			CreateTime:  time.Unix(list.CreateTime, 0).Format("2006-01-02 15:04:05"),
			CoverUrl:    list.Desc.Media[0].CoverUrl,
			ShortTitle:  list.Desc.ShortTitle[0].Title,
			Description: list.Desc.Description,
		})
	}

	return &videoList, nil
}
