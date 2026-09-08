// Copyright (C) 2026 XLR8discovery PBC
// See LICENSE for copying information.

package common

import (
	"strings"

	"xlr8d.io/bigbang-up/pkg/recipe"
)

// BuildDict stores the name of the container to build for Bigbang services.
var BuildDict = map[string]string{
	"authservice":          "app-edge",
	"gateway-mt":           "app-edge",
	"linksharing":          "app-edge",
	"satellite-admin":      "app-bigbang",
	"satellite-api":        "app-bigbang",
	"satellite-core":       "app-bigbang",
	"satellite-bf":         "app-bigbang",
	"satellite-gc":         "app-bigbang",
	"satellite-audit":      "app-bigbang",
	"satellite-rangedloop": "app-bigbang",
	"satellite-repair":     "app-bigbang",
	"storagenode":          "app-bigbang",
	"uplink":               "app-bigbang",
	"versioncontrol":       "app-bigbang",
}

// ResolveBuilds returns with the required docker images to build (as keys in the maps).
func ResolveBuilds(services []string) (map[string]string, error) {
	result := make(map[string]string)
	resolvedServices, err := ResolveServices(services)
	if err != nil {
		return result, err
	}
	for _, service := range resolvedServices {
		result[BuildDict[service]] = ""
	}
	return result, nil
}

// ResolveServices replaces group definition with exact services in the list.
func ResolveServices(selectors []string) ([]string, error) {
	var res []string
	st, err := recipe.GetStack()
	if err != nil {
		return res, err
	}
	for _, oneOrMoreSelector := range selectors {
		for selector := range strings.SplitSeq(oneOrMoreSelector, ",") {
			resolved := false
			for _, r := range st {
				if r.Name == selector {
					for _, rs := range r.Add {
						res = append(res, rs.Name)
						resolved = true
					}
				}
			}
			if !resolved {
				res = append(res, selector)
			}
		}
	}
	return res, nil
}

// ServiceMatches returns true if the service names match, regardless of the
// case or the numeric instance ID appended to the end of the service name.
func ServiceMatches(service1, service2 string) bool {
	// remove any numbers from the end of the strings
	service1 = strings.TrimRightFunc(service1, func(r rune) bool {
		return r >= '0' && r <= '9'
	})
	service2 = strings.TrimRightFunc(service2, func(r rune) bool {
		return r >= '0' && r <= '9'
	})
	return strings.EqualFold(service1, service2)
}
