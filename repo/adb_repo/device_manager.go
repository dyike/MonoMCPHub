package adb_repo

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/disintegration/imaging"
)

type AdbRepo interface {
	GetPackages() (string, error)
	GetPackageActionIndents(packageName string) ([]string, error)
	ExecuteAdbCommand(args ...string) (string, error)
	TakeScreenshot() error
	GetUILayout() (string, error)
}

type adbRepoImpl struct {
	DeviceName string
}

func NewAdbRepo(deviceName string) AdbRepo {
	return &adbRepoImpl{
		DeviceName: deviceName,
	}
}

func (r *adbRepoImpl) GetPackages() (string, error) {
	args := []string{"pm", "list", "packages"}
	output, err := r.runAdbCommand(args...)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		if strings.HasPrefix(line, "package:") {
			result.WriteString(line[8:])
			result.WriteString("\n")
		}
	}
	return result.String(), nil
}

func (r *adbRepoImpl) GetPackageActionIndents(packageName string) ([]string, error) {
	cmdArgs := []string{"dumpsys", "package", packageName}
	output, err := r.runAdbCommand(cmdArgs...)
	if err != nil {
		return nil, err
	}
	resolverTableStart := strings.Index(output, "Resolver Table:")
	if resolverTableStart == -1 {
		return []string{}, nil
	}

	resolverSection := output[resolverTableStart:]

	nonDataStart := strings.Index(resolverSection, "\n  Non-Data Actions:")
	if nonDataStart == -1 {
		return []string{}, nil
	}

	nonDataSection := resolverSection[nonDataStart:]
	sectionEnd := strings.Index(nonDataSection, "\n\n")
	if sectionEnd != -1 {
		nonDataSection = nonDataSection[:sectionEnd]
	}

	var actions []string
	for _, line := range strings.Split(nonDataSection, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "android.") || strings.HasPrefix(line, "com.") {
			actions = append(actions, line)
		}
	}

	return actions, nil
}

func (r *adbRepoImpl) ExecuteAdbCommand(args ...string) (string, error) {
	return r.runAdbCommand(args...)
}

func (r *adbRepoImpl) TakeScreenshot() error {
	// 截屏
	_, err := r.runAdbCommand("screencap", "-p", "/sdcard/screenshot.png")
	if err != nil {
		return err
	}

	// 拉取截图
	pullCmd := exec.Command("adb", "-s", r.DeviceName, "pull", "/sdcard/screenshot.png", "screenshot.png")
	err = pullCmd.Run()
	if err != nil {
		return err
	}

	// 删除截图
	_, err = r.runAdbCommand("shell", "rm", "/sdcard/screenshot.png")
	if err != nil {
		return err
	}

	// 压缩截图
	file, err := os.Open("screenshot.png")
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// 调整为原始尺寸的30%
	bounds := img.Bounds()
	newWidth := int(float64(bounds.Dx()) * 0.3)
	newHeight := int(float64(bounds.Dy()) * 0.3)

	resized := imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)

	outFile, err := os.Create("compressed_screenshot.png")
	if err != nil {
		return err
	}
	defer outFile.Close()

	err = png.Encode(outFile, resized)
	if err != nil {
		return err
	}

	return nil
}

// GetUILayout 获取并分析UI布局
func (r *adbRepoImpl) GetUILayout() (string, error) {
	// 使用uiautomator dump UI
	_, err := r.runAdbCommand("uiautomator dump")
	if err != nil {
		return "", err
	}

	// 拉取XML到本地
	pullCmd := exec.Command("adb", "-s", r.DeviceName, "pull", "/sdcard/window_dump.xml", "window_dump.xml")
	err = pullCmd.Run()
	if err != nil {
		return "", err
	}

	// 删除设备上的文件
	_, err = r.runAdbCommand("rm /sdcard/window_dump.xml")
	if err != nil {
		return "", err
	}

	// 解析XML
	file, err := os.Open("window_dump.xml")
	if err != nil {
		return "", err
	}
	defer file.Close()

	doc, err := xmlquery.Parse(file)
	if err != nil {
		return "", err
	}

	clickableElements := []string{}
	nodes := xmlquery.Find(doc, "//node[@clickable='true']")

	boundsRegex := regexp.MustCompile(`\[(\d+),(\d+)\]`)

	for _, node := range nodes {
		text := node.SelectAttr("text")
		contentDesc := node.SelectAttr("content-desc")
		bounds := node.SelectAttr("bounds")

		// 只包含有文本或描述的元素
		if text != "" || contentDesc != "" {
			elementInfo := "Clickable element:"

			if text != "" {
				elementInfo += fmt.Sprintf("\n  Text: %s", text)
			}

			if contentDesc != "" {
				elementInfo += fmt.Sprintf("\n  Description: %s", contentDesc)
			}

			elementInfo += fmt.Sprintf("\n  Bounds: %s", bounds)

			// 计算中心点
			matches := boundsRegex.FindAllStringSubmatch(bounds, -1)
			if len(matches) == 2 {
				x1, y1 := parseInt(matches[0][1]), parseInt(matches[0][2])
				x2, y2 := parseInt(matches[1][1]), parseInt(matches[1][2])
				centerX := (x1 + x2) / 2
				centerY := (y1 + y2) / 2
				elementInfo += fmt.Sprintf("\n  Center: (%d, %d)", centerX, centerY)
			}

			clickableElements = append(clickableElements, elementInfo)
		}
	}

	if len(clickableElements) == 0 {
		return "No clickable elements found with text or description", nil
	}

	return strings.Join(clickableElements, "\n\n"), nil
}

func (r *adbRepoImpl) runAdbCommand(args ...string) (string, error) {
	command := strings.Join(args, " ")
	var cmd *exec.Cmd

	if strings.HasPrefix(command, "adb shell") {
		command = strings.TrimPrefix(command, "adb shell")
		cmd = exec.Command("adb", "-s", r.DeviceName, "shell", command)
	} else if strings.HasPrefix(command, "adb ") {
		command = strings.TrimPrefix(command, "adb ")
		args = strings.Split(command, " ")
		cmdArgs := []string{"-s", r.DeviceName}
		cmdArgs = append(cmdArgs, args...)
		cmd = exec.Command("adb", cmdArgs...)
	} else {
		cmd = exec.Command("adb", "-s", r.DeviceName, "shell", command)
	}

	var output bytes.Buffer
	cmd.Stdout = &output
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run adb command(%s): %w", cmd.String(), err)
	}
	return output.String(), nil
}

// parseInt 将字符串转换为整数
func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}
