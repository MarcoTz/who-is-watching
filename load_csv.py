import csv 

def load_shows(shows_file_name:str) -> list[tuple[str,int]]:
    with open(shows_file_name) as shows_file:
        shows_reader = csv.reader(shows_file,delimiter=';')
    
        show_list : list[tuple[str,int]] = []
        for show_row in shows_reader:
            show_name = show_row[0].strip()
            if show_name =='Show':
                continue
    
            ep_nr : int = int(show_row[1].strip())
    
            show_list.append((show_name,ep_nr))
    
        shows_file.close()
    
        return show_list

def load_watchers(watchers_file_name:str) -> dict[str,list[str]]:
    with open(watchers_file_name) as watchers_file:
        watchers_reader = csv.reader(watchers_file,delimiter=';')
        
        watchers_dict : dict[str,list[str]] = {}
        for watcher_row in watchers_reader:
            watcher_name : str = watcher_row[0].strip()
    
            if watcher_name == 'Person':
                continue 
            
            show_name : str = watcher_row[1].strip()
            if watcher_name in watchers_dict:
                watchers_dict[watcher_name].append(show_name)
            else:
                watchers_dict[watcher_name] = [show_name]
    
        watchers_file.close() 
        return watchers_dict 
