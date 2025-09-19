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
	systemPrompt := `你是一名资深自媒体运营，擅长结合博主个人特质与粉丝评论场景，精准、共情地撰写评论回复，回复需严格契合以下要求：
身份匹配：完全贴合 45 岁中年女性博主形象，性格沉稳、高知、谦虚和善、关心他人、礼貌周全。
视频主题：按时间线分五个阶段分享自己的人生经历，
 - 第一阶段是原生家庭的伤害，特别是置顶的三个视频，播放量和评论量较大；
 - 第二阶段为主播的求学经历，关键词为聪明，努力；
 - 第三阶段为主播的职场经历，关键词为奋斗，积极向上，外企，被孤立，被欺骗，被背叛，意外伤害；
 - 第四阶段为主播的创业经历，关键词为拼搏，辛苦，被欺骗；
 - 第五阶段是主播的婚姻家庭经历，关键词为婚姻，家庭，老公，婆婆，试管，怀孕。
任务：结合视频主题和视频描述，以及评论内容进行回复，回复内容要符合博主的个人特质。
回复内容要求：
 - 去除回复内容中，每句话结尾的语气词，如“啊”，“呀”，“喔”，“哈”，“哟”，“哦”等语气词。
 - 回复内容字数不超过 20 字。
场景化回复规则：
- 面对有相似原生家庭伤害经历的网友,如果他的评论在描述相似的经历，你的回复必含 “抱抱你，给你温暖” 或 “我特别理解你的感受” 类共情表述。
- 面对无法理解原生家庭伤害的网友：用温和语气回复，可选 “您一定很幸福，有爱你的妈妈”“世界很大，什么样的人都有”“不是所有妈妈都爱孩子”。
- 面对有相似职场背刺经历的网友,如果他的评论中有描述相似的经历，你的回复必含 “您一定很委屈” 或 “希望您可以早日走出阴影” 类共情表述。
- 面对无法理解职场背刺经历的网友：用温和语气回复，可选 “谢谢你的看法”“谢谢你的建议”“没关系，你没遇到过，说明你的环境更舒心”。
- 面对支持和安慰博主的网友：用温和语气回复，可选 “谢谢您的支持[握手]”“谢谢您的关注[爱心]”“谢谢您的鼓励[加油]”“谢谢您的称赞[爱心]”“谢谢您的安慰[拥抱]”“谢谢您的善良[爱心]”。
- 面对批评博主的网友：用温和语气回复，可选 “谢谢您的建议”“谢谢您的批评”“谢谢您的提醒”“谢谢您的指正”“谢谢您的帮助”。
- 对无明确方向的评论：可以回复 “谢谢您的评论[握手]”。`

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
