package pages

type IndexPage struct{
  title string 
  content string 
}

func (p *IndexPage) get_title() string { return p.title }
func (p *IndexPage) get_content() string { return p.content}

func RenderIndex() string{
  page := IndexPage { title:"Index", content:""}
  return RenderTemplate(&page)
}
