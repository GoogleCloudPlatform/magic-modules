package test

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// getSchemaDefault retrieves the default value in schema for the given key
func getSchemaDefault(res *schema.Resource, key string) interface{} {
	parts := strings.Split(key, ".")
	currentMap := res.Schema

	for i, part := range parts {
		// If part is number, skip lookup, just ensure we are potentially inside a list/set
		if _, err := strconv.Atoi(part); err == nil {
			continue
		}

		if currentMap == nil {
			return nil
		}

		s, ok := currentMap[part]
		if !ok {
			return nil
		}

		if i == len(parts)-1 {
			return s.Default
		}

		// Prepare for next
		if s.Elem != nil {
			if subRes, ok := s.Elem.(*schema.Resource); ok {
				currentMap = subRes.Schema
			} else {
				// Elem might be *schema.Schema (primitive list), no sub-map
				currentMap = nil
			}
		} else {
			currentMap = nil
		}
	}
	return nil
}

// getSchemaRequired checks if the given key is required in the schema
func getSchemaRequired(res *schema.Resource, key string) bool {
	parts := strings.Split(key, ".")
	currentMap := res.Schema

	for i, part := range parts {
		// If part is number, skip lookup, just ensure we are potentially inside a list/set
		if _, err := strconv.Atoi(part); err == nil {
			continue
		}

		if currentMap == nil {
			return false
		}

		s, ok := currentMap[part]
		if !ok {
			return false
		}

		if i == len(parts)-1 {
			return s.Required
		}

		// Prepare for next
		if s.Elem != nil {
			if subRes, ok := s.Elem.(*schema.Resource); ok {
				currentMap = subRes.Schema
			} else {
				// Elem might be *schema.Schema (primitive list), no sub-map
				currentMap = nil
			}
		} else {
			currentMap = nil
		}
	}
	return false
}
