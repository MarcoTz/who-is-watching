use std::path::PathBuf;
use turso::{Builder, Connection};

pub mod errors;
pub mod schema;

use errors::Error;
use schema::Tables;

const DB_FILE: &str = "watchers_db.sqlite";

pub async fn connect_to_db() -> Result<Connection, Error> {
    let db_path = PathBuf::from(DB_FILE);
    if db_path.exists() {
        let db = Builder::new_local(DB_FILE).build().await?;
        let conn = db.connect()?;
        Ok(conn)
    } else {
        println!("Warning: Database could not be found, creating new one");
        let db = Builder::new_local(DB_FILE).build().await?;
        let conn = db.connect()?;

        for table in Tables::all() {
            conn.execute(&table.create_sql(), ()).await?;
        }
        Ok(conn)
    }
}
