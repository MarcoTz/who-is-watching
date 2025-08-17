package pages

import (
  "fmt"
  "database/sql"
  "rooxo/whoiswatching/database"
)

type ShowsPage struct {
  title string
  content string
}
func (p *ShowsPage) get_title() string { return p.title }
func (p *ShowsPage) get_content() string {return p.content} 

func RenderShows(db *sql.DB) string {
  shows, err := database.GetAllShows(db);
  if err != nil { return "Could not load shows" }

  shows_str := "<ul>\n"
  for _,show := range(shows){
    shows_str += fmt.Sprintf("<li>%s</li>\n",show.Name)
  }
  shows_str += "</ul>"

  page := ShowsPage { title: "Shows", content: shows_str}
  return RenderTemplate(&page)

}
