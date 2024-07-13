package playbook

import (
	"aidezone/auto-gen-golang-agent/generator"
)

type Playbook struct {

	Robots []struct{
		Actor string
		Number int
	}

	Steps []map[string]*generator.Robot

}