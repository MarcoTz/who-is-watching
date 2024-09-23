package types 

import ( 
  "fmt"
  "strings"
)

type WatchGroup struct {
  Id int 
  Watchers []Watcher 
  Show Show
  Current_ep int
  Done bool
}

func DisplayGroup(gr WatchGroup) string {
  watcher_names := make([]string,0)
  for _,watcher := range gr.Watchers {
    watcher_names = append(watcher_names,watcher.Name)
  }
  watcher_str := strings.Join(watcher_names,", ")

  group_str := fmt.Sprintf(`
  Group (ID %d) 
  show: %s, done: %t
    current episode: %d
    watchers: %s`, 
  gr.Id, gr.Show.Name, gr.Done, gr.Current_ep,watcher_str)
  return group_str
}

func (g *WatchGroup) GetWatchers() []string{
  watcher_names := make([]string,0)
  for _,watcher := range g.Watchers{
    watcher_names = append(watcher_names,watcher.Name)
  }
  return watcher_names
}
