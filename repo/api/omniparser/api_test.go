package omniparser

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/dyike/MonoMCPHub/repo/vlm"
)

func TestParseImage(t *testing.T) {
	imagePath := "/Users/bytedance/.mcp_tmp/screenshot.png"
	apiURL := "http://10.37.110.115:8000/process_image"

	resp, err := ProcessImage(imagePath, apiURL, 0.05, 0.1, true, 640)
	if err != nil {
		t.Fatalf("ProcessImage 失败: %v", err)
	}
	base64Image := resp.LabeledImage

	debugPath := "/Users/bytedance/.mcp_tmp/debug.png"
	debugImage, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		t.Fatalf("base64 解码失败: %v", err)
	}

	err = os.WriteFile(debugPath, debugImage, 0644)
	if err != nil {
		t.Fatalf("写入文件失败: %v", err)
	}

	parsedContent := resp.ParsedContent

	eles := ParseIconData(parsedContent)
	jsonData, _ := json.Marshal(eles)

	userRequest := "在应用宝中搜索得物APP并下载"
	resp1, err := vlm.DoubaoVision(context.Background(), base64Image, userRequest, string(jsonData))
	if err != nil {
		t.Fatalf("DoubaoVision 失败: %v", err)
	}
	fmt.Println(resp1)
}

func TestParseIconData(t *testing.T) {
	contentStr := `
	icon 0: {'type': 'text', 'bbox': [0.0729166641831398, 0.9444444179534912, 0.33125001192092896, 0.970370352268219], 'interactivity': False, 'content': '[2016M-CSG005905', 'source': 'box_ocr_content_ocr'}
icon 1: {'type': 'icon', 'bbox': [0.8482794165611267, 0.7801671028137207, 0.9184102416038513, 0.9347067475318909], 'interactivity': True, 'content': '12+ CADPA ', 'source': 'box_yolo_content_ocr'}
icon 2: {'type': 'icon', 'bbox': [0.9254415035247803, 0.39590007066726685, 0.9863064885139465, 0.49570611119270325], 'interactivity': True, 'content': 'Xiaohongshu (Little Red Book) app.', 'source': 'box_yolo_content_yolo'}
icon 3: {'type': 'icon', 'bbox': [0.018426287919282913, 0.033579021692276, 0.07202623039484024, 0.13014227151870728], 'interactivity': True, 'content': 'Settings or configuration options.', 'source': 'box_yolo_content_yolo'}
icon 4: {'type': 'icon', 'bbox': [0.9280473589897156, 0.15542253851890564, 0.9824953079223633, 0.2695449888706207], 'interactivity': True, 'content': 'QR code scanning', 'source': 'box_yolo_content_yolo'}
icon 5: {'type': 'icon', 'bbox': [0.2924197018146515, 0.7464330196380615, 0.49618199467658997, 0.8324514031410217], 'interactivity': True, 'content': 'A video streaming platform.', 'source': 'box_yolo_content_yolo'}
icon 6: {'type': 'icon', 'bbox': [0.9281579852104187, 0.033724021166563034, 0.9823734164237976, 0.1427483856678009], 'interactivity': True, 'content': 'Play button', 'source': 'box_yolo_content_yolo'}
icon 7: {'type': 'icon', 'bbox': [0.9270393252372742, 0.27518609166145325, 0.9819839596748352, 0.38946351408958435], 'interactivity': True, 'content': 'System', 'source': 'box_yolo_content_yolo'}
icon 8: {'type': 'icon', 'bbox': [0.5020170211791992, 0.7453713417053223, 0.712048351764679, 0.8331449627876282], 'interactivity': True, 'content': 'A video-related application or feature.', 'source': 'box_yolo_content_yolo'}
icon 9: {'type': 'icon', 'bbox': [0.1241692528128624, 0.9674462080001831, 0.1453464776277542, 0.9944751858711243], 'interactivity': True, 'content': 'Xray', 'source': 'box_yolo_content_yolo'}
icon 10: {'type': 'icon', 'bbox': [0.06745243072509766, 0.032150112092494965, 0.20324555039405823, 0.12497897446155548], 'interactivity': True, 'content': 'a user interface or account.', 'source': 'box_yolo_content_yolo'}
	`

	elements := ParseIconData(contentStr)

	for _, element := range elements {
		fmt.Println(element)
	}
}
