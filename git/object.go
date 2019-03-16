package git

type Object interface {
	Kind() string
	Serialize() string
	Deserialize(data string)
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
