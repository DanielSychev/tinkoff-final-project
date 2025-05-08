package ads

type Ad struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Text        string `json:"text"`
	AuthorID    int64  `json:"author_id"`
	Published   bool   `json:"published"`
	Deleted     bool   `json:"deleted"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
}

type User struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Deleted bool   `json:"deleted"`
}

type AdFilter struct {
	Pub   bool
	Auth  int64
	Title string
}
