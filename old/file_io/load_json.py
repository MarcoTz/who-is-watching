from common.types import *
from common.constants import * 
from common.Show import Show 
from common.Person import Person
from common.WatchGroup import WatchGroup

import json
import os 

def load_json(file_path:str): 
    with open(file_path,'r') as json_file:
        file_contents : str = json_file.read()
        return json.loads(file_contents)

def load_shows() -> list[Show]:
    shows_path : str = os.path.join(data_dir,shows_json)
    shows_data = load_json(shows_path)
    shows_info_list : list[ShowInfo] = list(map(lambda x:coalesce_show_info(x),shows_data))
    shows : list[Show] = list(map(lambda x:Show(x),shows_info_list))
    return shows

def load_people() -> list[Person]:
    people_path : str = os.path.join(data_dir,people_json)
    people_data = load_json(people_path)
    people_info_list : list[PersonInfo] = list(map(lambda x: coalesce_person_info(x),people_data))
    people : list[Person] = list(map(lambda x:Person(x),people_info_list))
    return people

def load_groups() -> list[WatchGroup]:
    groups_path : str = os.path.join(data_dir,watchgroups_json)
    groups_data = load_json(groups_path)
    groups_info_list : list[GroupInfo] = list(map(lambda x: coalesce_group_info(x),groups_data))
    groups : list[WatchGroup] = list(map(lambda x:WatchGroup(x),groups_info_list))
    return groups
