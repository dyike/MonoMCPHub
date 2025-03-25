package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ark"
	eino_mcp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"github.com/dyike/MonoMCPHub/client/android_mcp_client"
)

func main() {
	ctx := context.Background()
	cmd := "/Users/bytedance/Code/go/bin/android_mcp_server"
	env := []string{
		"ANDROID_DEVICE_ID=100.80.162.9:5555",
		"ANDROID_WORK_DIR=/Users/bytedance/.mcp_tmp",
	}
	args := []string{}
	cli, err := android_mcp_client.NewAndroidMcpClient(ctx, cmd, env, args...)
	if err != nil {
		log.Fatal(err)
	}

	tools, _ := eino_mcp.GetTools(ctx, &eino_mcp.Config{Cli: cli})
	for _, tool := range tools {
		info, err := tool.Info(ctx)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(info)
	}
	defer cli.Close()

	llm, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		BaseURL: "https://ark.cn-beijing.volces.com/api/v3",
		Model:   os.Getenv("ARK_MODEL"),
		APIKey:  os.Getenv("ARK_API_KEY"),
	})

	ragent, err := react.NewAgent(ctx, &react.AgentConfig{
		Model: llm,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: tools,
		},
	})

	if err != nil {
		fmt.Println("failed to create agent:", err)
		return
	}

	sr, err := ragent.Stream(ctx, []*schema.Message{
		{
			Role: schema.System,
			Content: `# Character:
			你是一个资深的Android手机玩家,根据用户的指令,执行Android手机相关的操作`,
		},
		{
			Role: schema.User,
			Content: `# Character:
			用Android MCP Server，帮我设置一下声音闹钟音量为最低。每次操作前，都先截图，然后获取UI布局，最后再操作指令。`,
		},
	}, agent.WithComposeOptions(compose.WithCallbacks(&LoggerCallback{})))

	if err != nil {
		fmt.Println("failed to stream:", err)
		return
	}

	for {
		msg, err := sr.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				// finish
				break
			}
			// error
			fmt.Println("failed to recv:", err)
			return
		}

		// 打字机打印
		for _, char := range msg.Content {
			fmt.Print(string(char))
			time.Sleep(10 * time.Millisecond) // 调整延迟时间以获得合适的打字速度
		} // 打印完成后换行
	}

}

type LoggerCallback struct {
	callbacks.HandlerBuilder // 可以用 callbacks.HandlerBuilder 来辅助实现 callback
}

func (cb *LoggerCallback) OnStart(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
	fmt.Println("==================")
	inputStr, _ := json.MarshalIndent(input, "", "  ") // nolint: byted_s_returned_err_check
	fmt.Printf("[OnStart] %s\n", string(inputStr))
	return ctx
}

func (cb *LoggerCallback) OnEnd(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
	fmt.Println("=========[OnEnd]=========")
	outputStr, _ := json.MarshalIndent(output, "", "  ") // nolint: byted_s_returned_err_check
	fmt.Println(string(outputStr))
	return ctx
}

func (cb *LoggerCallback) OnError(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
	fmt.Println("=========[OnError]=========")
	fmt.Println(err)
	return ctx
}

func (cb *LoggerCallback) OnEndWithStreamOutput(ctx context.Context, info *callbacks.RunInfo,
	output *schema.StreamReader[callbacks.CallbackOutput]) context.Context {

	// go func() {
	// 	defer func() {
	// 		if err := recover(); err != nil {
	// 			fmt.Println("[OnEndStream] panic err:", err)
	// 		}
	// 	}()

	// 	defer output.Close() // remember to close the stream in defer

	// 	fmt.Println("=========[OnEndStream]=========")
	// 	for {
	// 		frame, err := output.Recv()
	// 		if errors.Is(err, io.EOF) {
	// 			// finish
	// 			break
	// 		}
	// 		if err != nil {
	// 			fmt.Printf("internal error: %s\n", err)
	// 			return
	// 		}

	// 		s, err := json.Marshal(frame)
	// 		if err != nil {
	// 			fmt.Printf("internal error: %s\n", err)
	// 			return
	// 		}

	// 		fmt.Printf("%s: %s\n", info.Name, string(s))
	// 	}
	// }()
	return ctx
}

func (cb *LoggerCallback) OnStartWithStreamInput(ctx context.Context, info *callbacks.RunInfo,
	input *schema.StreamReader[callbacks.CallbackInput]) context.Context {
	defer input.Close()
	return ctx
}
