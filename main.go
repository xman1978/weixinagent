package main

import (
	"fmt"
	"log"

	"github.com/chromedp/cdproto/network"
)

func replyWeixinChannelComment(cookies []*network.Cookie) {
	log.Printf("获取视频列表（第 1 页）")
	pageSize := 20
	videoList, totalCount, err := getVideoList(cookies, pageSize, 1)
	if err != nil {
		log.Fatalf("获取视频列表失败 ## %v", err)
	}
	Pages := totalCount / pageSize
	if totalCount%pageSize != 0 {
		Pages++
	}

	replyCount := 0
	for currentPage := 1; currentPage <= Pages; currentPage++ {
		for _, video := range *videoList {
			log.Printf("视频: %s %s", video.ShortTitle, video.Description)
			log.Println("==============================================================")
			// log.Println("视频: ", video.ShortTitle, video.Description, video.ObjectId)
			// log.Println("获取视频评论：")
			lastBuff := ""
			for {
				comments, lastBuffPtr, err := getVideoComments(cookies, video.ObjectId, lastBuff)
				if err != nil {
					log.Printf("获取视频评论失败 ## %v", err)
					break
				}

				// 如果评论列表为空，则跳出循环，否则继续获取下一个视频的评论
				if len(*comments) == 0 {
					break
				}

				// fmt.Println("评论列表: ", *comments)
				for _, comment := range *comments {
					_, err := replyComment(cookies, video.ObjectId, comment)
					if err != nil {
						log.Printf("回复评论失败 ## %v", err)
						continue
					}
					replyCount++
				}

				// 如果 lastBuff 为空，则跳出循环，否则继续获取下一页评论
				if *lastBuffPtr == "" {
					break
				}
				lastBuff = *lastBuffPtr
			}
		}

		log.Printf("视频总数：%d，已成功回复评论数: %d", totalCount, replyCount)

		if currentPage == Pages {
			break
		}

		log.Printf("获取视频列表（第 %d 页）", currentPage+1)
		videoList, _, err = getVideoList(cookies, pageSize, currentPage+1)
		if err != nil {
			log.Fatalf("获取视频列表失败 ## %v", err)
		}
	}

}

func main() {
	cookies, err := loginWeixin()
	if err != nil {
		log.Fatal("登录失败 ## ", err)
	}

	finder, err := getLogfinderid(cookies)
	if err != nil {
		log.Fatal("获取LogfinderId失败 ## ", err)
	}
	finderUsername = *finder
	fmt.Println("finderUsername: ", finderUsername)

	uuid, err := getWxuuid(cookies)
	if err != nil {
		log.Fatal("获取Wxuuid失败 ## ", err)
	}
	wxuuid = *uuid
	fmt.Println("wxuuid: ", wxuuid)

	replyWeixinChannelComment(cookies)
}
