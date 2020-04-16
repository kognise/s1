package main

import (
	"encoding/gob"
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"os"
	"strings"
)

var tokens = make(map[string]struct{})

func storeTokens() {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(tokens)
	if err != nil {
		panic("Error encoding glob!")
	}

	ioutil.WriteFile(tokensFile, buffer.Bytes(), 0644)
}

func loadTokens() {
	fileBytes, err := ioutil.ReadFile(tokensFile)
	if err != nil {
		panic("Error reading tokens file!")
	}

	buffer := bytes.Buffer{}
	buffer.Write(fileBytes)

	decoder := gob.NewDecoder(&buffer)
	err = decoder.Decode(&tokens)
	if err != nil {
		panic("Error decoding tokens file!")
	}
}

func generateToken() string {
	token := generateRandomString(64)
	tokens[token] = struct{}{}
	storeTokens()
	return token
}

func getDataLoc(token string, key string) (string, string) {
	dataFileDir := dataDir + token + "/"
	dataFile := dataFileDir + base64.URLEncoding.EncodeToString([]byte(key)) + suffix

	return dataFileDir, dataFile
}

func storeData(token string, key string, data []byte) {
	dataFileDir, dataFile := getDataLoc(token, key)
	os.MkdirAll(dataFileDir, os.ModePerm)
	ioutil.WriteFile(dataFile, data, 0644)
}

func getData(token string, key string) []byte {
	_, dataFile := getDataLoc(token, key)
	
	fileBytes, err := ioutil.ReadFile(dataFile)
	if err != nil {
		return nil
	}
	return fileBytes
}

func deleteData(token string, key string) {
	_, dataFile := getDataLoc(token, key)
	os.Remove(dataFile)
}

func getKeys(token string) []string {
	dataFileDir, _ := getDataLoc(token, "")

	f, err := os.Open(dataFileDir)
	if err != nil {
		return nil
	}

	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil
	}

	keys := make([]string, len(files))
	for i, file := range files {
		name := file.Name()
		if !strings.HasSuffix(name, suffix) {
			continue
		}
		name = strings.TrimSuffix(name, suffix)
		decoded, err := base64.URLEncoding.DecodeString(name)
		if err != nil {
			continue
		}
		keys[i] = string(decoded)
	}
	return keys
}