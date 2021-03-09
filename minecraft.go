package dolphin

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/DataDrake/waterlog"
	"github.com/nxadm/tail"
)

// MinecraftWatcher watches for log lines from a Minecraft server.
type MinecraftWatcher struct {
	path          string
	deathKeywords []string
	uuidCache     map[string]string
	log           *waterlog.WaterLog
	tail          *tail.Tail
}

// NewWatcher creates a new watcher with all of the Minecraft death message keywords.
func NewWatcher(path string, logger *waterlog.WaterLog, customDeathKeywords []string) *MinecraftWatcher {
	var deathKeywords = []string{" shot", " pricked", " walked into a cactus", " roasted", " drowned", " kinetic", " blew up", " blown up", " killed", " hit the ground", " fell", " doomed", " squashed", " magic", " flames", " burned", " walked into fire", " burnt", " bang", " tried to swim in lava", " lightning", "floor was lava", "danger zone", " slain", " fireballed", " stung", " starved", " suffocated", " squished", " poked", " imapled", "didn't want to live", " withered", " pummeled", " died", " slain"}

	// Append any custom death keywords
	if len(customDeathKeywords) > 0 {
		deathKeywords = append(deathKeywords, customDeathKeywords...)
	}

	return &MinecraftWatcher{
		path:          path,
		log:           logger,
		deathKeywords: deathKeywords,
		uuidCache:     make(map[string]string),
	}
}

// Close stops the tail process and cleans up inotify file watches.
func (w *MinecraftWatcher) Close() error {
	err := w.tail.Stop()
	w.tail.Cleanup()
	return err
}

// GetUUID returns a UUID if present in the cache. See docs for Go maps.
// This is only meant to be used as a testing helper function.
func (w MinecraftWatcher) GetUUID(name string) (uuid string, ok bool) {
	uuid, ok = w.uuidCache[name]
	return
}

// Watch watches a log file for changes and sends Minecraft messages
// to the given channel.
func (w *MinecraftWatcher) Watch(c chan<- *MinecraftMessage) {
	// Check that the log file exists
	if _, err := os.Stat(w.path); err == nil {
		w.log.Infof("Using Minecraft log file at '%s'\n", w.path)

		// Start tailing the file
		var tailErr error
		w.tail, tailErr = tail.TailFile(w.path, tail.Config{
			Location: &tail.SeekInfo{
				Whence: io.SeekEnd,
			},
			ReOpen: true,
			Follow: true,
		})

		if tailErr != nil {
			w.log.Fatalf("Error trying to tail log file: %s\n", tailErr.Error())
		}

		w.log.Infoln("Log watcher started and waiting for lines")

		for {
			if line := <-w.tail.Lines; line != nil {
				if msg := w.ParseLine(line.Text); msg != nil {
					c <- msg
				}
			}
		}
	} else {
		w.log.Fatalf("Error opening log file: %s\n", err.Error())
	}
}

// ParseLine parses a log line for various types of messages and
// returns a MinecraftMessage struct if it is a message we care about.
func (w *MinecraftWatcher) ParseLine(line string) *MinecraftMessage {
	// Trim any line prefixes
	line = trimPrefix(line)
	if line == "" {
		return nil
	}

	// Trim trailing whitespace
	line = strings.TrimSpace(line)

	// Ignore villager death messages
	if strings.HasPrefix(line, "Villager") && strings.Contains(line, "died, message:") {
		return nil
	}

	// Check if the message is an auth message
	if strings.HasPrefix(line, "UUID of player") {
		parts := strings.Split(line, " ")
		name := parts[3]
		uuid := parts[5]
		w.uuidCache[name] = uuid
		return nil
	}

	// Check if the line is a chat message
	if strings.HasPrefix(line, "<") {
		// Split the message into parts
		parts := strings.SplitN(line, " ", 2)
		username := strings.TrimPrefix(parts[0], "<")
		username = strings.TrimSuffix(username, ">")
		message := parts[1]
		return &MinecraftMessage{
			Username: username,
			Content:  message,
			Source:   PlayerSource,
			UUID:     w.uuidCache[username],
		}
	} else if strings.Contains(line, "joined the game") || strings.Contains(line, "left the game") {
		// Remove from UUID cache when a player leaves
		if strings.Contains(line, "left the game") {
			name := strings.Fields(line)[0]
			delete(w.uuidCache, name)
		}
		return &MinecraftMessage{
			Username: "",
			Content:  line,
			Source:   ServerSource,
		}
	} else if isAdvancement(line) {
		return &MinecraftMessage{
			Username: "",
			Content:  fmt.Sprintf(":partying_face: %s", line),
			Source:   ServerSource,
		}
	} else if strings.HasPrefix(line, "Done (") {
		return &MinecraftMessage{
			Username: "",
			Content:  ":white_check_mark: Server has started",
			Source:   ServerSource,
		}
	} else if strings.HasPrefix(line, "Stopping the server") {
		return &MinecraftMessage{
			Username: "",
			Content:  ":x: Server is shutting down",
			Source:   ServerSource,
		}
	} else {
		// Check if the line is a death message
		for _, word := range w.deathKeywords {
			if strings.Contains(line, word) && line != "Found that the dragon has been killed in this world already." {
				return &MinecraftMessage{
					Username: "",
					Content:  fmt.Sprintf(":skull: %s", line),
					Source:   ServerSource,
				}
			}
		}
	}
	// Doesn't match anything we care about
	return nil
}

func isAdvancement(line string) bool {
	return strings.Contains(line, "has made the advancement") ||
		strings.Contains(line, "has completed the challenge") ||
		strings.Contains(line, "has reached the goal")
}

// trimPrefix trims the timestamp and thread prefix from incoming messages
// from the Minecraft server.
func trimPrefix(line string) string {
	// Some server plugins may log abnormal lines
	if !strings.HasPrefix(line, "[") || len(line) < 11 {
		return ""
	}

	start := strings.Index(line, "]: ") + 3

	// Trim the time prefix
	return line[start:]
}
