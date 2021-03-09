package dolphin

import (
	"testing"
)

var watcher = MinecraftWatcher{
	deathKeywords: make([]string, 0),
	uuidCache: map[string]string{
		"TestUser": "7f7c909b-24f1-49a4-817f-baa4f4973980",
	},
}

func TestParseVanillaChatLine(t *testing.T) {
	// Given
	input := "[12:32:45] [Server thread/INFO]: <TestUser> Sending a chat message"
	expected := &MinecraftMessage{
		Username: "TestUser",
		Content:  "Sending a chat message",
		Source:   PlayerSource,
		UUID:     "7f7c909b-24f1-49a4-817f-baa4f4973980",
	}
	// When
	actual := watcher.ParseLine(input)
	// Then
	if actual.Username != expected.Username {
		t.Errorf("Parsing chat line got incorrect username, got: %s, expected: %s", actual.Username, expected.Username)
	}
	if actual.Content != expected.Content {
		t.Errorf("Parsing chat line got incorrect message, got: %s, expected: %s", actual.Content, expected.Content)
	}
	if actual.Source != expected.Source {
		t.Errorf("Parsing chat line got incorrect source, got: %s, expected: %s", actual.Source, expected.Source)
	}
	if actual.UUID != expected.UUID {
		t.Errorf("Parsing chat line got incorrect UUID, got: %s, expected %s", actual.UUID, expected.UUID)
	}
}

func TestParseNonVanillaChatLine(t *testing.T) {
	// Given
	input := "[12:32:45] [Async Chat Thread - #0/INFO]: <TestUser> Sending a chat message"
	expected := &MinecraftMessage{
		Username: "TestUser",
		Content:  "Sending a chat message",
		Source:   PlayerSource,
		UUID:     "7f7c909b-24f1-49a4-817f-baa4f4973980",
	}
	// When
	actual := watcher.ParseLine(input)
	// Then
	if actual.Username != expected.Username {
		t.Errorf("Parsing chat line got incorrect username, got: %s, expected: %s", actual.Username, expected.Username)
	}
	if actual.Content != expected.Content {
		t.Errorf("Parsing chat line got incorrect message, got: %s, expected: %s", actual.Content, expected.Content)
	}
	if actual.Source != expected.Source {
		t.Errorf("Parsing chat line got incorrect source, got: %s, expected: %s", actual.Source, expected.Source)
	}
	if actual.UUID != expected.UUID {
		t.Errorf("Parsing chat line got incorrect UUID, got: %s, expected %s", actual.UUID, expected.UUID)
	}
}

func TestParseVanillaJoinLine(t *testing.T) {
	// Given
	input := "[12:32:45] [Server thread/INFO]: TestUser joined the game"
	expected := &MinecraftMessage{
		Username: "",
		Content:  "TestUser joined the game",
		Source:   ServerSource,
	}
	// When
	actual := watcher.ParseLine(input)
	// Then
	if actual.Username != expected.Username {
		t.Errorf("Parsing chat line got incorrect username, got: %s, expected: %s", actual.Username, expected.Username)
	}
	if actual.Content != expected.Content {
		t.Errorf("Parsing chat line got incorrect message, got: %s, expected: %s", actual.Content, expected.Content)
	}
	if actual.Source != expected.Source {
		t.Errorf("Parsing chat line got incorrect source, got: %s, expected: %s", actual.Source, expected.Source)
	}
}

func TestParseLeaveLine(t *testing.T) {
	// Given
	input := "[12:32:45] [Server thread/INFO]: TestUser left the game"
	expected := &MinecraftMessage{
		Username: "",
		Content:  "TestUser left the game",
		Source:   ServerSource,
	}
	// When
	actual := watcher.ParseLine(input)
	// Then
	if actual.Username != expected.Username {
		t.Errorf("Parsing chat line got incorrect username, got: %s, expected: %s", actual.Username, expected.Username)
	}
	if actual.Content != expected.Content {
		t.Errorf("Parsing chat line got incorrect message, got: %s, expected: %s", actual.Content, expected.Content)
	}
	if actual.Source != expected.Source {
		t.Errorf("Parsing chat line got incorrect source, got: %s, expected: %s", actual.Source, expected.Source)
	}
	if uuid, ok := watcher.GetUUID("TestUser"); ok && uuid != "" {
		t.Error("Player UUID still present in cache after player left")
	}
}

