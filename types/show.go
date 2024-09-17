package types 

import "math/rand" 

type Show struct {
  Id int
  Name string
}

func RandomShow(shows []Show) Show {
  ind := rand.Intn(len(shows))
  return shows[ind]
}
