from telegram.ext import Application,ApplicationBuilder,CommandHandler,ContextTypes
from telegram import Update
from ShowManager import ShowManager
from Person import Person 
from Show import Show

def load_api_key() -> str:
    config_file_path = 'bot.cfg' 
    config_file = open(config_file_path)
    return config_file.read().strip()

class WatcherBot: 
    api_key         : str
    application     : Application
    show_manager    : ShowManager

    def __init__(self, manager:ShowManager):
        self.show_manager = manager
        self.api_key : str = load_api_key()
        self.application = ApplicationBuilder().token(self.api_key).build()

        self.cmd_actions = [
                ('help',            '',              'show help message',                            self.get_help),
                ('possible_shows',  '$people',       'get list of possible shows to watch',          self.get_possible),
                ('add_watcher',     '$person;$show', 'add person watching show',                     self.add_watcher),
                ('remove_watcher',  '$person;$show', 'remove person watching show',                  self.remove_watcher),
                ('update_show',     '$show;$nr',     'update episode number for show',               self.update_show),
                ('add_show',        '$show',         'add new show',                                 self.add_show),
                ('show_shows',      '$person',       'show shows person is watching',                self.show_shows),
                ('show_people',     '',              'show people',                                  self.show_people),
                ('remove_show',     '$show',         'reove a show',                                 self.remove_show)
            ]

        for (cmd,_,_,action) in self.cmd_actions:
            new_handler : CommandHandler = CommandHandler(cmd,action)
            self.application.add_handler(new_handler)

    def run(self):
        self.application.run_polling()

    async def send_message(self,update:Update,context:ContextTypes.DEFAULT_TYPE, msg:str) -> None: 
        if update.effective_chat is None or update.effective_chat.id is None:
            return 
        chat_id : int = update.effective_chat.id
        await context.bot.send_message(chat_id=chat_id,text=msg)

    def get_message_text(self,update:Update) -> str:
        if update.effective_message is None or update.effective_message.text is None:
            return ''
        else:
            msg : str = update.effective_message.text
            msg_cmd : str = msg.split(' ')[0]
            return msg.replace(msg_cmd,'')

    async def get_help(self,update,context) -> None:

        help_template : str = '/%s %s - %s' 
        help_strs : list[str] = []
        for (cmd,cmd_args,cmd_info,_) in self.cmd_actions:
            help_strs.append(help_template % (cmd,cmd_args,cmd_info))


        help_message : str = 'Possible commands:\n ' + '\n '.join(help_strs) 

        await self.send_message(update,context,help_message)
    
    async def get_possible(self,update,context) -> None:
        watcher_str : str = self.get_message_text(update)
        if watcher_str == '':
            ret_msg : str = 'Please provide at least one person to see shows'
            await self.send_message(update,context,ret_msg)
            return

        watchers : list[str] = list(map(lambda x: x.strip(),watcher_str.split(', ')))
        
        people_list : list[Person] = [] 
        not_found : list[str] = []
        for watcher_name in watchers:
            watcher : Person | None = self.show_manager.get_person(watcher_name)
            if watcher is None:
                not_found.append(watcher_name)
                continue 
            people_list.append(watcher)

        if len(people_list) == 0:
            ret_msg : str = 'Could not find people %s' % ', '.join(watchers)
        shows : list[Show] = self.show_manager.get_possible_shows(people_list)

        watcher_names : list[str] = list(map(lambda x: x.name,people_list))
        if len(shows) == 0:
            ret_msg : str = 'No show to watch with %s ' % (', '.join(watcher_names))
            await self.send_message(update,context,ret_msg)
            return

        show_strs : list[str] = list(map(lambda x: '%s (%s)' %(x.name,str(x.current_ep)),shows))
        ret_template : str = 'Possible shows to watch with %s:\n%s'
        ret_msg : str =  ret_template % (', '.join(watcher_names),'\n'.join(show_strs))
        if len(not_found) > 0:
            ret_msg += '\n\nCould not find %s' % ','.join(not_found)

        await self.send_message(update,context,ret_msg)

    async def add_watcher(self,update,context) -> None:
        msg_content : str = self.get_message_text(update)
        msg_args : list[str] = msg_content.split(';')
        if len(msg_args) != 2:
            ret_msg : str = 'Malformed command, please try again'
            await self.send_message(update,context,ret_msg)
            return
        
        new_watcher : str = msg_args[0].strip()
        new_show : str = msg_args[1].strip()

        self.show_manager.add_watcher(new_watcher,new_show)
        ret_msg : str = 'Added %s to show %s' % (new_watcher,new_show)
        await self.send_message(update,context,ret_msg)
    
    async def remove_watcher(self,update,context) -> None:
        msg_content : str = self.get_message_text(update)
        msg_args : list[str] = msg_content.split(';')
        if len(msg_args) != 2:
            ret_msg : str = 'Malformed command, please try again'
            await self.send_message(update,context,ret_msg)
            return 

        new_watcher : str = msg_args[0].strip()
        new_show : str = msg_args[1].strip()
        self.show_manager.remove_watcher(new_watcher,new_show)
        ret_msg : str = 'Removed %s from %s' % (new_watcher,new_show)
        await self.send_message(update,context,ret_msg)

        
    async def update_show(self,update,context) -> None:
        msg_content : str = self.get_message_text(update)
        msg_args = msg_content.split(';')
        if len(msg_args) != 2: 
            ret_msg : str = 'Malformed command, please try again'
            await self.send_message(update,context,ret_msg)
            return
        
        show_name : str = msg_args[0].strip()
        ep_nr : int = -1 
        try:
            ep_nr = int(msg_args[1].strip())
        except:
            ret_msg:str='Cannot parse episode number,please try again'
            await self.send_message(update,context,ret_msg)
            return

        show : Show | None = self.show_manager.get_show(show_name)
        if show is None:
            ret_msg : str = 'Cannot find show %s ' % show_name
            await self.send_message(update,context,ret_msg)
            return
        
        show.update_ep(ep_nr)
        ret_msg:str = 'Updated show %s' %show_name
        await self.send_message(update,context,ret_msg)


    async def add_show(self,update,context) -> None:
        show_name : str = self.get_message_text(update).strip()
        self.show_manager.get_show_add(show_name)
        ret_msg : str = 'Added show %s' % show_name
        await self.send_message(update,context,ret_msg)

    async def show_shows(self,update,context) -> None:
        person_name : str = self.get_message_text(update).strip()
        person : Person | None = self.show_manager.get_person(person_name)
    
        if person is None:
            show_strs : list[str] =list(map(lambda x: '%s (%s)' % (x.name,str(x.current_ep)),self.show_manager.shows))
            ret_msg : str = 'All shows:\n%s' % ('\n'.join(show_strs))
            await self.send_message(update,context,ret_msg)
            return

        person_shows : list[str] = list(map(lambda x: x.name,person.watching))
        ret_msg : str = 'Shows for %s: %s' % (person_name,'\n'.join(person_shows))
        await self.send_message(update,context,ret_msg)

    async def show_people(self,update,context) -> None:
        people_names : list[str] = list(map(lambda x: x.name,self.show_manager.people))
        ret_msg : str = 'All watching people:\n%s' % '\n'.join(people_names)
        await self.send_message(update,context,ret_msg)

    async def remove_show(self,update,context) -> None:
        show_name : str = self.get_message_text(update).strip()
        self.show_manager.remove_show(show_name)
        ret_msg: str = 'Removed show %s' % show_name
        await self.send_message(update,context,ret_msg)
