package prompttree

import (
	"goformit/internal/logging"
	"sort"
	"strings"
)

// used to check if a response *qualifies* for a new model
type Qualifier interface {
	Qualifies(input []string) bool
	ModelID() string
}

// base impl qualifier
type BasicQualifier struct {
	Options     [][]string
	nextModelId string
}

func createQualifierOption(str string) []string {
	items := strings.Split(str, "+")
	sort.Strings(items)
	logging.AppLogger.Log("qualifier_init", "%v", items)
	return items
}

func NewBasicQualifier(rawStr string, nextModelId string) *BasicQualifier {
	optsRaw := strings.Split(rawStr, ",")
	opts := make([][]string, 0)
	for _, opt := range optsRaw {
		opts = append(opts, createQualifierOption(opt))
	}

	return &BasicQualifier{
		Options:     opts,
		nextModelId: nextModelId,
	}
}

func (b *BasicQualifier) Qualifies(input []string) bool {
	sort.Strings(input)

	qStatus := false
	for _, opt := range b.Options {
		logging.AppLogger.Log("qualifier", "checking %v against %v", input, opt)

		if opt[0] == "_" {
			return true
		}

		if len(input) != len(opt) {
			continue
		}

		passes := true
		for i, item := range input {
			if item != opt[i] {
				passes = false
			}
		}

		if passes {
			return true
		}
	}
	return qStatus
}

func (b *BasicQualifier) ModelID() string {
	return b.nextModelId
}
