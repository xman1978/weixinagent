package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
)

func replyWeixinChannelComment(cookies []*network.Cookie, out chan<- string, done chan<- error) {
	// 确保关闭通道
	defer close(out)

	out <- "获取视频列表（第 1 页）"
	pageSize := 20
	videoList, totalCount, err := getVideoList(cookies, pageSize, 1)
	if err != nil {
		done <- fmt.Errorf("获取视频列表失败 ## %v", err)
		close(done)
		return
	}
	Pages := totalCount / pageSize
	if totalCount%pageSize != 0 {
		Pages++
	}

	replyCount := 0
	for currentPage := 1; currentPage <= Pages; currentPage++ {
		for _, video := range *videoList {
			out <- fmt.Sprintf("视频: %s %s", video.ShortTitle, video.Description)
			// fmt.Println("视频: ", video.ShortTitle, video.Description, video.ObjectId)
			// fmt.Println("获取视频评论：")
			comments, err := getVideoComments(cookies, video.ObjectId)
			if err != nil {
				done <- fmt.Errorf("获取视频评论失败 ## %v", err)
				close(done)
				return
			}
			// fmt.Println("评论列表: ", *comments)
			out <- "--------------------------------"

			for _, comment := range *comments {
				if len(comment.LevelTwoComment) > 0 || strings.TrimSpace(comment.CommentContent) == "" || strings.TrimSpace(comment.CommentNickname) == nickName {
					continue
				}

				out <- fmt.Sprintf("评论内容: %s", comment.CommentContent)
				replyContent, err := generateReplyContentV2(comment.CommentContent)
				if err != nil {
					done <- fmt.Errorf("生成回复内容失败 ## %v", err)
					close(done)
					return
				}
				out <- fmt.Sprintf("AI 生成的回复内容: %s", *replyContent)

				replyComment, err := replyComment(cookies, video.ObjectId, comment, *replyContent)
				if err != nil {
					out <- fmt.Sprintf("回复评论失败 ## %v", err)
					continue
				}
				replyCount++
				out <- fmt.Sprintf("已回复，回复内容: %s", replyComment.CommentContent)
				out <- "=========================================="

				// 在 30 到 90 秒之间随机等待
				out <- "等待 30 到 60 秒之间随机等待..."
				time.Sleep(time.Duration(rand.Intn(30)+30) * time.Second)
			}
		}

		out <- fmt.Sprintf("视频总数：%d，已成功回复评论数: %d", totalCount, replyCount)

		if currentPage == Pages {
			break
		}

		out <- fmt.Sprintf("获取视频列表（第 %d 页）", currentPage+1)
		videoList, _, err = getVideoList(cookies, pageSize, currentPage+1)
		if err != nil {
			done <- fmt.Errorf("获取视频列表失败 ## %v", err)
			close(done)
			return
		}
	}

	done <- nil
	close(done)
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

	msgs := make(chan string, 40)
	done := make(chan error, 1)
	go replyWeixinChannelComment(cookies, msgs, done)

	for msg := range msgs {
		fmt.Println(msg)
	}

	if err := <-done; err != nil {
		log.Fatal("回复评论失败 ## ", err)
	}
}
