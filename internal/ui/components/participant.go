package components

type Participant struct {
	Name    string
	IsPrime bool
}

func (p Participant) Title() string {
	return p.Name
}

func (p Participant) Description() string {
	return ""
}

func (p Participant) FilterValue() string {
	return p.Name
}
