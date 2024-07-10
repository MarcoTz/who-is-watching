from common.Person     import Person
from common.Show       import Show 
from common.WatchGroup import WatchGroup
from common.types import * 
from file_io.load_json import load_shows,load_people,load_groups
from file_io.save_json import save_shows,save_people,save_groups

class ShowManager:
    people : list[Person]
    shows  : list[Show]
    groups : list[WatchGroup]

    def __init__(self) -> None:
        self.shows = load_shows()
        self.people = load_people()
        self.groups = load_groups()

    def save_all(self) -> None: 
        save_shows(self.shows)
        save_people(self.people)
        save_groups(self.groups)

    def get_person_by_name(self,person_name:str) -> Person | WatcherError:
        for person in self.people:
            if person.name.lower().strip() == person_name.lower().strip():
                return person
        return person_not_found(person_name) 

    def get_show_by_name(self,show_name:str) -> Show | WatcherError:
        for show in self.shows: 
            if show.name.lower().strip() == show_name.lower().strip():
                return show
        return show_not_found(show_name)

    def get_show_by_id(self,show_id:int) -> Show | WatcherError:
        for show in self.shows:
            if show.show_id == show_id:
                return show
        return show_not_found('id: %s' % str(show_id)) 

    def get_watchgroups_by_show_id(self,show_id:int) -> list[WatchGroup] | WatcherError:
        show_exists : WatcherError | None = is_watcher_error(self.get_show_by_id(show_id))
        if show_exists is not None:
            return show_exists

        group_list : list[WatchGroup] = []
        for watch_group in self.groups:
            if watch_group.show_id == show_id:
                group_list.append(watch_group)

        return group_list

    def get_watchgroups_by_people(self,watchers:list[Person]) -> list[WatchGroup]:
        people_ids : set[int] = set(map(lambda x: x.person_id,watchers))
        group_list : list[WatchGroup] = []
        for watch_group in self.groups:
            if set(watch_group.people_ids) == people_ids:
                group_list.append(watch_group)
        return group_list

    def get_watchgroups_by_person(self,watcher:Person) -> list[WatchGroup]:
        group_list : list[WatchGroup] = []
        for watch_group in self.groups:
            if watcher.person_id in watch_group.people_ids:
                group_list.append(watch_group)
        return group_list


    def get_next_show_id(self) -> int:
        next_id : int = 0
        show_ids : list[int] = list(map(lambda x:x.show_id,self.shows))
        while next_id in show_ids:
            next_id += 1
        return next_id

    def get_possible(self,names_list:list[str]) -> list[Show] | WatcherError:
        watcher_list : list[Person] = [] 
        for name in names_list:
            person : Person | WatcherError = self.get_person_by_name(name)
            match person:
                case Person():
                    watcher_list.append(person)
                case WatcherError():
                    return person


        groups : list[WatchGroup] = self.get_watchgroups_by_people(watcher_list)
        show_list : list[Show] = [] 
        for watch_group in groups:
            show_id : int = watch_group.show_id
            show : Show | WatcherError = self.get_show_by_id(show_id)
            match show:
                case Show():
                    show_list.append(show)
                case WatcherError():
                    return show

        return show_list

    def add_watcher_show(self,watcher_name:str,show_name:str) -> WatchGroup | WatcherError:
        person : Person | WatcherError = self.get_person_by_name(watcher_name)
        match person:
            case WatcherError():
                return person

        show : Show | WatcherError = self.get_show_by_name(show_name)
        match show:
            case WatcherError():
                return show
        
        watch_groups : list[WatchGroup] | WatcherError = self.get_watchgroups_by_show_id(show.show_id)
        match watch_groups:
            case WatcherError():
                return watch_groups

        for watch_group in watch_groups:
            watch_group.people_ids.append(person.person_id)

        self.save_all()
        return watch_groups[0]

    def remove_watcher_show(self,watcher_name:str,show_name:str) -> None | WatcherError:
        person : Person | WatcherError = self.get_person_by_name(watcher_name)
        match person:
            case WatcherError():
                return person

        show : Show | WatcherError = self.get_show_by_name(show_name)
        match show:
            case WatcherError():
                return show
        watch_groups : list[WatchGroup] | WatcherError = self.get_watchgroups_by_show_id(show.show_id)
        match watch_groups: 
            case WatcherError():
                return watch_groups

        for watch_group in watch_groups:
            watch_group.people_ids.remove(person.person_id)
        self.save_all()

    def update_show_episode(self,show_name:str,ep_nr:int) -> None | WatcherError:
        show : Show | WatcherError = self.get_show_by_name(show_name)
        match show:
            case WatcherError():
                return show

        watch_groups : list[WatchGroup] | WatcherError = self.get_watchgroups_by_show_id(show.show_id)
        match watch_groups:
            case WatcherError():
                return watch_groups 

        for watch_group in watch_groups:
            watch_group.episode_nr = ep_nr
        self.save_all()

    def add_show(self,show_name:str) -> None | WatcherError:
        show : Show | WatcherError = self.get_show_by_name(show_name)
        match show: 
            case Show():
                return show_exists(show_name)

        show_id : int = self.get_next_show_id()
        new_show_info : ShowInfo = { 
                                    'show_id':show_id,
                                    'show_name':show_name
                                    }
        new_show : Show = Show(new_show_info)
        self.shows.append(new_show)
        self.save_all()

    def get_shows_person(self,person_name:str) -> list[Show] | WatcherError:
        shows_list : list[Show] = []
        person : Person | WatcherError = self.get_person_by_name(person_name)
        match person:
            case WatcherError():
                return person

        
        person_groups : list[WatchGroup] = self.get_watchgroups_by_person(person)
        shows_list : list[Show] = []
        for watch_group in person_groups:
            show : Show | WatcherError = self.get_show_by_id(watch_group.show_id)
            match show:
                case WatcherError():
                    return show
            shows_list.append(show)

        return shows_list

    def remove_show(self,show_name:str) -> None | WatcherError:
        show : Show | WatcherError = self.get_show_by_name(show_name)
        match show:
            case WatcherError():
                return show

        show_groups : list[WatchGroup] | WatcherError = self.get_watchgroups_by_show_id(show.show_id)
        match show_groups:
            case WatcherError():
                return show_groups

        for show_group in show_groups:
            self.groups.remove(show_group)

        self.shows.remove(show)
        self.save_all()