func TestParseAdvancement1Line(t *testing.T) {
	// Given
	input := "[12:32:45] [Server thread/INFO]: TestUser has made the advancement [MonsterHunter]"
	expected := &MinecraftMessage{
		Username: "",
		Content:  ":partying_face: TestUser has made the advancement [MonsterHunter]",
		Source:   ServerSource,
	}
	// When
	actual := watcher.ParseLine(input)
	// Then
	if actual.Username != expected.Username {
		t.Errorf("Parsing chat line got incorrect username, got: %s, expected: %s", actual.Username, expected.Username)
	}
	if actual.Content != expected.Content {
		t.Errorf("Parsing chat line got incorrect message, got: %s, expected: %s", actual.Content, expected.Content)
	}
	if actual.Source != expected.Source {
		t.Errorf("Parsing chat line got incorrect source, got: %s, expected: %s", actual.Source, expected.Source)
	}
}

func TestParseAdvancement2Line(t *testing.T) {
	// Given
	input := "[12:32:45] [Server thread/INFO]: TestUser has completed the challenge [MonsterHunter]"
	expected := &MinecraftMessage{
		Username: "",
		Content:  ":partying_face: TestUser has completed the challenge [MonsterHunter]",
		Source:   ServerSource,
	}
	// When
	actual := watcher.ParseLine(input)
	// Then
	if actual.Username != expected.Username {
		t.Errorf("Parsing chat line got incorrect username, got: %s, expected: %s", actual.Username, expected.Username)
	}
	if actual.Content != expected.Content {
		t.Errorf("Parsing chat line got incorrect message, got: %s, expected: %s", actual.Content, expected.Content)
	}
	if actual.Source != expected.Source {
		t.Errorf("Parsing chat line got incorrect source, got: %s, expected: %s", actual.Source, expected.Source)
	}
}

func TestParseServerStartLine(t *testing.T) {
	// Given
	input := "[12:32:45] [Server thread/INFO]: Done (21.3242s)! For help, type \"help\""
	expected := &MinecraftMessage{
		Username: "",
		Content:  ":white_check_mark: Server has started",
		Source:   ServerSource,
	}
	// When
	actual := watcher.ParseLine(input)
	// Then
	if actual.Username != expected.Username {
		t.Errorf("Parsing chat line got incorrect username, got: %s, expected: %s", actual.Username, expected.Username)
	}
	if actual.Content != expected.Content {
		t.Errorf("Parsing chat line got incorrect message, got: %s, expected: %s", actual.Content, expected.Content)
	}
	if actual.Source != expected.Source {
		t.Errorf("Parsing chat line got incorrect source, got: %s, expected: %s", actual.Source, expected.Source)
	}
}

func TestParseServerStopLine(t *testing.T) {
	// Given
	input := "[12:32:45] [Server thread/INFO]: Stopping the server"
	expected := &MinecraftMessage{
		Username: "",
		Content:  ":x: Server is shutting down",
		Source:   ServerSource,
	}
	// When
	actual := watcher.ParseLine(input)
	// Then
	if actual.Username != expected.Username {
		t.Errorf("Parsing chat line got incorrect username, got: %s, expected: %s", actual.Username, expected.Username)
	}
	if actual.Content != expected.Content {
		t.Errorf("Parsing chat line got incorrect message, got: %s, expected: %s", actual.Content, expected.Content)
	}
	if actual.Source != expected.Source {
		t.Errorf("Parsing chat line got incorrect source, got: %s, expected: %s", actual.Source, expected.Source)
	}
}

func TestIgnoreVillagerDeath(t *testing.T) {
	// Given
	input := "[12:32:45] [Server thread/INFO]: Villager axw['Villager'/85, l='world', x=-147.30, y=57.00, z=-190.70] died, message: 'Villager was squished too much'"

	// When
	result := watcher.ParseLine(input)

	// Then
	if result != nil {
		t.Errorf("Parsing line failed to ignore villager death message, got: %s", result)
	}
}

func TestAddsUUID(t *testing.T) {
	// Given
	input := "[19:54:56] [User Authenticator #1/INFO]: UUID of player Bob is 7f7c909b-24f1-49a4-817f-baa4f4973980"

	// When
	result := watcher.ParseLine(input)

	// Then
	if result != nil {
		t.Errorf("Parsing line returned a message for an authentication line")
	}
	if uuid, ok := watcher.GetUUID("Bob"); !ok || uuid != "7f7c909b-24f1-49a4-817f-baa4f4973980" {
		t.Error("Player UUID not present in cache after auth message")
	}
}
