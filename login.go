package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

var nickName string = ""

func loginWeixin() ([]*network.Cookie, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("start-maximized", true),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	loginURL := "https://channels.weixin.qq.com/login.html"
	if err := chromedp.Run(ctx,
		chromedp.Navigate(loginURL),
	); err != nil {
		return nil, fmt.Errorf("打开登录页失败: %v", err)
	}

	if err := chromedp.Run(ctx,
		// 等待 <h2 data-v-686a94f1="" class="finder-nickname">xxxx</h2> 可见
		chromedp.WaitVisible(`.finder-nickname`, chromedp.ByQuery),
		// 读取 <h2 data-v-686a94f1="" class="finder-nickname">xxxx</h2> 中的 xxxx
		chromedp.Text(`.finder-nickname`, &nickName, chromedp.ByQuery),
	); err != nil {
		return nil, fmt.Errorf("等待登录失败: %v", err)
	}

	log.Println("登录成功，用户昵称:", nickName)

	// 获取 cookies
	var cookies []*network.Cookie
	if err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			cookies, err = network.GetCookies().Do(ctx)
			return err
		}),
	); err != nil {
		return nil, fmt.Errorf("获取 cookies 失败: %v", err)
	}

	return cookies, nil
}

func getLogfinderid(cookies []*network.Cookie) (*string, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	url := "/cgi-bin/mmfinderassistant-bin/auth/auth_data"
	payload := `{"timestamp":"%s","_log_finder_uin":"","_log_finder_id":"","rawKeyBuff":null,"pluginSessionId":null,"scene":7,"reqScene":7}`
	payload = fmt.Sprintf(payload, timestamp)

	bodyString, err := requestUrl(url, payload, cookies)
	if err != nil {
		return nil, fmt.Errorf("请求 logfinderid url 失败: %v", err)
	}

	var v struct {
		Data struct {
			FinderUser struct {
				FinderUsername string `json:"finderUsername"`
			} `json:"finderUser"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(*bodyString), &v); err != nil {
		return nil, fmt.Errorf("解析 LogfinderId JSON失败: %v", err)
	}

	return &v.Data.FinderUser.FinderUsername, nil
}

func getWxuuid(cookies []*network.Cookie) (*string, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	url := "/cgi-bin/mmfinderassistant-bin/helper/helper_upload_params"
	payload := `{"timestamp":"%s","_log_finder_uin":"","_log_finder_id":"","rawKeyBuff":null,"pluginSessionId":null,"scene":7,"reqScene":7}`
	payload = fmt.Sprintf(payload, timestamp)

	bodyString, err := requestUrl(url, payload, cookies)
	if err != nil {
		return nil, fmt.Errorf("请求 wxuuid url 失败: %v", err)
	}

	var v struct {
		Data struct {
			Wxuuid int64 `json:"uin"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(*bodyString), &v); err != nil {
		return nil, fmt.Errorf("解析 Wxuuid JSON失败: %v", err)
	}

	wxuuid = strconv.FormatInt(v.Data.Wxuuid, 10)

	return &wxuuid, nil
}
