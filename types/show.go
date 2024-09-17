package types 

import "math/rand" 

type Show struct {
  id int
  Name string
}

func NewShow(id int,name string) Show{
  return Show {id:id,Name:name}
}

func RandomShow(shows []Show) Show {
  ind := rand.Intn(len(shows))
  return shows[ind]
}
