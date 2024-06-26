import csv
def save_shows(shows:list[tuple[str,str]],shows_file_name:str) -> None:
    with open(shows_file_name,'w+') as shows_file:
        header_row : tuple[str,str] = ('Show','Episode')
        writer = csv.writer(shows_file,delimiter=';')
        writer.writerow(header_row)
        for show in shows:
            writer.writerow(show)
        shows_file.close()

def save_watchers(watchers:list[tuple[str,str]],watchers_file_name:str) -> None:
    with open(watchers_file_name,'w+') as watchers_file:
        header_row : tuple[str,str] = ('Person','Show')
        writer = csv.writer(watchers_file,delimiter=';')
        writer.writerow(header_row)
        for watcher in watchers:
            writer.writerow(watcher)
        watchers_file.close()
