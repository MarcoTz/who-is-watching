package pages 

import ( 
  "fmt"
  "database/sql"
  "rooxo/whoiswatching/database"
  "rooxo/whoiswatching/types"
)

type GroupsPage struct {
  title string
  content string
}

func (p *GroupsPage) get_title() string { return p.title }
func (p *GroupsPage) get_content() string { return p.content}

func RenderGroups(db *sql.DB) string {
  groups, err := database.GetAllGroups(db)
  if err != nil { return "Could not load groups" }

  groups_str := "<ul>\n"
  for _,group := range groups{
    groups_str += fmt.Sprintf("<li>%s</li>\n",types.DisplayGroup(group))
  }
  groups_str += "</ul>"

  page := GroupsPage { title:"Groups", content:groups_str }
  return RenderTemplate(&page) 
}
