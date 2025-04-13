package secret

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

func File(encodingkey []byte, file string) *Vault {
	return &Vault{
		encodingKey: encodingkey,
		filepath:    file,
		keyValues:   make(map[string]string),
	}
}

type Vault struct {
	encodingKey []byte
	keyValues   map[string]string
	filepath    string
	mutex       sync.Mutex
}

func (v *Vault) loadKeyValues() error {
	f, err := os.Open(v.filepath)
	if err != nil {
		v.keyValues = make(map[string]string)
		return nil
	}
	defer f.Close()

	var sb strings.Builder
	_, err = io.Copy(&sb, f)
	if err != nil {
		return err
	}

	decryptedString, err := decrypt(sb.String(), v.encodingKey)
	if err != nil {
		return err
	}

	r := strings.NewReader(decryptedString)
	dec := json.NewDecoder(r)
	err = dec.Decode(&v.keyValues)
	if err != nil {
		return err
	}

	return nil
}

func (v *Vault) saveKeyValues() error {
	var sb strings.Builder
	enc := json.NewEncoder(&sb)
	err := enc.Encode(v.keyValues)
	if err != nil {
		return err
	}

	encryptedString, err := encrypt([]byte(sb.String()), v.encodingKey)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(v.filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = fmt.Fprint(f, encryptedString)
	return err
}

func (v *Vault) Get(key string) (string, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	err := v.loadKeyValues()
	if err != nil {
		return "", err
	}

	value, ok := v.keyValues[key]
	if !ok {
		return "", errors.New("Secret: no value for that key")
	}

	return string(value), nil
}

func (v *Vault) Set(key, value string) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	err := v.loadKeyValues()
	if err != nil {
		return err
	}

	v.keyValues[key] = value
	err = v.saveKeyValues()
	return err
}
