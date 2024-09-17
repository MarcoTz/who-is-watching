package types 

type Show struct {
  id int
  name string
}

func NewShow(id int,name string) *Show{
  return &Show {id:id,name:name}
}
