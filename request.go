package main

import (
	"compress/flate"
	"compress/gzip"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/chromedp/cdproto/network"
)

var baseUrl string = "https://channels.weixin.qq.com"

const fingerPrintDeviceId string = "b3d80e79cd3c6f230e92fe539c18f3a5"

var finderUsername string = ""
var wxuuid string = "0000000000"

func cookiesToHeader(cookies []*network.Cookie) string {
	var sb []string
	for _, c := range cookies {
		sb = append(sb, fmt.Sprintf("%s=%s", c.Name, c.Value))
	}
	return strings.Join(sb, "; ")
}

func generateRidtring() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	part1 := hex.EncodeToString(bytes)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	part2 := hex.EncodeToString(bytes)

	return fmt.Sprintf("%s-%s", part1, part2), nil
}

func requestUrl(requestUrl string, payload string, cookies []*network.Cookie) (*string, error) {
	parsedUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("解析URL失败: %v", err)
	}

	rid, err := generateRidtring()
	if err != nil {
		return nil, fmt.Errorf("生成rid失败: %v", err)
	}

	uri := baseUrl + requestUrl + "?_aid=e8bbc517-c9aa-4788-acba-0e7459085a7e&_pageUrl=https:%2F%2Fchannels.weixin.qq.com%2Fmicro%2Finteraction%2Fcomment&_rid=" + rid
	req, err := http.NewRequest("POST", uri, strings.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("请求URL失败: %v", err)
	}
	req.Header.Set("Cookie", cookiesToHeader(cookies))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", parsedUrl.Hostname())
	req.Header.Set("Origin", baseUrl)
	req.Header.Set("finger-print-device-id", fingerPrintDeviceId)
	req.Header.Set("X-WECHAT-UIN", wxuuid)

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求URL失败: %v", err)
	}
	defer response.Body.Close()

	// 读取响应体数据
	var reader io.Reader = response.Body

	// 检查是否使用压缩算法
	if response.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(response.Body)
		if err != nil {
			return nil, fmt.Errorf("创建 gzip 解压器失败: %v", err)
		}
		defer gzipReader.Close()
		reader = gzipReader
	} else if response.Header.Get("Content-Encoding") == "deflate" {
		deflateReader := flate.NewReader(response.Body)
		defer deflateReader.Close()
		reader = deflateReader
	} else if response.Header.Get("Content-Encoding") == "br" {
		brReader := brotli.NewReader(response.Body)
		reader = brReader
	} else {
		reader = response.Body
	}

	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal("读取响应体失败 ## ", err)
	}
	bodyString := string(bodyBytes)

	// fmt.Println("bodyString: ", bodyString)

	// 统一处理后端错误: 如果形如 {"errCode":-1, "errMsg":"..."}
	var apiErr struct {
		Error   *int   `json:"errCode"`
		Message string `json:"errMsg"`
	}
	if err := json.Unmarshal(bodyBytes, &apiErr); err == nil {
		if apiErr.Error != nil && *apiErr.Error != 0 {
			return nil, fmt.Errorf("接口错误: %d, %s", *apiErr.Error, apiErr.Message)
		}
	}

	return &bodyString, nil
}
