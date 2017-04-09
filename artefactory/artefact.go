package artefactory

type Artefact struct {
	Project     string
	Branch      string
	Commit      string
	Stage       string
	Job         string
	Success     bool
	Description string
}
