from telegram.ext import Application,ApplicationBuilder,CommandHandler,ContextTypes
from telegram import Update
from common.types import WatcherError,show_err
from common.ShowManager import ShowManager
from common.Person import Person
from common.Show import Show

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
                ('add_watcher',     '$person;$show[;$group]', 'add person watching show to group (default 0)',                     self.add_watcher),
                ('remove_watcher',  '$person;$show[;$group]', 'remove person watching show from group (default first)',            self.remove_watcher),
                ('update_show',     '$show;$nr[;$group]',     'update episode number for show for gropu (default 0)',               self.update_show),
                ('add_show',        '$show',         'add new show',                                 self.add_show),
                ('show_shows',      '$person',       'show shows person is watching',                self.show_shows),
                ('show_people',     '',              'show all people',                              self.show_people),
                ('remove_show',     '$show',         'remove a show',                                self.remove_show),
                ('recommend_show',  '$people', 'recommend a show for given people', self.recommend_show)
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
        shows    : list[Show] | WatcherError = self.show_manager.get_possible(watchers)

        match shows:
            case WatcherError():
                ret_msg : str = show_err(shows)
                await self.send_message(update,context,ret_msg)
                return
        
        if len(shows) == 0:
            ret_msg : str = 'No show to watch with %s' % ', '.join(watchers)

        shows_strs : list[str] = list(map(lambda x: x.name,shows))

        ret_msg : str = 'Shows to watch with %s:\n%s' % (', '.join(watchers),'\n'.join(shows_strs))
        await self.send_message(update,context,ret_msg)



    async def add_watcher(self,update,context) -> None:
        msg_content : str = self.get_message_text(update)
        msg_args : list[str] = msg_content.split(';')
        if len(msg_args) < 2:
            ret_msg : str = 'Malformed command, please try again'
            await self.send_message(update,context,ret_msg)
            return
        
        new_watcher : str = msg_args[0].strip()
        new_show : str = msg_args[1].strip()
        group_id : int | None = None
        if len(msg_args) > 2:
            group_id : int | None = int(msg_args[2].strip())

        add_res : WatcherError | None = self.show_manager.add_watcher_show(new_watcher,new_show,group_id)
        match add_res: 
            case WatcherError():
                ret_msg : str = show_err(add_res)
                await self.send_message(update,context,ret_msg)
                return

        ret_msg : str = 'Added %s to show %s' % (new_watcher,new_show)
        await self.send_message(update,context,ret_msg)
    
    async def remove_watcher(self,update,context) -> None:
        msg_content : str = self.get_message_text(update)
        msg_args : list[str] = msg_content.split(';')
        if len(msg_args) < 2:
            ret_msg : str = 'Malformed command, please try again'
            await self.send_message(update,context,ret_msg)
            return 

        new_watcher : str = msg_args[0].strip()
        new_show : str = msg_args[1].strip()
        group_id : int | None = None
        if len(msg_args) > 2: 
            group_id : int | None = int(msg_args[2].strip())

        rem_res : WatcherError | None = self.show_manager.remove_watcher_show(new_watcher,new_show,group_id)
        match rem_res:
            case WatcherError():
                ret_msg : str = show_err(rem_res)
                await self.send_message(update,context,ret_msg)
                return

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

        update_res : WatcherError | None = self.show_manager.update_show_episode(show_name,ep_nr)
        match update_res:
            case WatcherError():
                ret_msg : str = show_err(update_res)
                await self.send_message(update,context,ret_msg)
                return 

        ret_msg:str = 'Updated show %s' %show_name
        await self.send_message(update,context,ret_msg)


    async def add_show(self,update,context) -> None:
        show_name : str = self.get_message_text(update).strip()
        add_res : WatcherError | None = self.show_manager.add_show(show_name)
        match add_res: 
            case WatcherError():
                ret_msg : str = show_err(add_res)
                await self.send_message(update,context,ret_msg)
                return 

        ret_msg : str = 'Added show %s' % show_name
        await self.send_message(update,context,ret_msg)

    async def show_shows(self,update,context) -> None:
        person_name : str = self.get_message_text(update).strip()
        shows_list : list[Show] | WatcherError = self.show_manager.get_shows_person(person_name)
        match shows_list:
            case WatcherError():
                ret_msg : str = show_err(shows_list)
                await self.send_message(update,context,ret_msg)
                return 
        shows_strs : list[str] = list(map(lambda x: x.name,shows_list))
    
        ret_msg : str = 'Shows for %s: %s' % (person_name,'\n'.join(shows_strs))
        await self.send_message(update,context,ret_msg)

    async def show_people(self,update,context) -> None:
        people_names : list[str] = list(map(lambda x: x.name,self.show_manager.people))
        ret_msg : str = 'All watching people:\n%s' % '\n'.join(people_names)
        await self.send_message(update,context,ret_msg)

    async def remove_show(self,update,context) -> None:
        show_name : str = self.get_message_text(update).strip()
        rem_res : WatcherError | None = self.show_manager.remove_show(show_name)
        match rem_res:
            case WatcherError():
                ret_msg : str = show_err(rem_res)
                await self.send_message(update,context,ret_msg)
                return
        ret_msg: str = 'Removed show %s' % show_name
        await self.send_message(update,context,ret_msg)

    async def recommend_show(self,update,context) -> None:
        people_names : list[str] = self.get_message_text(update).split(',')
        people_names = list(map(lambda x: x.strip(),people_names))
        people : list[Person] = []
        for person_name in people_names:
            maybe_person : Person | WatcherError = self.show_manager.get_person_by_name(person_name)
            match maybe_person:
                case WatcherError():
                    ret_msg : str = show_err(maybe_person)
                    await self.send_message(update,context,ret_msg)
                    return
            people.append(maybe_person)

        people_ids : list[int] = list(map(lambda x:x.person_id,people))
        recommended_show : Show | WatcherError = self.show_manager.recommend_show(people_ids)
        match recommended_show: 
            case WatcherError():
                ret_msg : str = show_err(recommended_show)
                await self.send_message(update,context,ret_msg)
                return
            case Show():
                ret_msg : str = 'Recommended Show for %s:\n%s' % (', '.join(people_names),recommended_show.name)
                await self.send_message(update,context,ret_msg)
                return
