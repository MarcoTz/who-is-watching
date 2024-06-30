from common.ShowManager import ShowManager
from watcher_bot import WatcherBot

if __name__ == '__main__':
    manager : ShowManager = ShowManager()
    bot = WatcherBot(manager)
    bot.run()
