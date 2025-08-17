package pages 

import ( 
  "fmt"
  "database/sql"
  "rooxo/whoiswatching/database"
)

type WatchersPage struct{
  title string
  content string
} 

func (p *WatchersPage) get_title() string { return p.title }
func (p *WatchersPage) get_content() string { return p.content }

func RenderWatchers(db *sql.DB) string {
  watchers, err := database.GetAllWatchers(db)
  if err != nil { return "Could not load watchers" }

  watchers_str := "<ul>\n"
  for _,watcher := range watchers {
    watchers_str += fmt.Sprintf("<li>%s</li>\n",watcher.Name)
  }
  watchers_str += "</ul>"

  page := WatchersPage { title: "Watchers", content:watchers_str}
  return RenderTemplate(&page)
}
