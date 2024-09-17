package database

import "fmt"

type ShowIdDoesNotExist struct { show_id int }
func (e *ShowIdDoesNotExist) Error() string { return fmt.Sprintf("Show with id %d does not exist",e.show_id) }

type ShowNameDoesNotExist struct { show_name string }  
func (e *ShowNameDoesNotExist) Error() string { return fmt.Sprintf("Show %s does not exist",e.show_name) }

type WatcherExistsErr struct { watcher_name string }
func (e *WatcherExistsErr) Error() string{ return fmt.Sprintf("Watcher %s already exist",e.watcher_name) }

type WatcherNameDoesNotExistErr struct {watcher_name string }
func (e *WatcherNameDoesNotExistErr) Error() string{ return fmt.Sprintf("Watcher %s does not exist",e.watcher_name) }

type WatcherIdDoesNotExistErr struct { watcher_id int }
func (e *WatcherIdDoesNotExistErr) Error() string{ return fmt.Sprintf("Watcher with id %d does not exist", e.watcher_id) }
