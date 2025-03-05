package config

type ChatStringMap map[string]*Chat
type ChatIntMap map[int64]*Chat

type Config struct {
	Config      *Scheme
	Path        string
	ChatsPath   string
	PromptsPath string
}

type Scheme struct {
	EnableLogs     bool          `json:"enable_logs"`
	APIKey         string        `json:"api_key"`
	AvailableChats []string      `json:"available_chats"`
	Chats          ChatStringMap `json:"-"`
	IDChats        ChatIntMap    `json:"-"`
}

type Chat struct {
	ChatName          string              `json:"-"`
	ChatID            int64               `json:"chat_id"`
	ScheduledMessages []*ScheduledMessage `json:"scheduled_messages"`
	Replies           []*Reply            `json:"replies"`
	AiReplies         []*AiReply          `json:"ai_replies"`
}

type ScheduledMessage struct {
	When    []string `json:"when"`
	Days    []string `json:"days"`
	Message string   `json:"message"`
}

type Reply struct {
	If      string `json:"if"`
	Regex   string `json:"regex"`
	Message string `json:"message"`
}

type AiReply struct {
	If            string `json:"if"`
	Regex         string `json:"regex"`
	Prompt        string `json:"prompt"`
	ResponseRegex string `json:"responseRegex"`
	Model         string `json:"model"`
	PromptText    string `json:"-"`
}
