package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

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

	fmt.Println("获取视频列表：")
	videoList, err := getVideoList(cookies)
	if err != nil {
		log.Fatal("获取视频列表失败 ## ", err)
	}

	error_times := 0
	for _, video := range *videoList {
		fmt.Println("视频: ", video.ShortTitle, video.Description)
		// fmt.Println("视频: ", video.ShortTitle, video.Description, video.ObjectId)
		// fmt.Println("获取视频评论：")
		comments, err := getVideoComments(cookies, video.ObjectId)
		if err != nil {
			log.Fatal("获取视频评论失败 ## ", err)
		}
		// fmt.Println("评论列表: ", *comments)
		fmt.Println("--------------------------------")

		for _, comment := range *comments {
			if len(comment.LevelTwoComment) > 0 || strings.TrimSpace(comment.CommentContent) == "" || strings.TrimSpace(comment.CommentNickname) == nickName {
				continue
			}

			fmt.Println("评论内容: ", comment.CommentContent)
			replyContent, err := generateReplyContent(video.Description, comment.CommentContent)
			if err != nil {
				log.Fatal("生成回复内容失败 ## ", err)
			}
			fmt.Println("AI 生成的回复内容: ", *replyContent)

			replyComment, err := replyComment(cookies, video.ObjectId, comment, *replyContent)
			if err != nil {
				log.Println("回复评论失败 ## ", err)
				error_times++
				if error_times > 3 {
					log.Fatal("回复评论失败次数超过3次，退出程序")
				}
				continue
			}
			fmt.Println("已回复，回复内容: ", replyComment.CommentContent)
			fmt.Println("==========================================")

			// 在 60 到 180 秒之间随机等待
			fmt.Println("等待 60 到 180 秒之间随机等待...")
			time.Sleep(time.Duration(rand.Intn(120)+60) * time.Second)
		}
	}

}
