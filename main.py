from Person import Person
from Show import Show
from ShowManager import ShowManager
import load_csv
from watcher_bot import WatcherBot

if __name__ == '__main__':
    manager : ShowManager = ShowManager()
    manager.load_from_csv()
    bot = WatcherBot(manager)
    bot.run()
