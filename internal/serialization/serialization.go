package serialization

import (
	"encoding/json"
	"os"
)

type PromptJSON struct {
	Id          string            `json:"id"`
	Type        string            `json:"type"`
	AnswerType  string            `json:"answer_type"`
	Title       string            `json:"title"`
	Desc        *string           `json:"description"`
	Placeholder *string           `json:"placeholder"`
	Choices     []string          `json:"choices,omitempty"`
	Next        map[string]string `json:"next"`
	Group       string            `json:"group"`
	BindSubmit  *string           `json:"bind_submit"`
}

type FormJSON struct {
	First   string                 `json:"first_prompt"`
	Prompts map[string]*PromptJSON `json:"prompts"`
	Vars    map[string]string      `json:"vars"`
}

func NewFormJSON(path string) (*FormJSON, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var formJSON FormJSON
	err = json.Unmarshal(contents, &formJSON)
	return &formJSON, err
}
