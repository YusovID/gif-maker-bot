package domain

type FileType string

const (
	Video FileType = "video"
	GIF   FileType = "gif"
	Photo FileType = "photo"
)

type Message struct {
	From    User
	Chat    Chat
	File    File
	ReplyTo int
}

type User struct {
	ID       int64
	Username string
}

type Chat struct {
	ChatID int64
}

type File struct {
	ID       string
	Type     FileType
	FileName string
	Height   int
	Width    int
	MymeType string
}
