package serialization

type PromptResult struct {
	ID             string `json:"id"`
	Group          string `json:"group"`
	Response       string `json:"response"`
	ResponseNumber int
}

type FormResult struct {
	PromptGroups map[string][]*PromptResult `json:"response_groups"`
}

func NewFormResult() *FormResult {
	return &FormResult{
		PromptGroups: make(map[string][]*PromptResult, 0),
	}
}

func (f *FormResult) AddPromptResult(promptResult *PromptResult) {
	if _, exists := f.PromptGroups[promptResult.Group]; !exists {
		f.PromptGroups[promptResult.Group] = make([]*PromptResult, 0)
	}

	iterCount := 1
	for _, p := range f.PromptGroups[promptResult.ID] {
		if p.ID == promptResult.ID {
			iterCount++
		}
	}

	promptResult.ResponseNumber = iterCount

	f.PromptGroups[promptResult.Group] = append(f.PromptGroups[promptResult.Group], promptResult)
}
