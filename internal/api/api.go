// api.go
package api

import (
	"fmt"

	"github.com/rreubenreyes/goth/internal/db"
	"github.com/rreubenreyes/goth/internal/rbac"
)

// Expand lists who has the given relationship to the given object.
func Expand(subjSet rbac.SubjectSet) (subjects []rbac.Subject) {
	// make sure the object exists
	resource, err := db.GetNode(map[string]string{
		"ID":    fmt.Sprintf("%s:%s", subjSet.Namespace, subjSet.Object),
		"Label": "Object",
	})
	if err != nil {
		panic("object does not exist")
	}

	// list all outgoing relationships to that object
	edges, _ := db.QueryEdges(map[string]string{
		"To":    resource.ID,
		"Label": subjSet.Relation,
	})

	for _, e := range edges {
		subject := rbac.NewSubject(e.From.ID)

		// if subject refers to a subject set, then also expand that subject set
		ss, err := subject.SubjectSet()
		if err == nil {
			return Expand(ss)
		}

		// else this subject is concrete, nothing else to expand
		subjects = append(subjects, subject)
	}

	return
}

// List lists what objects the subject has access to.
//
// Note that this function assumes subject has a concrete ID.
// To query on a SubjectSet, use Expand.
func List(subj string, rel string) (objects []string) {
	// get the target subject from the db
	resource, err := db.GetNode(map[string]string{
		"ID":    subj,
		"Label": "Subject",
	})
	if err != nil {
		panic("subject does not exist")
	}

	// list all relationships from that subject
	edges, _ := db.QueryEdges(map[string]string{
		"From":  resource.ID,
		"Label": rel,
	})

	for _, e := range edges {
		objects = append(objects, e.To.ID)
	}

	return
}

// Check checks if the given subject has access to the given object.
func Check(subj string, subjSet rbac.SubjectSet) bool {
	// make sure the object exists
	resource, err := db.GetNode(map[string]string{
		"ID":    fmt.Sprintf("%s:%s", subjSet.Namespace, subjSet.Object),
		"Label": "Object",
	})
	if err != nil {
		panic("object does not exist")
	}

	// list all outgoing relationships to that object
	edges, _ := db.QueryEdges(map[string]string{
		"To":    resource.ID,
		"Label": subjSet.Relation,
	})

	for _, e := range edges {
		subject := rbac.NewSubject(e.From.ID)

		// if this edge refers to a concrete subject, return immediately
		// if that subject matches the one in question
		id, err := subject.ID()
		if err == nil && id == subj {
			return true
		}

		// else, the edge refers to another subject set, evaluate those
		ss, _ := subject.SubjectSet()
		if Check(subj, ss) {
			return true
		}
	}

	return false
}
