from typing import TypedDict 

class ShowInfo(TypedDict):
    show_id   : int
    show_name : str 

def coalesce_show_info(json_dict) -> ShowInfo:
    info_dict : ShowInfo = {
            'show_id': int(json_dict['id']),
            'show_name' : str(json_dict['name']).strip() 
    }
    return info_dict


class PersonInfo(TypedDict):
    person_id   : int 
    person_name : str

def coalesce_person_info(json_dict) -> PersonInfo:
    info_dict : PersonInfo = {
            'person_id': int(json_dict['id']),
            'person_name': str(json_dict['name'])
            }
    return info_dict


class GroupInfo(TypedDict):
    show_id : int
    people_ids : list[int]
    episode_nr : int

def coalesce_group_info(json_dict) -> GroupInfo:
    people_ids : list[int] = list(map(lambda x:int(x),list(json_dict['people_ids'])))
    info_dict : GroupInfo = {
            'show_id'    : int(json_dict['show_id']),
            'people_ids' : people_ids,
            'episode_nr' : int(json_dict['episode_nr'])
            }
    return info_dict

