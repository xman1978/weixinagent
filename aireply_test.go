package main

import (
	"fmt"
	"testing"
)

func TestGenerateReplyContentV2(t *testing.T) {
	commentContent := "你太过敏感，只顾自己的感受"
	replyContent, err := generateReplyContentV2(commentContent)
	if err != nil {
		t.Errorf("生成回复内容失败: %v", err)
	}
	fmt.Println("生成回复内容: ", *replyContent)
}
