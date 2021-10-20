package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

type JSONStore struct {
	mu       sync.Mutex
	filePath string

	kv map[string]json.RawMessage
}

func OpenJSONStore(filePath string) (*JSONStore, error) {
	data, ferr := ioutil.ReadFile(filePath)

	var kv map[string]json.RawMessage

	if ferr != nil {
		if os.IsNotExist(ferr) {
			kv = make(map[string]json.RawMessage)
		} else {
			return nil, ferr
		}
	} else {
		if len(data) == 0 {
			// accept empty file
			kv = make(map[string]json.RawMessage)
		} else {
			err := json.Unmarshal(data, &kv)
			if err != nil {
				return nil, err
			}
		}
	}

	s := JSONStore{
		filePath: filePath,
		kv:       kv,
	}

	return &s, nil
}

func (s *JSONStore) Get(key string, val interface{}) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, ok := s.kv[key]
	if !ok {
		return false, nil
	}

	err := json.Unmarshal(data, val)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *JSONStore) Set(key string, val interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	s.kv[key] = data
	return s.saveFile()
}

func (s *JSONStore) saveFile() error {
	f, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")

	return enc.Encode(s.kv)

}
