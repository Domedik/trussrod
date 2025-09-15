package events

type Topic string

const (
	InterrogationUpdate Topic = "interrogation.update"
	NoteCreation        Topic = "note.creation"
	ExplorationCreation Topic = "interrogation.creation"
	ProfileUpdate       Topic = "profile.update"
)
