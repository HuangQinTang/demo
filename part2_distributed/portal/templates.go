package portal

import (
	"html/template"
)

var rootTemplate *template.Template

func ImportTemplates() error {
	var err error
	rootTemplate, err = template.ParseFiles(
		"./part2_distributed/portal/students.html",
		"./part2_distributed/portal/student.html")

	if err != nil {
		return err
	}

	return nil
}
