package omniparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// 定义 API 响应的结构体
type APIResponse struct {
	Status           string                 `json:"status"`
	LabeledImage     string                 `json:"labeled_image"` // base64 encoded image
	ParsedContent    string                 `json:"parsed_content"`
	LabelCoordinates map[string]interface{} `json:"label_coordinates"`
	Message          string                 `json:"message"`
}

// 定义图标数据的结构体
type ScreenElement struct {
	ID            int       `json:"id"`
	Type          string    `json:"type"`
	Content       string    `json:"content"`
	Interactivity bool      `json:"interactivity"`
	Position      Position  `json:"position"`
	Bbox          []float64 `json:"bbox"`
}

// 定义位置的结构体
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// 处理图像
func ProcessImage(imagePath, apiURL string, boxThreshold, iouThreshold float64, usePaddleOCR bool, imgsz int) (APIResponse, error) {
	// 打开图片文件
	file, err := os.Open(imagePath)
	if err != nil {
		return APIResponse{}, fmt.Errorf("无法打开图片文件: %v", err)
	}
	defer file.Close()

	// 创建 multipart 表单
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "image.png")
	if err != nil {
		return APIResponse{}, fmt.Errorf("无法创建表单文件: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return APIResponse{}, fmt.Errorf("无法复制文件内容: %v", err)
	}

	// 添加查询参数
	_ = writer.WriteField("box_threshold", fmt.Sprintf("%f", boxThreshold))
	_ = writer.WriteField("iou_threshold", fmt.Sprintf("%f", iouThreshold))
	_ = writer.WriteField("use_paddleocr", fmt.Sprintf("%v", usePaddleOCR)) // ignore_security_alert
	_ = writer.WriteField("imgsz", fmt.Sprintf("%d", imgsz))

	// 关闭表单
	err = writer.Close()
	if err != nil {
		return APIResponse{}, fmt.Errorf("无法关闭表单: %v", err)
	}

	// 发送 POST 请求
	resp, err := http.Post(apiURL, writer.FormDataContentType(), body) // ignore_security_alert
	if err != nil {
		return APIResponse{}, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var result APIResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return APIResponse{}, fmt.Errorf("无法解析响应: %v", err)
	}
	return result, nil
}

// 解析图标数据
func ParseIconData(contentStr string) []ScreenElement {
	var elements []ScreenElement
	lines := strings.Split(strings.TrimSpace(contentStr), "\n")
	for _, line := range lines {
		if !strings.HasPrefix(line, "icon ") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		iconID := strings.TrimSpace(strings.TrimPrefix(parts[0], "icon "))
		contentText := strings.TrimSpace(parts[1])

		// Convert Python format to JSON format
		contentText = strings.ReplaceAll(contentText, "'", "\"")
		contentText = strings.ReplaceAll(contentText, "True", "true")
		contentText = strings.ReplaceAll(contentText, "False", "false")

		var content map[string]interface{}
		err := json.Unmarshal([]byte(contentText), &content)
		if err != nil {
			log.Printf("解析错误: %v", err)
			continue
		}
		// Calculate center point for interaction
		bbox := toFloat64Slice(content["bbox"])
		centerX := (bbox[0] + bbox[2]) / 2
		centerY := (bbox[1] + bbox[3]) / 2

		element := ScreenElement{
			ID:            atoi(iconID),
			Type:          toString(content["type"]),
			Content:       toString(content["content"]),
			Interactivity: toBool(content["interactivity"]),
			Position: Position{
				X: centerX,
				Y: centerY,
			},
			Bbox: toFloat64Slice(content["bbox"]),
		}

		elements = append(elements, element)
	}
	return elements
}

// 辅助函数：将字符串转换为整数
func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

// 辅助函数：将接口转换为字符串
func toString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// 辅助函数：将接口转换为布尔值
func toBool(v interface{}) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}

// 辅助函数：将接口转换为浮点数
func toFloat64(v interface{}) float64 {
	if f, ok := v.(float64); ok {
		return f
	}
	return 0.0
}

// 辅助函数：将接口转换为浮点数切片
func toFloat64Slice(v interface{}) []float64 {
	if slice, ok := v.([]interface{}); ok {
		result := make([]float64, len(slice))
		for i, val := range slice {
			result[i] = toFloat64(val)
		}
		return result
	}
	return nil
}
