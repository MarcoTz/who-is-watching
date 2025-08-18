use sqlite::DBDriver;
use turso::Builder;
use watchers::Watcher;

#[tokio::main]
async fn main() -> Result<(), sqlite::errors::Error> {
    let old = Builder::new_local("old/watchers.db").build().await.unwrap();
    let old_conn = old.connect().unwrap();
    let mut stmt = old_conn
        .prepare("SELECT name FROM watchers;")
        .await
        .unwrap();
    let mut rows = stmt.query(()).await.unwrap();
    let mut old_watchers = vec![];
    while let Some(row) = rows.next().await.unwrap() {
        let name_val = row.get_value(0).unwrap();
        let name = name_val.as_text().unwrap();
        println!("watcher {name}");
        old_watchers.push(name.clone());
    }

    let drv = DBDriver::connect().await?;
    let watchers = drv.get_watchers().await?;
    let watcher_ids = watchers
        .iter()
        .map(|watcher| watcher.id)
        .collect::<Vec<u32>>();

    for watcher_name in old_watchers {
        if drv.watcher_exists(&watcher_name).await? {
            continue;
        }

        let new_watcher = Watcher::create(&watcher_name, &watcher_ids);
        drv.create_watcher(&new_watcher).await.unwrap();
    }
    Ok(())
}
