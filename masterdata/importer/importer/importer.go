package importer

type Importer struct {
	path string
}

func NewImporter(path string) *Importer {
	return &Importer{
		path: path,
	}
}

func (i *Importer) Run() error {
	return nil
}
