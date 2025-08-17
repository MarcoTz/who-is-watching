use sqlite::connect_to_db;

#[tokio::main]
async fn main() -> Result<(), sqlite::errors::Error> {
    let _ = connect_to_db().await?;
    Ok(())
}
