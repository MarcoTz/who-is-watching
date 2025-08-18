use shows::ShowProgress;
use sqlite::DBDriver;
use turso::Builder;
use watchers::Watcher;

#[tokio::main]
async fn main() -> Result<(), sqlite::errors::Error> {
    let old = Builder::new_local("old/watchers.db").build().await.unwrap();
    let old_conn = old.connect().unwrap();
    let mut stmt = old_conn
        .prepare("select shows.name, watchers.name, watchgroups.current_ep from shows,watchers,watchgroups,watchers_groups where watchgroups.show_id=shows.rowid and watchers_groups.watcher_id=watchers.rowid and watchers_groups.group_id=watchgroups.rowid;")
        .await
        .unwrap();
    let mut rows = stmt.query(()).await.unwrap();
    let mut old_watching = vec![];
    while let Some(row) = rows.next().await.unwrap() {
        let show_name_val = row.get_value(0).unwrap();
        let show_name = show_name_val.as_text().unwrap().clone();
        let watcher_name_val = row.get_value(1).unwrap();
        let watcher_name = watcher_name_val.as_text().unwrap().clone();
        let ep_value = row.get_value(2).unwrap();
        let ep = u32::try_from(*ep_value.as_integer().unwrap()).unwrap();
        old_watching.push((show_name, watcher_name, ep));
    }

    let drv = DBDriver::connect().await?;
    for (show_name, watcher_name, ep_nr) in old_watching {
        println!("Updating {watcher_name} {show_name} {ep_nr}");
        let show_id = drv.get_show_id(&show_name.trim()).await.unwrap();
        let watcher_id = drv.get_watcher_id(&watcher_name.trim()).await.unwrap();
        let progress = ShowProgress::new(show_id, 1, ep_nr);
        drv.update_progress(watcher_id, &progress).await.unwrap()
    }
    Ok(())
}
