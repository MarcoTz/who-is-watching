package types 

import "fmt"

type WatchGroup struct {
  Id int 
  Watchers []Watcher 
  Show Show
  Current_ep int
}

func DisplayGroup(gr WatchGroup) string {
  group_str := fmt.Sprintf("Group for show : %s, current episode: %d, watchesr:\n", gr.Show.Name, gr.Current_ep)
  for _,watcher := range gr.Watchers{
    group_str += "\n"+watcher.Name
  }
  return group_str
}
