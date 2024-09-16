from common.types import ShowInfo

class Show:
    name    : str
    show_id : int 

    def __init__(self,info_dict:ShowInfo) -> None:
        self.name    : str = info_dict['show_name']
        self.show_id : int = info_dict['show_id']

    def get_json_dict(self) -> ShowInfo:
        return {
                'show_name':self.name,
                'show_id':self.show_id            
                }
