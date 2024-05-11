package db

// table
type Node struct {
	ID    string // pk component
	Label string // pk component
}

// table
// indexes:
//   - `Label#From`
//   - `Label#To`
type Edge struct {
	From  *Node  // auxiliary field
	To    *Node  // auxiliary field
	ID    string // pk - `Label#From#To`
	Label string // auxiliary field
}

// GetItem queries for a Node from the database.
func GetNode(q map[string]string) (*Node, error) {
	// need a db adapter or something
	return nil, nil
}

// GetItem queries for an Edge from the database.
func QueryEdges(q map[string]string) ([]*Edge, error) {
	// need a db adapter or something
	return nil, nil
}
