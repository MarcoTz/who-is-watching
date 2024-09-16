import jinja2
import os 

from file_io.load_json import load_shows,load_people,load_groups
from common.Show import Show 
from common.Person import Person
from common.WatchGroup import WatchGroup

html_dir : str = 'overview_html'
index_template_path : str = 'index_template.html'
index_path : str = 'index.html'

def get_jinja_template() -> jinja2.Template:

    env : jinja2.Environment = jinja2.Environment(loader=jinja2.FileSystemLoader(html_dir),autoescape=False)
    index_template : jinja2.Template = env.get_template(index_template_path)
    return index_template

def get_show_group(group:WatchGroup, shows_list:list[Show]) -> Show | None: 
    for show in shows_list:
        if show.show_id == group.show_id:
            return show
    return None

def get_people_group(group:WatchGroup,people_list:list[Person]) -> list[Person]:
    watchers : list[Person] = []
    for person in people_list:
        if person.person_id in group.people_ids:
            watchers.append(person)
    return watchers

def render_group(group : WatchGroup, shows_list:list[Show],people_list:list[Person]) -> str: 
    group_div_template : str = '<div class="group_item"><div class="group_header">%s (%s)</div><br/></br>%s</div>'
    group_show : Show | None = get_show_group(group,shows_list)
    if group_show is None:
        return ''

    group_people : list[Person] = get_people_group(group,people_list)
    people_names : list[str] = list(map(lambda x:x.name,group_people))
    people_str : str = ', '.join(people_names)

    return group_div_template % (group_show.name,group.episode_nr,people_str)

def write_html(html_str : str) -> None:
    index_out : str = os.path.join(html_dir,index_path)
    with open(index_out,'w+') as out_file:
        out_file.write(html_str)
        out_file.close()

if __name__=='__main__':
    index_template : jinja2.Template = get_jinja_template()

    shows_list : list[Show]         = load_shows()
    shows_names : list[str] = list(map(lambda x:'"%s"' % x.name,shows_list))
    people_list : list[Person]      = load_people()
    groups_list : list[WatchGroup]  = load_groups()

    group_strs : list[str] = list(map(lambda x:render_group(x,shows_list,people_list),groups_list))

    info_dict : dict[str,str] = {
            'shows_list':'\n'.join(group_strs),
            'show_names':', '.join(shows_names)}
    html_out : str = index_template.render(info_dict)
    write_html(html_out)
