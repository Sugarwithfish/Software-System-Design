package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sugarwithfish/Software-System-Design/models/request"
	"github.com/Sugarwithfish/Software-System-Design/models/response"
	"log"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const (
	SysMsgOfChoice = "你是一名卓越的程序员，精通 Go 和 Javascript，并且擅长出一些相关的选择题\n" +
		"# 要求\n" +
		"1. 题目形式为单选题或多选题，且每道题目都只能包含四个选项\n" +
		"2. 每次生成指定数量题目，答案是随机分布的，以 JSON 数组的格式返回，示例如下：\n" +
		"   [{\"title\":\"Go 语言的优势有哪些\",\"options\":[\"具有强大的并发支持\",\"语法简洁\",\"不适合开发 Web 应用程序\",\"编译速度慢\"],\"answer\":[0,1]}]\n" +
		"   注意：不管单选还是多选，正确答案\"answer\"都用数组返回，下标从0开始\n" +
		"3. 除此之外，不要输出任何内容，必须是 JSON数组格式，以 `[` 开头。\n" +
		"4. 接下来我将提供要生成的题目类型、编程语言和题目数量"

	SysMsgOfCoding = "你是一名卓越的程序员，精通 Go 和Javascript，并且擅长出一些相关的编程题\n" +
		"# 要求\n" +
		"1. 根据要求出指定数量的编程题，以 JSON 数组格式返回，示例如下：\n" +
		"   [{\"title\":\"两数之和\",\"content\":\"给定指定的两个数，返回其相加之和\"}]\n" +
		"2. 除此之外，不要输出任何内容，必须是 JSON 数组格式，以 `[` 开头。\n" +
		"3. 接下来我将提供要生成题目的编程语言和题目数量"
)

type LLMConfig struct {
	BaseURL string `json:"aibase_url"`
	APIKey  string `json:"ai_api_key"`
	Model   string `json:"model"`
}

func Dochat(config *LLMConfig, systemMessage string, userMessage string) (string, error) {
	client := openai.NewClient(
		option.WithAPIKey(config.APIKey),
		option.WithBaseURL(config.BaseURL),
	)

	chatCompletion, err := client.Chat.Completions.New(context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(systemMessage),
				openai.UserMessage(userMessage),
			},
			Model: openai.ChatModel(config.Model),
		})
	if err != nil {
		return "", err
	}

	if len(chatCompletion.Choices) == 0 {
		return "", errors.New("大模型响应出错")
	}

	return chatCompletion.Choices[0].Message.Content, nil
}

func handlePrompt(params request.AIGenQuestionParams) (string, string, bool) {
	var systemMessage, userMessage string
	var isCoding = false

	if params.Type == 3 { // 编程题
		isCoding = true
		systemMessage = SysMsgOfCoding
		userMessage = fmt.Sprintf("编程语言: %s，题目数量：%d", params.Language, params.Amount)
	} else { // 选择题
		var questionTypeStr string
		switch params.Type {
		case 1:
			questionTypeStr = "单选题"
		case 2:
			questionTypeStr = "多选题"
		}
		systemMessage = SysMsgOfChoice
		userMessage = fmt.Sprintf("题目类型：%s，编程语言：%s，题目数量：%d", questionTypeStr, params.Language, params.Amount)
	}

	return systemMessage, userMessage, isCoding
}

func AIGenerateQuestion(params request.AIGenQuestionParams, cfg *LLMConfig) ([]*response.AIGenQuestion, error) {
	systemMessage, userMessage, isCoding := handlePrompt(params)
	chatResponse, err := Dochat(cfg, systemMessage, userMessage)
	if err != nil {
		return nil, fmt.Errorf("大模型服务调用失败: %w", err)
	}

	var questions []*response.AIGenQuestion
	if isCoding {
		var generatedQuestions []response.CodingQuestion
		if err := json.Unmarshal([]byte(chatResponse), &generatedQuestions); err != nil {
			return nil, fmt.Errorf("大模型响应格式有误: %w", err)
		}
		for _, genQuestion := range generatedQuestions {
			questions = append(questions, &response.AIGenQuestion{
				Title:   genQuestion.Title,
				Content: genQuestion.Content,
			})
		}
	} else {
		var generatedQuestions []response.ChoiceQuestion
		if err := json.Unmarshal([]byte(chatResponse), &generatedQuestions); err != nil {
			return nil, fmt.Errorf("大模型响应格式有误: %w", err)
		}
		for _, genQuestion := range generatedQuestions {
			optionsStr, _ := json.Marshal(genQuestion.Options)
			answerStr, _ := json.Marshal(genQuestion.Answer)
			questions = append(questions, &response.AIGenQuestion{
				Title:   genQuestion.Title,
				Options: string(optionsStr),
				Answer:  string(answerStr),
			})
		}
	}

	if len(questions) == 0 {
		return nil, fmt.Errorf("大模型未能生成任何题目")
	}

	log.Println("大模型生成题目成功！")

	return questions, nil
}
