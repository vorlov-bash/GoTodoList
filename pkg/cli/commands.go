package cli

type Command string

const (
	CommandAdd             Command = "add"
	CommandRemove          Command = "remove"
	CommandMarkAsCompleted Command = "complete"
	CommandShow            Command = "show"
	CommandHelp            Command = "help"
)
