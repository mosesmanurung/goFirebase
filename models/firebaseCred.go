package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type FirebaseConfig struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain"`
}

func FirebaseCred() ([]byte, error) {
	filePath := "./fir-file-6a929-firebase-adminsdk-qnpgx-54c1e392f8.txt"

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("Error open file: %v", err)
	}
	defer file.Close()

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Error reading file: %v", err)
	}

	var config FirebaseConfig
	err = json.Unmarshal(fileContent, &config)
	if err != nil {
		return nil, fmt.Errorf("Error decoding JSON: %v", err)
	}

	jsonData, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("Error encoding JSON: %v", err)
	}

	return jsonData, nil
}
