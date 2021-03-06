package dolphin

// MessageSource is the source of a Minecraft message,
// either a player or the server.
type MessageSource string

const (
	// PlayerSource indicates that a message came from a player, e.g.
	// a chat message.
	PlayerSource MessageSource = "Player"
	// ServerSource indicates that a message came from the server,
	// such as when a player joins.
	ServerSource MessageSource = "Server"
)

// MinecraftMessage represents a message from Minecraft to be sent to Discord.
type MinecraftMessage struct {
	Username string        `json:"name"`
	Content  string        `json:"content"`
	Source   MessageSource `json:"source"`
	UUID     string        `json:"uuid"`
}
