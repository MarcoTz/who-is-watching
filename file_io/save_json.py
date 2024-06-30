from common.Show import Show 
from common.Person import Person
from common.WatchGroup import WatchGroup
from common.types import *
from common.constants import * 

import json 
import os

def save_json(file_path:str, json_data) -> None:
    with open(file_path,'w+') as json_file:
        json_file.write(json.dumps(json_data))

def save_shows(shows:list[Show]) -> None:
    json_list : list[ShowInfo] = list(map(lambda x:x.get_json_dict(),shows))
    shows_file_path : str = os.path.join(data_dir,shows_json)
    save_json(shows_file_path,json_list)

def save_people(people:list[Person]) -> None:
    json_list : list[PersonInfo] = list(map(lambda x: x.get_json_dict(),people))
    people_file_path : str = os.path.join(data_dir,people_json)
    save_json(people_file_path,json_list)

def save_groups(groups:list[WatchGroup]) -> None:
    group_list : list[GroupInfo] = list(map(lambda x: x.get_json_dict(),groups))
    group_file_path : str = os.path.join(data_dir,watchgroups_json)
    save_json(group_file_path,group_list)

