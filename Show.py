class Show:
    name : str
    show_id : int 
    current_ep : int

    def __init__(self,name:str,show_id:int) -> None:
        self.name = name
        self.current_ep = 0
        self.show_id = show_id

    def update_ep(self,new_ep:int) -> None:
        self.current_ep : int  = new_ep
