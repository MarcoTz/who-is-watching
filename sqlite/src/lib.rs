use shows::Show;
use std::path::PathBuf;
use turso::{Builder, Connection};

pub mod errors;
pub mod schema;
pub mod show_management;

use errors::{Error, SqlAction, StatementType};
use schema::{ColumnName, Table};

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

    pub async fn get_max_season(&self, show_id: u32) -> Result<u32, Error> {
        show_management::get_max_season(self, show_id).await
    }

    pub async fn get_show_ids(&self) -> Result<Vec<u32>, Error> {
        show_management::get_show_ids(self).await
    }

    pub async fn get_show(&self, show_id: u32) -> Result<Show, Error> {
        show_management::get_show(self, show_id).await
    }

    pub async fn show_exists(&self, show_name: &str) -> Result<bool, Error> {
        show_management::show_exists(self, show_name).await
    }

    pub async fn add_season(
        &self,
        show_id: u32,
        season_nr: u32,
        num_episodes: u32,
    ) -> Result<(), Error> {
        show_management::add_season(self, show_id, season_nr, num_episodes).await
    }

    pub async fn season_exists(&self, show_id: u32, season_nr: u32) -> Result<bool, Error> {
        show_management::season_exists(self, show_id, season_nr).await
    }

    pub async fn get_show_name(&self, show_id: u32) -> Result<String, Error> {
        show_management::get_show_name(self, show_id).await
    }
}
