package main

import (
	"bytes"
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
	LastBuff          string `json:"lastBuff"`
	DownContinueFlag  int    `json:"downContinueFlag"`
	VisibleFlag       int    `json:"visibleFlag"`
	ReadFlag          bool   `json:"readFlag"`
	DisplayFlag       int    `json:"displayFlag"`
	BlacklistFlag     int    `json:"blacklistFlag"`
	LikeFlag          int    `json:"likeFlag"`
}

type CommentItem struct {
	LevelTwoComment   []LevelTwoCommentItem `json:"levelTwoComment"`
	CommentId         string                `json:"commentId"`
	CommentNickname   string                `json:"commentNickname"`
	CommentContent    string                `json:"commentContent"`
	CommentHeadurl    string                `json:"commentHeadurl"`
	CommentCreatetime string                `json:"commentCreatetime"`
	CommentLikeCount  int                   `json:"commentLikeCount"`
	LastBuff          string                `json:"lastBuff"`
	DownContinueFlag  int                   `json:"downContinueFlag"`
	VisibleFlag       int                   `json:"visibleFlag"`
	ReadFlag          bool                  `json:"readFlag"`
	DisplayFlag       int                   `json:"displayFlag"`
	BlacklistFlag     int                   `json:"blacklistFlag"`
	LikeFlag          int                   `json:"likeFlag"`
}

func getVideoComments(cookies []*network.Cookie, objectId string, lastBuff string) (*[]CommentItem, *string, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	url := "/micro/interaction/cgi-bin/mmfinderassistant-bin/comment/comment_list"
	payload := `{"lastBuff":"%s","exportId":"%s","commentSelection":false,"forMcn":false,"timestamp":"%s","_log_finder_uin":"","_log_finder_id":"%s","rawKeyBuff":null,"pluginSessionId":null,"scene":7,"reqScene":7}`
	payload = fmt.Sprintf(payload, lastBuff, objectId, timestamp, finderUsername)

	bodyString, err := requestUrl(url, payload, cookies)
	if err != nil {
		return nil, nil, fmt.Errorf("获取视频评论列表失败: %v", err)
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
		return nil, nil, fmt.Errorf("解析 CommentList JSON失败: %v", err)
	}

	return &v.Data.Comment, &v.Data.LastBuff, nil
}

type ReplyCommentItem struct {
	LevelTwoComment   []LevelTwoCommentItem `json:"levelTwoComment"`
	CommentId         string                `json:"commentId"`
	CommentContent    string                `json:"commentContent"`
	CommentCreatetime string                `json:"commentCreatetime"`
	CommentLikeCount  int                   `json:"commentLikeCount"`
	ReplyCommentId    string                `json:"replyCommentId"`
	ReplyContent      string                `json:"replyContent"`
	VisibleFlag       int                   `json:"visibleFlag"`
	DisplayFlag       int                   `json:"displayFlag"`
	Username          string                `json:"username"`
}

type ReplyCommentPayload struct {
	ReplyCommentId  string      `json:"replyCommentId"`
	Content         string      `json:"content"`
	ClientId        string      `json:"clientId"`
	RootCommentId   string      `json:"rootCommentId"`
	Comment         CommentItem `json:"comment"`
	ExportId        string      `json:"exportId"`
	Timestamp       string      `json:"timestamp"`
	Log_finder_uin  string      `json:"_log_finder_uin"`
	Log_finder_id   string      `json:"_log_finder_id"`
	RawKeyBuff      *string     `json:"rawKeyBuff"`
	PluginSessionId *string     `json:"pluginSessionId"`
	Scene           int         `json:"scene"`
	ReqScene        int         `json:"reqScene"`
}

func sanitizeControlChars(s string) string {
	var buf bytes.Buffer
	for _, r := range s {
		// 清理掉非法控制字符
		if r < 0x20 && r != '\n' && r != '\r' && r != '\t' {
			continue
		}
		buf.WriteRune(r)
	}

	return buf.String()
}

func replyComment(cookies []*network.Cookie, objectId string, commentItem CommentItem, replyContent string) (*ReplyCommentItem, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	clientId := uuid.New().String()

	// payload := `{"replyCommentId":"%s","content":"%s","clientId":"%s","rootCommentId":"%s","comment":{"levelTwoComment":[],"commentId":"%s","commentNickname":"%s","commentContent":"%s","commentHeadurl":"%s","commentCreatetime":"%s","commentLikeCount":0,"lastBuff":"","downContinueFlag":0,"visibleFlag":0,"readFlag":true,"displayFlag":0,"blacklistFlag":0,"likeFlag":0},"exportId":"%s","timestamp":"%s","_log_finder_uin":"","_log_finder_id":"%s","rawKeyBuff":null,"pluginSessionId":null,"scene":7,"reqScene":7}`
	// payload = fmt.Sprintf(payload, commentItem.CommentId, replyContent, clientId, commentItem.CommentId, commentItem.CommentId, commentItem.CommentNickname, commentItem.CommentContent, commentItem.CommentHeadurl, commentItem.CommentCreatetime, objectId, timestamp, finderUsername)

	// log.Println("回复评论请求参数: ", payload)

	url := "/micro/interaction/cgi-bin/mmfinderassistant-bin/comment/create_comment"
	payload := ReplyCommentPayload{
		ReplyCommentId:  commentItem.CommentId,
		Content:         sanitizeControlChars(replyContent),
		ClientId:        clientId,
		RootCommentId:   commentItem.CommentId,
		Comment:         commentItem,
		ExportId:        objectId,
		Timestamp:       timestamp,
		Log_finder_uin:  "",
		Log_finder_id:   finderUsername,
		RawKeyBuff:      nil,
		PluginSessionId: nil,
		Scene:           7,
		ReqScene:        7,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("序列化 ReplyCommentPayload 失败: %v", err)
	}

	payloadString := string(payloadBytes)

	bodyString, err := requestUrl(url, payloadString, cookies)
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
