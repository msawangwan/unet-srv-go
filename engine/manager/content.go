package manager

type ContentHandler struct {
	*LoreGenerator
	*NameGenerator
}

func NewContentHandler(filepath *string) (*ContentHandler, error) {
	var p string

	if filepath == nil {
		p = "../../lore.json" // TODO: load from config/params.json
	} else {
		p = *filepath
	}

	lg, e := NewLoreGenerator(p, nil)
	if e != nil {
		return nil, e
	}

	ng, e := NewNameGenerator()
	if e != nil {
		return nil, e
	}

	ch := &ContentHandler{
		LoreGenerator: lg,
		NameGenerator: ng,
	}

	return ch, nil
}
