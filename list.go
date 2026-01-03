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

func getVideoList(cookies []*network.Cookie, pageSize int, currentPage int) (*[]VideoList, int, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	url := "/micro/content/cgi-bin/mmfinderassistant-bin/post/post_list"
	payload := `{"pageSize":%d,"currentPage":%d,"onlyUnread":false,"needAllCommentCount":true,"userpageType":13,"forMcn":false,"lastBuff":"","timestamp":"%s","_log_finder_uin":"%s","_log_finder_id":"","rawKeyBuff":null,"pluginSessionId":null,"scene":7,"reqScene":7}`
	payload = fmt.Sprintf(payload, pageSize, currentPage, timestamp, finderUsername)

	bodyString, err := requestUrl(url, payload, cookies)
	if err != nil {
		return nil, 0, fmt.Errorf("获取视频列表失败: %v", err)
	}

	var v struct {
		Data struct {
			TotalCount int `json:"totalCount"`
			List       []struct {
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
		return nil, 0, fmt.Errorf("解析 VideoList JSON失败: %v", err)
	}

	// fmt.Printf("获取视频列表: %v\n", v.Data.List)

	var videoList []VideoList
	for _, list := range v.Data.List {
		// 安全获取 CoverUrl
		coverUrl := ""
		if len(list.Desc.Media) > 0 {
			coverUrl = list.Desc.Media[0].CoverUrl
		}

		// 安全获取 ShortTitle
		shortTitle := "未知标题"
		if len(list.Desc.ShortTitle) > 0 {
			shortTitle = list.Desc.ShortTitle[0].Title
		}

		// 安全获取 Description（字符串类型，如果为空则已经是空字符串）
		description := list.Desc.Description

		// 安全获取 CreateTime（检查是否为 0 或无效值）
		createTime := "2000-01-01 00:00:00"
		if list.CreateTime > 0 {
			createTime = time.Unix(list.CreateTime, 0).Format("2006-01-02 15:04:05")
		}

		videoList = append(videoList, VideoList{
			ObjectId:    list.ObjectId,
			CreateTime:  createTime,
			CoverUrl:    coverUrl,
			ShortTitle:  shortTitle,
			Description: description,
		})
	}

	return &videoList, v.Data.TotalCount, nil
}
