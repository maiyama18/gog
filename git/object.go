package git

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Object interface {
	Kind() string
	Serialize() []byte
	Deserialize(data []byte)
}

func ObjectHash(filePath, kind string, dryRun bool) (string, error) {
	repo, err := FindRepository(filePath, true)
	if err != nil {
		return "", err
	}

	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	var obj Object
	switch kind {
	case "commit":
	case "tree":
	case "tag":
	case "blob":
		obj = NewBlob(data)
	default:
		return "", fmt.Errorf("unknown type: %s", kind)
	}

	sha, err := repo.WriteObject(obj, dryRun)
	if err != nil {
		return "", err
	}
	return sha, nil
}

type Blob struct {
	data []byte
}

func NewBlob(data []byte) *Blob {
	return &Blob{data: data}
}

func (b *Blob) Kind() string {
	return "blob"
}
func (b *Blob) Serialize() []byte {
	return b.data
}
func (b *Blob) Deserialize(data []byte) {
	b.data = data
}

type Commit struct {
	kvlm *KVLM
}

func (c *Commit) Kind() string {
	return "commit"
}
func (c *Commit) Serialize() []byte {
	return c.kvlm.Serialize()
}
func (c *Commit) Deserialize(data []byte) {
	kvlm := NewKVLM()
	_ = ParseKVLM(data, 0, kvlm)
	c.kvlm = kvlm
}

type KVLM struct {
	keys      []string
	keyValues map[string][]string
}

func NewKVLM() *KVLM {
	return &KVLM{
		keys:      make([]string, 0),
		keyValues: make(map[string][]string),
	}
}

func ParseKVLM(raw []byte, start int, kvlm *KVLM) error {
	spi := bytes.Index(raw, []byte(" "))
	nli := bytes.Index(raw, []byte("\n"))

	// blank line followed by commit message
	if spi < 0 || nli < spi {
		if nli != start {
			return fmt.Errorf("invalid format of kvlm: %v", raw)
		}
		// key for commit message is blank
		kvlm.add("", string(raw[start+1:]))
		return nil
	}

	// usual key-value pair
	key := string(raw[start:spi])
	end := start
	for {
		end = bytes.Index(raw, []byte("\n"))
		if raw[end+1] == ' ' {
			break
		}
	}
	value := strings.Replace(string(raw[spi+1:end]), "\n ", "\n", -1)
	kvlm.add(key, value)

	return ParseKVLM(raw, end+1, kvlm)
}

func (k *KVLM) add(key, value string) {
	k.keys = append(k.keys, key)
	k.keyValues[key] = append(k.keyValues[key], value)
}

func (k *KVLM) Serialize() []byte {
	var out bytes.Buffer
	for _, key := range k.keys {
		values := k.keyValues[key]
		for _, value := range values {
			out.WriteString(key)
			out.WriteString(" ")
			out.WriteString(strings.Replace(value, "\n", "\n ", -1))
			out.WriteString("\n")
		}
	}
	return out.Bytes()
}
