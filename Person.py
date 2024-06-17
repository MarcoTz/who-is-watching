from Show import Show 

class Person:
    name : str
    person_id : int 
    watching : list[Show]

    def __init__(self,name:str,person_id:int) -> None:
        self.name : str = name
        self.person_id : int = person_id 
        self.watching : list[Show] = [] 

    def add_show(self,new_show:Show) -> None:
        if new_show in self.watching: 
            return  
        self.watching.append(new_show)

    def remove_show(self,old_show:Show) -> None:
        self.watching.remove(old_show)
