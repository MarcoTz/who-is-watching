use shows::{Show, ShowProgress};
use std::path::PathBuf;
use turso::{Builder, Connection};
use watchers::Watcher;

pub mod errors;
pub mod schema;
pub mod show_management;
pub mod watcher_management;

use errors::{Error, SqlAction};
use schema::Table;

const DB_FILE: &str = "watchers_db.sqlite";

pub struct DBDriver {
    conn: Connection,
}

impl DBDriver {
    pub async fn connect() -> Result<Self, Error> {
        let db_path = PathBuf::from(DB_FILE);
        if db_path.exists() {
            let db = Builder::new_local(DB_FILE)
                .build()
                .await
                .map_err(|err| Error::turso(err, SqlAction::BuildDatabase, "connect"))?;
            let conn = db
                .connect()
                .map_err(|err| Error::turso(err, SqlAction::ConnectToDatabase, "connect"))?;
            Ok(DBDriver { conn })
        } else {
            println!("Warning: Database could not be found, creating new one");
            let db = Builder::new_local(DB_FILE)
                .build()
                .await
                .map_err(|err| Error::turso(err, SqlAction::BuildDatabase, "connect"))?;
            let conn = db
                .connect()
                .map_err(|err| Error::turso(err, SqlAction::ConnectToDatabase, "connect"))?;

            for table in Table::all() {
                println!("creating table {table}");
                println!("{}", table.create_sql());
                conn.execute(&table.create_sql(), ()).await.map_err(|err| {
                    Error::turso(
                        err,
                        SqlAction::CreateTable(table.name().to_owned()),
                        "connect",
                    )
                })?;
            }
            Ok(DBDriver { conn })
        }
    }

    pub async fn show_exists(&self, show_name: &str) -> Result<bool, Error> {
        show_management::show_exists(self, show_name).await
    }

    pub async fn get_show_id(&self, show_name: &str) -> Result<u32, Error> {
        show_management::get_show_id(self, show_name).await
    }

    pub async fn get_show_name(&self, show_id: u32) -> Result<String, Error> {
        show_management::get_show_name(self, show_id).await
    }

    pub async fn watcher_exists(&self, watcher_name: &str) -> Result<bool, Error> {
        watcher_management::watcher_exists(self, watcher_name).await
    }

    pub async fn get_watcher_id(&self, watcher_name: &str) -> Result<u32, Error> {
        watcher_management::get_watcher_id(self, watcher_name).await
    }

    pub async fn get_watchers(&self) -> Result<Vec<Watcher>, Error> {
        watcher_management::get_watchers(self).await
    }

    pub async fn get_shows(&self) -> Result<Vec<Show>, Error> {
        show_management::get_shows(self).await
    }

    pub async fn update_progress(
        &self,
        watcher_id: u32,
        progress: &ShowProgress,
    ) -> Result<(), Error> {
        show_management::update_progress(self, watcher_id, progress).await
    }

    pub async fn create_show(&self, show: &Show) -> Result<(), Error> {
        show_management::create_show(self, show).await
    }

    pub async fn create_watcher(&self, watcher: &Watcher) -> Result<(), Error> {
        watcher_management::create_watcher(self, watcher).await
    }
}
