package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sugarwithfish/Software-System-Design/models/request"
	"github.com/Sugarwithfish/Software-System-Design/models/response"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"log"
	"sync"
)

//const (
//	SysMsgOfChoice = "你是一名卓越的程序员，精通 Go 和 Javascript，并且擅长出一些相关的选择题\n" +
//		"# 要求\n" +
//		"1. 题目形式为单选题或多选题，且每道题目都只能包含四个选项\n" +
//		"2. 每次生成指定数量题目，答案是随机分布的，以 JSON 数组的格式返回，示例如下：\n" +
//		"   [{\"title\":\"Go 语言的优势有哪些\",\"options\":[\"具有强大的并发支持\",\"语法简洁\",\"不适合开发 Web 应用程序\",\"编译速度慢\"],\"answer\":[0,1]}]\n" +
//		"   注意：不管单选还是多选，正确答案\"answer\"都用数组返回，下标从0开始\n" +
//		"3. 除此之外，不要输出任何内容，必须是 JSON数组格式，以 `[` 开头。\n" +
//		"4. 接下来我将提供要生成的题目类型、编程语言和题目数量"
//
//	SysMsgOfCoding = "你是一名卓越的程序员，精通 Go 和Javascript，并且擅长出一些相关的编程题\n" +
//		"# 要求\n" +
//		"1. 根据要求出指定数量的编程题，以 JSON 数组格式返回，示例如下：\n" +
//		"   [{\"title\":\"两数之和\",\"content\":\"给定指定的两个数，返回其相加之和\"}]\n" +
//		"2. 除此之外，不要输出任何内容，必须是 JSON 数组格式，以 `[` 开头。\n" +
//		"3. 接下来我将提供要生成题目的编程语言和题目数量"
//)
const (
	SysMsgOfChoice = "你是一名卓越的程序员，精通 Go 和 Javascript，并且擅长出相关的选择题。\n" +
		"# 要求\n" +
		"1. 题目形式为单选题或多选题，且每道题目都只能包含四个选项。\n" +
		"2. 只生成一道题目，答案随机分布，以 JSON 格式返回，示例如下：\\n" +
		"{\"title\":\"Go 语言的优势有哪些\",\"options\":[\"具有强大的并发支持\",\"语法简洁\",\"不适合开发 Web 应用程序\",\"编译速度慢\"],\"answer\":[0,1]}\n" +
		"3. 除此之外，不要输出任何内容，必须是 JSON 格式，以 { 开头。\n" +
		"4. 题目类型：%s，编程语言：%s"

	SysMsgOfCoding = "你是一名卓越的程序员，精通 Go 和Javascript，并且擅长出一些相关的编程题\n" +
		"# 要求\n" +
		"1. 只生成一道编程题，以 JSON 格式返回，示例如下：\n" +
		"{\"title\":\"两数之和\",\"content\":\"给定指定的两个数，返回其相加之和\"}\n" +
		"2. 除此之外，不要输出任何内容，必须是 JSON 格式，以 `{` 开头。\n" +
		"3. 编程语言：%s"
)

type LLMConfig struct {
	BaseURL string `json:"aibase_url"`
	APIKey  string `json:"ai_api_key"`
	Model   string `json:"model"`
}

var cfg *LLMConfig

func InitLLM(key, url, model string) {
	cfg = &LLMConfig{
		APIKey:  key,
		BaseURL: url,
		Model:   model,
	}
}

func Dochat(systemMessage string, userMessage string) (string, error) {
	client := openai.NewClient(
		option.WithAPIKey(cfg.APIKey),
		option.WithBaseURL(cfg.BaseURL),
	)

	chatCompletion, err := client.Chat.Completions.New(context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(systemMessage),
				openai.UserMessage(userMessage),
			},
			Model: openai.ChatModel(cfg.Model),
		})
	if err != nil {
		return "", err
	}

	if len(chatCompletion.Choices) == 0 {
		return "", errors.New("大模型响应出错")
	}

	return chatCompletion.Choices[0].Message.Content, nil
}

// handlePromptSingle 根据请求参数生成系统消息和用户消息（每次只生成一道题）
func handlePromptSingle(params request.AIGenQuestionParams) (string, string, bool) {
	var systemMessage, userMessage string
	var isCoding = false

	if params.Type == 3 { // 编程题
		isCoding = true
		systemMessage = SysMsgOfCoding
		userMessage = fmt.Sprintf("编程语言：%s", params.Language)
	} else { // 选择题
		var questionTypeStr string
		switch params.Type {
		case 1:
			questionTypeStr = "单选题"
		case 2:
			questionTypeStr = "多选题"
		}
		systemMessage = SysMsgOfChoice
		userMessage = fmt.Sprintf("题目类型：%s，编程语言：%s", questionTypeStr, params.Language)
	}
	return systemMessage, userMessage, isCoding
}

// AIGenerateQuestion 处理AI生成题目的请求（并发，每次生成一道题）
func AIGenerateQuestion(params request.AIGenQuestionParams) ([]*response.AIGenQuestion, error) {
	var (
		wg             sync.WaitGroup
		mu             sync.Mutex
		results        []*response.AIGenQuestion
		errs           []error
		maxConcurrency = 10 // 可根据需要调整
		sem            = make(chan struct{}, maxConcurrency)
	)

	for i := 0; i < params.Amount; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			// 生成单题 prompt
			systemMessage, userMessage, isCoding := handlePromptSingle(params)
			chatResponse, err := Dochat(systemMessage, userMessage)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}
			var question *response.AIGenQuestion
			if isCoding {
				var genQuestion response.CodingQuestion
				if err := json.Unmarshal([]byte(chatResponse), &genQuestion); err != nil {
					mu.Lock()
					errs = append(errs, err)
					mu.Unlock()
					return
				}
				question = &response.AIGenQuestion{
					Title:   genQuestion.Title,
					Content: genQuestion.Content,
				}
			} else {
				var genQuestion response.ChoiceQuestion
				if err := json.Unmarshal([]byte(chatResponse), &genQuestion); err != nil {
					mu.Lock()
					errs = append(errs, err)
					mu.Unlock()
					return
				}
				optionsStr, _ := json.Marshal(genQuestion.Options)
				answerStr, _ := json.Marshal(genQuestion.Answer)
				question = &response.AIGenQuestion{
					Title:   genQuestion.Title,
					Options: string(optionsStr),
					Answer:  string(answerStr),
				}
			}
			mu.Lock()
			results = append(results, question)
			mu.Unlock()
		}()
	}
	wg.Wait()
	if len(results) == 0 {
		return nil, fmt.Errorf("全部生成失败: %v", errs)
	}

	log.Println("大模型生成题目成功！")
	return results, nil
}
