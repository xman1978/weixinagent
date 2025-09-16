package main

import (
	"context"
	"fmt"
	"os"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

func generateReplyContent(videoDescription string, commentContent string) (*string, error) {
	systemPrompt := `你是一个新媒体运营官，专门负责回复用户的评论。 
你的目标：通过你的回复，增加和粉丝的互动，提升视频号的活力值。 
你的语言风格：真诚，礼貌，并且不失幽默。 
你的任务：
1. 根据视频的内容，对用户的评论进行回复。
2. 如果用户评论的内容与视频无关，或者视频描述为空，则从下列回答中随机选择一个合适的回复：
 - 感谢评论[抱拳][抱拳]
 - 感谢支持[抱拳][抱拳]
 - 感谢评论[握手][握手]
 - 感谢支持[握手][握手]
 - 感谢评论[爱心][爱心]
 - 感谢支持[爱心][爱心]
 - 感谢您一直以来的支持[握手]
 - 感谢您如此用心的评论[抱拳]
 - 感谢您每天都来支持我[握手]
 - 感谢您如此用心的评论[爱心]
输出要求：
1. 只能输出回复内容，不要输出其他内容；
2. 你发布的回复内容不能违反国家法律法规，不能违反微信视频号平台的规定；
3. 你发布的回复内容不能超过 20 个字。`

	userPrompt := `视频描述：` + videoDescription + `\n评论内容：` + commentContent

	os.Setenv("ARK_API_KEY", "b8c7ff23-1fae-4e5a-9385-6de7d95dc0b9")
	os.Setenv("ARK_ENDPOINT_ID", "ep-20250916104034-k9prt")

	// 请确保您已将 API Key 存储在环境变量 ARK_API_KEY 中
	// 初始化Ark客户端，从环境变量中读取您的API Key
	client := arkruntime.NewClientWithApiKey(
		// 从环境变量中获取您的 API Key。此为默认方式，您可根据需要进行修改
		os.Getenv("ARK_API_KEY"),
		// 此为默认路径，您可根据业务所在地域进行配置
		arkruntime.WithBaseUrl("https://ark.cn-beijing.volces.com/api/v3"),
	)

	ctx := context.Background()

	temperature, topP, frequencyPenalty, maxTokens := float32(1), float32(0.7), float32(1.2), int(4096)
	req := model.CreateChatCompletionRequest{
		// 您可以前往 开通管理页 开通服务，并在 在线推理页 创建ep后进行查看
		Model: os.Getenv("ARK_ENDPOINT_ID"),
		Messages: []*model.ChatCompletionMessage{
			{
				Role: model.ChatMessageRoleSystem,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String(systemPrompt),
				},
			},
			{
				Role: model.ChatMessageRoleUser,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String(userPrompt),
				},
			},
		},
		Temperature:      &temperature,
		TopP:             &topP,
		MaxTokens:        &maxTokens,
		FrequencyPenalty: &frequencyPenalty,
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("生成回复内容失败: %v", err)
	}

	return resp.Choices[0].Message.Content.StringValue, nil

}
