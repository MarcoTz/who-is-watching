from common.types import *
class Person:
    name      : str
    person_id : int 

    def __init__(self,info_dict:PersonInfo) -> None:
        self.name : str = info_dict['person_name']
        self.person_id : int = info_dict['person_id']

    def get_json_dict(self) -> PersonInfo:
        return {
                'person_id' : self.person_id,
                'person_name': self.name
        }
