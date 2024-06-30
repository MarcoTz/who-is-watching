from common.types import * 
class WatchGroup:
    show_id : int 
    people_ids : list[int]
    episode_nr : int 

    def __init__(self,json_dict:GroupInfo):
        self.show_id    : int       = json_dict['show_id']
        self.people_ids : list[int] = json_dict['people_ids']
        self.episode_nr : int       = json_dict['episode_nr']

    def to_json_dict(self) -> GroupInfo:
        return {
                'show_id':self.show_id,
                'people_ids':self.people_ids,
                'episode_nr':self.episode_nr
                }
