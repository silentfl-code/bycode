package main

import (
	"fmt"
	"github.com/go-martini/martini"
)

func main() {
	m := martini.Classic()
	m.Get("/:id", func(params martini.Params) string {
		return fmt.Sprintf(`
		<html>
			<head>
				<title>%s</title>
			</head>
			<body>
				Lorem ipsum
			</body>
		</html>`,
			params["id"])
	})
	m.Run()
}
