from typing import TypedDict 
from enum import Enum 

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
    group_id   : int
    show_id    : int
    people_ids : list[int]
    episode_nr : int

def coalesce_group_info(json_dict) -> GroupInfo:
    people_ids : list[int] = list(map(lambda x:int(x),list(json_dict['people_ids'])))
    info_dict : GroupInfo = {
            'group_id'   : int(json_dict['group_id']),
            'show_id'    : int(json_dict['show_id']),
            'people_ids' : people_ids,
            'episode_nr' : int(json_dict['episode_nr'])
            }
    return info_dict

class ErrorType(Enum):
    PERSON_NOT_FOUND = 1
    SHOW_NOT_FOUND = 2
    GROUP_NOT_FOUND = 3 
    SHOW_EXISTS = 4 
    NO_GROUP_SHOW = 5 
    PERSON_NOT_IN_GROUP = 6

all_error_types : list[ErrorType] = [ErrorType.PERSON_NOT_FOUND, 
                             ErrorType.SHOW_NOT_FOUND,
                             ErrorType.GROUP_NOT_FOUND,
                             ErrorType.SHOW_EXISTS,
                             ErrorType.NO_GROUP_SHOW,
                            ErrorType.PERSON_NOT_IN_GROUP]

class WatcherError(TypedDict):
    error_type : ErrorType
    error_message : str

def show_err(err:WatcherError) -> str:
    match err['error_type']:
        case ErrorType.PERSON_NOT_FOUND:
            return 'Could not find person\n%s' % err['error_message']
        case ErrorType.SHOW_NOT_FOUND:
            return 'Could not find show\n%s' % err['error_message']
        case ErrorType.GROUP_NOT_FOUND:
            return 'Could not find watch group\n%s' % err['error_message']
        case ErrorType.SHOW_EXISTS:
            return 'Show already exists\n%s' % err['error_message']
        case ErrorType.NO_GROUP_SHOW:
            return 'No group for show\n%s' % err['error_message']
        case ErrorType.PERSON_NOT_IN_GROUP:
            return 'Person not found in watchgroup\n%s' % err['error_message']

def person_not_found(msg:str) -> WatcherError:
    return {
      'error_type':ErrorType.PERSON_NOT_FOUND,
      'error_message':msg
      }

def show_not_found(msg:str) -> WatcherError:
    return {
      'error_type':ErrorType.SHOW_NOT_FOUND,
      'error_message':msg
    }

def group_not_found(msg:str) -> WatcherError:
    return {
      'error_type':ErrorType.GROUP_NOT_FOUND,
      'error_message':msg
      }

def no_show_for_group(msg:str) -> WatcherError:
    return {
      'error_type':ErrorType.NO_GROUP_SHOW,
      'error_message':msg
    }

def show_exists(msg:str) -> WatcherError:
    return {
      'error_type':ErrorType.SHOW_EXISTS,
      'error_message':msg
      }
def not_in_watchgroup(msg:str) -> WatcherError:
    return {
      'error_type':ErrorType.PERSON_NOT_IN_GROUP,
      'error_message':msg
      }

def is_watcher_error(maybe_err) -> WatcherError | None:
    if 'error_type' not in maybe_err:
        return None 
    if ('error_message' in maybe_err) and maybe_err['error_type'] in all_error_types:
        return maybe_err
    return None

