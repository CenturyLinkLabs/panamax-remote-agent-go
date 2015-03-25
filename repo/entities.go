package repo

type Deployment struct {
	ID         int
	Name       string
	ServiceIDs string
	Template   string
}

type Persister interface {
	FindByID(int) (Deployment, error)
	All() ([]Deployment, error)
	Save(*Deployment) error
	Remove(int) error
}
