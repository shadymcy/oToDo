package entity

type FileAccessType string

const (
	FileTypePublic FileAccessType = "public" // set RelatedID to empty
	FileTypeTodo   FileAccessType = "todo"   // set RelatedID to TodoID
)

type File struct {
	Entity

	FileName     string `json:"fileName"`
	FileServerID string `json:"-" gorm:"size:15"`
	FilePath     string `json:"-" gorm:"size:128"`
	AccessType   string `json:"accessType" gorm:"size:8"` // FileAccessType
	RelatedID    int64  `json:"relatedID"`                // Depend on access type
}
