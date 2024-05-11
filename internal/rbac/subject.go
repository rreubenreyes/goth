package rbac

import (
	"errors"
	"regexp"
)

var subjectIndirectionExpr = regexp.MustCompile("(?P<Namespace>.*):(?P<Object>.*)#(?P<Relation>.*)")

type Subject struct {
	id string
}

type SubjectSet struct {
	Namespace string
	Object    string
	Relation  string
}

func NewSubject(idOrSubjectSet string) Subject {
	return Subject{
		id: idOrSubjectSet,
	}
}

func (s Subject) ID() (string, error) {
	_, err := s.SubjectSet()
	if err != nil {
		return "", errors.New("subject is subject set")
	}

	return s.id, nil
}

func (s Subject) SubjectSet() (SubjectSet, error) {
	match := subjectIndirectionExpr.FindStringSubmatch(s.id)
	matches := make(map[string]string)
	for i, name := range subjectIndirectionExpr.SubexpNames() {
		if i != 0 && name != "" {
			matches[name] = match[i]
		}
	}

	if len(matches) != 3 {
		return SubjectSet{}, errors.New("subject is concrete")
	}

	return SubjectSet{
		Namespace: matches["Namespace"],
		Object:    matches["Object"],
		Relation:  matches["Relation"],
	}, nil
}
