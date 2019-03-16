package git

import (
	"fmt"
	"io/ioutil"
	"os"
)

type Object interface {
	Kind() string
	Serialize() string
	Deserialize(data string)
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
		obj = NewBlob(string(data))
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
	blobData string
}

func NewBlob(blobData string) *Blob {
	return &Blob{blobData: blobData}
}

func (b *Blob) Kind() string {
	return "blob"
}
func (b *Blob) Serialize() string {
	return b.blobData
}
func (b *Blob) Deserialize(blobData string) {
	b.blobData = blobData
}
