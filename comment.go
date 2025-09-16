package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/google/uuid"
)

type LevelTwoCommentItem struct {
	CommentId         string `json:"commentId"`
	CommentNickname   string `json:"commentNickname"`
	CommentContent    string `json:"commentContent"`
	CommentHeadurl    string `json:"commentHeadurl"`
	CommentCreatetime string `json:"commentCreatetime"`
	ReplyCommentId    string `json:"replyCommentId"`
	ReplyContent      string `json:"replyContent"`
}

type CommentItem struct {
	LevelTwoComment   []LevelTwoCommentItem `json:"levelTwoComment"`
	CommentId         string                `json:"commentId"`
	CommentNickname   string                `json:"commentNickname"`
	CommentContent    string                `json:"commentContent"`
	CommentHeadurl    string                `json:"commentHeadurl"`
	CommentCreatetime string                `json:"commentCreatetime"`
}

func getVideoComments(cookies []*network.Cookie, objectId string) (*[]CommentItem, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	url := "/micro/interaction/cgi-bin/mmfinderassistant-bin/comment/comment_list"
	payload := `{"lastBuff":"","exportId":"%s","commentSelection":false,"forMcn":false,"timestamp":"%s","_log_finder_uin":"","_log_finder_id":"%s","rawKeyBuff":null,"pluginSessionId":null,"scene":7,"reqScene":7}`
	payload = fmt.Sprintf(payload, objectId, timestamp, finderUsername)

	bodyString, err := requestUrl(url, payload, cookies)
	if err != nil {
		return nil, fmt.Errorf("获取视频评论列表失败: %v", err)
	}

	var v struct {
		Data struct {
			Comment          []CommentItem `json:"comment"`
			LastBuff         string        `json:"lastBuff"`
			CommentCount     int           `json:"commentCount"`
			DownContinueFlag int           `json:"downContinueFlag"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(*bodyString), &v); err != nil {
		return nil, fmt.Errorf("解析 CommentList JSON失败: %v", err)
	}

	return &v.Data.Comment, nil
}

type ReplyCommentItem struct {
	CommentId      string `json:"commentId"`
	CommentContent string `json:"commentContent"`
	ReplyCommentId string `json:"replyCommentId"`
	ReplyContent   string `json:"replyContent"`
}

func replyComment(cookies []*network.Cookie, objectId string, commentItem CommentItem, replyContent string) (*ReplyCommentItem, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	clientId := uuid.New().String()

	url := "/micro/interaction/cgi-bin/mmfinderassistant-bin/comment/create_comment"
	payload := `{"replyCommentId":"%s","content":"%s","clientId":"%s","rootCommentId":"%s","comment":{"levelTwoComment":[],"commentId":"%s","commentNickname":"%s","commentContent":"%s","commentHeadurl":"%s","commentCreatetime":"%s","commentLikeCount":0,"lastBuff":"","downContinueFlag":0,"visibleFlag":0,"readFlag":true,"displayFlag":0,"blacklistFlag":0,"likeFlag":0},"exportId":"%s","timestamp":"%s","_log_finder_uin":"","_log_finder_id":"%s","rawKeyBuff":null,"pluginSessionId":null,"scene":7,"reqScene":7}`
	payload = fmt.Sprintf(payload, commentItem.CommentId, replyContent, clientId, commentItem.CommentId, commentItem.CommentId, commentItem.CommentNickname, commentItem.CommentContent, commentItem.CommentHeadurl, commentItem.CommentCreatetime, objectId, timestamp, finderUsername)

	// log.Println("回复评论请求参数: ", payload)

	bodyString, err := requestUrl(url, payload, cookies)
	if err != nil {
		return nil, fmt.Errorf("回复评论失败: %v", err)
	}

	// log.Println("回复评论响应: ", *bodyString)

	var v struct {
		Data struct {
			Comment ReplyCommentItem `json:"comment"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(*bodyString), &v); err != nil {
		return nil, fmt.Errorf("解析 ReplyCommentItem JSON失败: %v", err)
	}

	return &v.Data.Comment, nil
}
