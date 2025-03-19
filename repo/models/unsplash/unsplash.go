package models

// Photo represents an Unsplash photo
type Photo struct {
	ID          string   `json:"id"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	Width       int      `json:"width"`
	Height      int      `json:"height"`
	Color       string   `json:"color"`
	Description string   `json:"description"`
	AltDesc     string   `json:"alt_description"`
	URLs        PhotoURL `json:"urls"`
	Links       Links    `json:"links"`
	User        User     `json:"user"`
}

// PhotoURL contains various sizes of a photo
type PhotoURL struct {
	Raw     string `json:"raw"`
	Full    string `json:"full"`
	Regular string `json:"regular"`
	Small   string `json:"small"`
	Thumb   string `json:"thumb"`
}

// Links contains related links
type Links struct {
	Self     string `json:"self"`
	HTML     string `json:"html"`
	Download string `json:"download"`
}

// User represents an Unsplash user
type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Links     Links  `json:"links"`
}

// SearchResult represents the response from a search API request
type SearchResult struct {
	Total      int     `json:"total"`
	TotalPages int     `json:"total_pages"`
	Results    []Photo `json:"results"`
}
