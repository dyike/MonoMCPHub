package tools

import (
	"context"
	"fmt"

	"github.com/dyike/MonoMCPHub/repo/adb_repo"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewGetPackagesTool() mcp.Tool {
	return mcp.NewTool("get_packages",
		mcp.WithDescription("Get all packages of your android device"),
		mcp.WithString("package_option",
			// -3: 只显示第三方应用包（用户安装的应用）
			// -s: 只显示系统应用包
			// -f: 显示应用包名及其关联的 APK 文件路径
			// -d: 只显示已禁用的应用包
			// -e: 只显示已启用的应用包
			// -u: 也包括已卸载但数据未清除的应用包
			mcp.Description(`The option to get packages ('', -3, -s, -f, -d, -e, -u)
				-3: 只显示第三方应用包（用户安装的应用）
				-s: 只显示系统应用包
				-f: 显示应用包名及其关联的 APK 文件路径
				-d: 只显示已禁用的应用包
				-e: 只显示已启用的应用包
				-u: 也包括已卸载但数据未清除的应用包
			`),
		),
	)
}

func HandleGetPackages(adbRepo adb_repo.AdbRepo) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		packageOption := req.Params.Arguments["package_option"].(string)
		if packageOption != "" &&
			packageOption != "-3" &&
			packageOption != "-s" &&
			packageOption != "-f" &&
			packageOption != "-d" &&
			packageOption != "-e" &&
			packageOption != "-u" {
			packageOption = ""
		}
		packages, err := adbRepo.GetPackages(packageOption)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to get packages: %v", err)
			return mcp.NewToolResultError(errMsg), nil
		}
		return mcp.NewToolResultText(packages), nil
	}
}
