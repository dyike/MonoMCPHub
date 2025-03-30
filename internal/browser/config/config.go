package config

type BrowserConfig struct {
	Headless        bool
	Timeout         int
	Proxy           string
	UserAgent       string
	DefaultLanguage string
	URLTimeout      int
	CSSTimeout      int
	DataPath        string
}

func NewBrowserConfig() *BrowserConfig {
	return &BrowserConfig{
		Headless:        true,
		Timeout:         30,
		UserAgent:       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		DefaultLanguage: "zh-CN",
		URLTimeout:      30,
		CSSTimeout:      30,
		DataPath:        "/Users/ityike/.browser_data",
	}
}
