package angmar

import "strings"

type AngmarMessage struct {
	Url    string
	SHA    string
	Pusher string
	Tasks  []string
}

func (m AngmarMessage) String() string {
	var builder strings.Builder
	builder.WriteString("URL: " + m.Url + "\n")
	builder.WriteString("SHA: " + m.SHA + "\n")
	builder.WriteString("Pusher: " + m.Pusher + "\n")
	return builder.String()
}
