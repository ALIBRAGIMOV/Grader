package grader

type Config struct {
	ExternalGrader string       `json:"external_grader"`
	Files          []FileConfig `json:"files"`
	GraderPayload  Payload      `json:"grader_payload"`
}

type FileConfig struct {
	Label    string `json:"label"`
	FileName string `json:"filename"`
}

type Payload struct {
	Container string `json:"container"`
	PartID    string `json:"partId"`
}
