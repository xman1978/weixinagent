package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Prompt PromptConfig `yaml:"prompt"`
	Llmapi LlmapiConfig `yaml:"llmapi"`
}

type PromptConfig struct {
	SystemPrompt string `yaml:"systemPrompt"`
	UserPrompt   string `yaml:"userPrompt"`
}

type LlmapiConfig struct {
	ModelName           string  `yaml:"modelName"`
	APIKey              string  `yaml:"apiKey"`
	BaseURL             string  `yaml:"baseURL"`
	FrequencyPenalty    float64 `yaml:"frequencyPenalty"`
	TopP                float64 `yaml:"topP"`
	Temperature         float64 `yaml:"temperature"`
	MaxCompletionTokens int64   `yaml:"maxCompletionTokens"`
}

func generateReplyContentV2(commentContent string) (*string, error) {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		panic(fmt.Errorf("无法读取配置文件: %w", err))
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(fmt.Errorf("YAML 解析失败: %w", err))
	}

	systemPrompt := config.Prompt.SystemPrompt

	userPrompt := config.Prompt.UserPrompt + commentContent

	params := OpenAIParams{
		ModeName:            config.Llmapi.ModelName,
		APIKey:              config.Llmapi.APIKey,
		BaseURL:             config.Llmapi.BaseURL,
		FrequencyPenalty:    config.Llmapi.FrequencyPenalty,
		TopP:                config.Llmapi.TopP,
		MaxCompletionTokens: config.Llmapi.MaxCompletionTokens,
		Temperature:         config.Llmapi.Temperature,
	}

	reply, err := openaiReply(&params, systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("生成回复内容失败: %v", err)
	}

	return reply, nil
}
