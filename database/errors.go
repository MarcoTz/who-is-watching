package database

import "fmt"

type ShowIdNotFoundErr struct {
  show_id int
}

func (e *ShowIdNotFoundErr) Error() string {
  return fmt.Sprintf("Could not find show with id %d",e.show_id)
}

type UserIdNotFoundErr struct {
  user_id int 
}

func (e *UserIdNotFoundErr) Error() string {
  return fmt.Sprintf("Could not find user with id %d",e.user_id)
}
