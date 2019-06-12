package angmar

import "strings"

// AngmarMessage is a struct that encapsulates the message that Angmar
// listens to on a queue for.
type AngmarMessage struct {
	Url    string
	SHA    string
	Pusher string
	Tasks  []string
}

// String returns a stringified version of AngmarMessage, but doesn't
// stringify the list of Tasks.
func (m AngmarMessage) String() string {
	var builder strings.Builder
	builder.WriteString("URL: " + m.Url + "\n")
	builder.WriteString("SHA: " + m.SHA + "\n")
	builder.WriteString("Pusher: " + m.Pusher + "\n")
	return builder.String()
}
