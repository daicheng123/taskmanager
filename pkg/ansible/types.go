package ansible

import (
	"encoding/json"
	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
	errors "github.com/apenella/go-common-utils/error"
)

type AnsibleModule string

const (
	CommandModule         = "command"
	StrictHostKeyChecking = "-o StrictHostKeyChecking=no"
	SmartConnection       = "smart"
	ForkNumber            = "10"
	BecomeUser            = "root"

	AnsibleRetryFilesEnabled = "ANSIBLE_RETRY_FILES_ENABLED"
)

type AnsibleJSONResults struct {
	Playbook          string                                              `json:"-"`
	CustomStats       interface{}                                         `json:"custom_stats"`
	GlobalCustomStats interface{}                                         `json:"global_custom_stats"`
	Plays             []AnsiblePlaybookJSONResultsPlay                    `json:"plays"`
	Stats             map[string]*results.AnsiblePlaybookJSONResultsStats `json:"stats"`
}

// AnsiblePlaybookJSONResultsPlay
type AnsiblePlaybookJSONResultsPlay struct {
	Play  *results.AnsiblePlaybookJSONResultsPlaysPlay `json:"play"`
	Tasks []*AnsiblePlaybookJSONResultsPlayTask        `json:"tasks"`
}

type AnsiblePlaybookJSONResultsPlayTask struct {
	Task  *results.AnsiblePlaybookJSONResultsPlayTaskItem         `json:"task"`
	Hosts map[string]*AnsiblePlaybookJSONResultsPlayTaskHostsItem `json:"hosts"`
}

type AnsiblePlaybookJSONResultsPlayTaskHostsItem struct {
	*results.AnsiblePlaybookJSONResultsPlayTaskHostsItem `json:",inline"`
	*AnsiblePlaybookJSONResultsPlayTaskHostsItemResults  `json:",inline"`
}

type AnsiblePlaybookJSONResultsPlayTaskHostsItemResults struct {
	Results []*AnsiblePlaybookJSONResultsPlayTaskHostsItemResult `json:"results"`
}

type AnsiblePlaybookJSONResultsPlayTaskHostsItemResult struct {
	Message     interface{} `json:"msg"`
	Unreachable bool        `json:"unreachable"`
	Failed      bool        `json:"failed"`
}

// JSONParse return an AnsiblePlaybookJSONResults from
func JSONParse(data []byte) (*AnsibleJSONResults, error) {

	result := &AnsibleJSONResults{}

	err := json.Unmarshal(data, result)
	if err != nil {
		return nil, errors.New("(results::JSONParser)", "Unmarshall error", err)
	}

	return result, nil
}
