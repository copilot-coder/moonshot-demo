package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/copilot-coder/moonshot-demo/chatengine"
)

func main() {
	cfg := chatengine.Config{
		ApiKey: os.Getenv("MOONSHOT_API_KEY"),
	}
	engine := chatengine.NewEngine(cfg)

	req := chatengine.ChatReq{
		Model: "moonshot-v1-8k",
		Messages: []chatengine.Message{
			{Role: "user", Content: "北京有什么好玩的景点？用中文回答"},
		},
		MaxTokens: 1000,
		Stream:    true,
	}

	// stream request
	fmt.Println(">>> test for stream request")
	ch, err := engine.StreamRequest(context.TODO(), req)
	if err != nil {
		log.Panic(err)
	}
	for resp := range ch {
		if resp.Err != nil {
			log.Panic("recv fail.", resp.Err)
		}
		fmt.Print(resp.Choices[0].Delta.Content)
	}

	// non-stream request
	fmt.Println("\n>>> test for non-stream request")
	req.Stream = false

	resp, err := engine.ChatRequest(context.TODO(), req)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(resp.Choices[0].Message.Content)
}
