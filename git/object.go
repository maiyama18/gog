package git

type Object interface {
	Serialize()
	Deserialize()
}
