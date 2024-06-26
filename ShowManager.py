from Person import Person
from Show import Show 
import load_csv
import save_csv

class ShowManager:
    people : list[Person]
    shows : list[Show]
    shows_file_name : str 
    watchers_file_name : str


    def __init__(self) -> None:
        self.people : list[Person] = []
        self.shows : list[Show] = []
        self.shows_file_name    : str = 'shows.csv'
        self.watchers_file_name : str = 'shows_people.csv'


    def get_next_show_id(self) -> int:
        show_ids : list[int] = list(map(lambda x:x.show_id,self.shows))
        return self.get_next_id(show_ids)

    def get_next_person_id(self) -> int:
        people_ids : list[int] = list(map(lambda x:x.person_id,self.people))
        return self.get_next_id(people_ids)

    def get_next_id(self,used_ids:list[int]) -> int:
        new_id : int = 0 
        while new_id in used_ids:
            new_id += 1
        return new_id

    def load_from_csv(self) -> None:
        show_list = load_csv.load_shows(self.shows_file_name)
        self.load_shows(show_list)
        watchers_list = load_csv.load_watchers(self.watchers_file_name)
        self.load_watchers(watchers_list)

    def load_shows(self,shows : list[tuple[str,int]]) -> None:
        for (show_name,ep_nr) in shows:
            if self.get_show(show_name=show_name) is not None:
                continue
            show_id : int = self.get_next_show_id()
            new_show : Show = Show(name=show_name,show_id=show_id)
            new_show.update_ep(ep_nr)
            self.shows.append(new_show)

    def load_watchers(self,watchers:dict[str,list[str]]) -> None:
        for person_name in watchers.keys():
            person : Person = self.get_person_add(person_name)

            for show_name in watchers[person_name]:
                show : Show | None = self.get_show(show_name)
                if show is None:
                    continue 
                person.add_show(show)

    def save_to_csv(self) -> None:
        shows_list : list[tuple[str,str]] = []
        for show in self.shows:
            shows_list.append((show.name,str(show.current_ep)))
        save_csv.save_shows(shows_list,self.shows_file_name)
        watchers_list : list[tuple[str,str]] = []
        for person in self.people:
            for show in person.watching:
                watchers_list.append((person.name,show.name))
        save_csv.save_watchers(watchers_list,self.watchers_file_name)
    
    def get_person(self,person_name:str) -> Person | None:
        for person in self.people:
            if person.name == person_name:
                return person
        return None

    def get_person_add(self,person_name:str) -> Person:
        mPerson : Person | None = self.get_person(person_name)
        if mPerson is not None:
            return mPerson 
        new_id : int = self.get_next_person_id()
        new_person : Person = Person(name=person_name,person_id=new_id)
        self.people.append(new_person)
        self.save_to_csv()
        return new_person

    def get_show(self,show_name:str) -> Show | None:
        for show in self.shows:
            if show_name == show.name:
                return show 
        return None

    def get_show_add(self,show_name:str) -> Show:
        show : Show | None = self.get_show(show_name)
        if show is not None:
            return show
        new_id : int = self.get_next_show_id()
        new_show : Show = Show(show_name,new_id)
        self.shows.append(new_show)
        self.save_to_csv()
        return new_show

    def remove_show(self,show_name:str) -> None:
        show : Show | None = self.get_show(show_name)
        if show is None:
            return 
        self.shows.remove(show)
        self.save_to_csv()

    def add_watcher(self,person_name:str,show_name:str) -> None:
        show    : Show   | None = self.get_show_add(show_name=show_name)
        person  : Person | None = self.get_person_add(person_name=person_name)

        if show is None or person is None:
            return 
        
        person.add_show(new_show = show)
        self.save_to_csv()

    def remove_watcher(self,person_name:str,show_name:str) -> None:
        show : Show | None = self.get_show_add(show_name=show_name)
        person : Person | None = self.get_person_add(person_name=person_name)
        
        if show is None or person is None:
            return 

        person.remove_show(old_show=show)
        self.save_to_csv()


    def get_watchers(self,show_name:str) -> list[Person] :
        watchers : list[Person] = [] 
        show : Show | None = self.get_show(show_name=show_name)

        for person in self.people:
            if show in person.watching:
                watchers.append(person)

        return watchers 

    def get_possible_shows(self,watchers:list[Person]) -> list[Show] :
         possible_shows : list[Show] = self.shows

         for person in watchers:
             possible_shows = list(filter(lambda x: x in person.watching,possible_shows))

         return possible_shows

    def get_possible_shows_names(self,watchers:list[str]) -> list[Show]:
        people_list : list[Person] = []
        for watcher_name in watchers:
            watcher : Person | None = self.get_person(watcher_name)
            if watcher is None:
                continue 
            people_list.append(watcher)
        return self.get_possible_shows(watchers=people_list)
