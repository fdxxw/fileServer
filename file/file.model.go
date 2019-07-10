package file

// FileInfo 文件基本信息
type FileInfo struct {
	Filename string
	Metadata struct {
		UserID   string `json:"userId" bson:"userId"`
		Username string
		Tags     []string
		Public   bool
	}
}
