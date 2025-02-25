package main

import (
	"fmt"
	"strings"

	link "github.com/Kacperbrozyna/learning_go_repo/http_link_parser/link_parser"
)

const example_html = `
<html>
<head>
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css">
</head>
<body>
  <h1>Social stuffs</h1>
  <div>
    <a href="https://www.linkedin.com/in/kacper-brozyna/">
      Check me out on <strong>LinkedIn</strong>!
    </a>
    <a href="https://github.com/Kacperbrozyna">
      Check me out on <strong>Github</strong>!
    </a>
  </div>
</body>
</html>
`

func main() {
	string_reader := strings.NewReader(example_html)
	links, err := link.Parse(string_reader)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", links)
}
