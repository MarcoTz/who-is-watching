package database

import "fmt"

type ShowIdDoesNotExistErr struct { show_id int }
func (e *ShowIdDoesNotExistErr) Error() string { return fmt.Sprintf("Show with id %d does not exist",e.show_id) }

type ShowNameDoesNotExistErr struct { show_name string }  
func (e *ShowNameDoesNotExistErr) Error() string { return fmt.Sprintf("Show %s does not exist",e.show_name) }

type ShowExistsErr struct {show_name string}
func (e *ShowExistsErr) Error() string { return fmt.Sprintf("Show %s already exists",e.show_name) }

type WatcherExistsErr struct { watcher_name string }
func (e *WatcherExistsErr) Error() string{ return fmt.Sprintf("Watcher %s already exist",e.watcher_name) }

type WatcherNameDoesNotExistErr struct {watcher_name string }
func (e *WatcherNameDoesNotExistErr) Error() string{ return fmt.Sprintf("Watcher %s does not exist",e.watcher_name) }

type WatcherIdDoesNotExistErr struct { watcher_id int }
func (e *WatcherIdDoesNotExistErr) Error() string{ return fmt.Sprintf("Watcher with id %d does not exist", e.watcher_id) }

type GroupIdDoesNotExistErr struct {group_id int}
func (e *GroupIdDoesNotExistErr) Error() string { return fmt.Sprintf("Group with id %d does not exist",e.group_id) }

type NotAWatcherErr struct { 
  watcher_name string 
  group_id int
}
func (e *NotAWatcherErr) Error() string { return fmt.Sprintf("%s is not in group %d",e.watcher_name,e.group_id)} 
