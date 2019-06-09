// +build integration

package main

import (
	"net/http"
	"testing"

	"gopkg.in/h2non/baloo.v3"
)

const (
	// TODO: from env?
	localURL = "http://localhost:8000"

	statsSchema = `{
		"title": "stats resource",
		"type":"object",
		"properties": {
			"hits": {
				"type":"integer"
			},
			"request": {
				"type":"object",
				"properties":{
					"int1":{
						"type":"integer"
					},
					"int2":{
						"type":"integer"
					},
					"limit":{
						"type":"integer"
					},
					"str1":{
						"type":"string"
					},
					"str2":{
						"type":"string"
					}
				},
				"required":["int1", "int2", "limit", "str1", "str2"]
			}
		},
		"required": ["hits", "request"]
	}`

	fizzbuzzSchema = `{
		"title": "fizzbuzz resource",
		"type":"object",
		"properties": {
			"result": {
				"type":"array"
			}
		},
		"required": ["result"]
	}`

	errorSchema = `{
		"title": "error resource",
		"type": "object",
		"properties": {
			"error": {
				"type": "object",
				"properties": {
					"message": {
						"type": "string"
					},
					"fields": {
						"type": "array"
					}
				},
				"required": ["message"]
			}
		}	
	}`
)

func Test_Integration_GetStats(t *testing.T) {
	client := baloo.New(localURL)
	tests := []struct {
		name               string
		expectedStatusCode int
		expectedSchema     string
	}{
		{
			name:               "OK",
			expectedStatusCode: http.StatusOK,
			expectedSchema:     statsSchema,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.Get("/v1/fizzbuzz/stats/").
				Expect(t).
				Status(tt.expectedStatusCode).
				JSONSchema(tt.expectedSchema).
				Done()
		})
	}
}

func Test_Integration_GetFizzbuzz(t *testing.T) {
	client := baloo.New(localURL)
	tests := []struct {
		name               string
		params             map[string]string
		expectedStatusCode int
		expectedSchema     string
	}{
		{
			name: "OK",
			params: map[string]string{
				"int1":  "3",
				"int2":  "5",
				"limit": "15",
				"str1":  "fizz",
				"str2":  "buzz",
			},
			expectedStatusCode: http.StatusOK,
			expectedSchema:     fizzbuzzSchema,
		},
		{
			name: "KO_InvalidRequest",
			params: map[string]string{
				"int1":  "not_a_int",
				"int2":  "5",
				"limit": "15",
				"str1":  "fizz",
				"str2":  "buzz",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedSchema:     errorSchema,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.Get("/v1/fizzbuzz/").
				SetQueryParams(tt.params).
				Expect(t).
				Status(tt.expectedStatusCode).
				JSONSchema(tt.expectedSchema).
				Done()
		})
	}
}
