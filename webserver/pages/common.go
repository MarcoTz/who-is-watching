package pages 

import (
  "fmt"
)

type Page interface {
  get_title() string
  get_content() string
}

const PAGE_TEMPLATE = `
  <!doctype html>
  <html>
    <head><title>%s</title></head>
    <body>
      %s
      %s
    </body>
  </html>`

const HEADER_TEMPLATE = `
  <div id="header">
    <a href="shows">Shows</a>
    <a href="watchers">Watchers</a>
    <a href="groups">Groups</a>
  </div>
  `

func RenderTemplate(page Page) string{
  return fmt.Sprintf(PAGE_TEMPLATE,
    page.get_title(),
    HEADER_TEMPLATE,
    page.get_content())
}
