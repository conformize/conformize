// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package blueprint

import (
	"fmt"

	"github.com/conformize/conformize/internal/blueprint/elements"
	"gopkg.in/yaml.v3"
)

type Blueprint struct {
	Version    float64                                 `json:"version" yaml:"version"`
	Sources    map[string]elements.ConfigurationSource `json:"sources,omitempty" yaml:"sources"`
	References map[string]string                       `json:"$refs,omitempty" yaml:"$refs"`
	Ruleset    []elements.Rule                         `json:"ruleset,omitempty" yaml:"ruleset"`
}

func (b Blueprint) MarshalYAML() (any, error) {
	rootNode := yaml.Node{
		Kind: yaml.MappingNode,
	}

	rootNode.Content = append(rootNode.Content, &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: "version",
	})
	rootNode.Content = append(rootNode.Content, &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: fmt.Sprintf("%v", b.Version),
	})

	rootNode.Content = append(rootNode.Content, &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: "sources",
	})
	if len(b.Sources) > 0 {
		sourceNode := &yaml.Node{
			Kind: yaml.MappingNode,
		}
		for k, v := range b.Sources {
			sourceNode.Content = append(sourceNode.Content, &yaml.Node{
				Kind:  yaml.ScalarNode,
				Value: k,
			})
			sourceValueNode := &yaml.Node{}
			err := sourceValueNode.Encode(v)
			if err != nil {
				return nil, err
			}
			sourceNode.Content = append(sourceNode.Content, sourceValueNode)
		}
		rootNode.Content = append(rootNode.Content, sourceNode)
	} else {
		rootNode.Content = append(rootNode.Content, &yaml.Node{
			Kind: yaml.ScalarNode,
		})
	}

	rootNode.Content = append(rootNode.Content, &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: "$refs",
	})
	if len(b.References) > 0 {
		refNode := &yaml.Node{
			Kind: yaml.MappingNode,
		}
		for k, v := range b.References {
			refNode.Content = append(refNode.Content, &yaml.Node{
				Kind:  yaml.ScalarNode,
				Value: k,
			})
			refNode.Content = append(refNode.Content, &yaml.Node{
				Kind:  yaml.ScalarNode,
				Value: v,
			})
		}
		rootNode.Content = append(rootNode.Content, refNode)
	} else {
		rootNode.Content = append(rootNode.Content, &yaml.Node{
			Kind: yaml.ScalarNode,
		})
	}

	rootNode.Content = append(rootNode.Content, &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: "ruleset",
	})
	if len(b.Ruleset) > 0 {
		rulesetNode := &yaml.Node{
			Kind: yaml.SequenceNode,
		}
		for _, rule := range b.Ruleset {
			ruleNode := &yaml.Node{}
			err := ruleNode.Encode(rule)
			if err != nil {
				return nil, err
			}
			rulesetNode.Content = append(rulesetNode.Content, ruleNode)
		}
		rootNode.Content = append(rootNode.Content, rulesetNode)
	} else {
		rootNode.Content = append(rootNode.Content, &yaml.Node{
			Kind: yaml.ScalarNode,
		})
	}
	return &rootNode, nil
}
