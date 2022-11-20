package portal

import (
	"fmt"
	"html/template"
	"os"
)

var rootTemplate *template.Template

func ImportTemplates() error {
	var err error
	str, _ := os.Getwd()
	rootTemplate, err = template.ParseFiles(
		fmt.Sprintf("%s%s", str, "/portal/students.html"),
		fmt.Sprintf("%s%s", str, "/portal/student.html"))

	if err != nil {
		return err
	}

	return nil
}
