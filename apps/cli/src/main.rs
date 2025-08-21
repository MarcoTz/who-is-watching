use sqlite::DBDriver;

#[tokio::main]
async fn main() -> Result<(), sqlite::errors::Error> {
    let drv = DBDriver::connect().await?;
    let watchers = drv.get_watchers().await?;
    println!("All watchers");
    for watcher in watchers {
        println!("{} ({})", watcher.name, watcher.id);
        for progress in watcher.watching {
            let show = drv.get_show_name(progress.show_id).await?;
            println!("\t{} {:?}", show, progress.last_watched);
        }
        println!("");
    }

    Ok(())
}
