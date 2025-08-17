package pages

import ( 
  "rooxo/whoiswatching/anilist_api"
)
type IndexPage struct{
  title string 
  content string 
}

func (p *IndexPage) get_title() string { return p.title }
func (p *IndexPage) get_content() string { return p.content}

func RenderIndex() string{
  page := IndexPage { title:"Index", content:anilist_api.QueryApi()}
  return RenderTemplate(&page)
}
