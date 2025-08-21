use crate::{
    DBDriver,
    errors::{Error, SqlAction, StatementType},
    schema::{ColumnName, Table},
    show_management,
};
use shows::ShowProgress;
use watchers::Watcher;

pub async fn get_watchers(drv: &DBDriver) -> Result<Vec<Watcher>, Error> {
    let ids = get_watcher_ids(drv).await?;
    let mut watchers = vec![];
    for watcher_id in ids {
        let watcher = get_watcher(drv, watcher_id).await?;
        watchers.push(watcher);
    }
    Ok(watchers)
}

pub async fn get_watcher_ids(drv: &DBDriver) -> Result<Vec<u32>, Error> {
    let query = format!("SELECT {} FROM {}", ColumnName::Id, Table::Watchers);
    let mut stmt = drv.conn.prepare(&query).await.map_err(|err| {
        Error::turso(
            err,
            SqlAction::PrepareStatement(StatementType::Select),
            "get_watcher_ids",
        )
    })?;
    let mut rows = stmt.query(()).await.map_err(|err| {
        Error::turso(
            err,
            SqlAction::ExecuteStatement(StatementType::Select),
            "get_watcher_ids",
        )
    })?;

    let mut ids = vec![];

    while let Some(row) = rows
        .next()
        .await
        .map_err(|err| Error::turso(err, SqlAction::GetNextRow, "get_watcher_ids"))?
    {
        let id_val = row.get_value(0).map_err(|err| {
            Error::turso(err, SqlAction::GetValue("id".to_owned()), "get_watcher_ids")
        })?;
        let id = u32::try_from(*id_val.as_integer().ok_or(Error::cast(
            &Table::Watchers,
            &ColumnName::Id,
            "integer",
        ))?)
        .map_err(|_| Error::neg_id(&Table::Watchers, &ColumnName::Id))?;
        ids.push(id);
    }
    Ok(ids)
}

pub async fn get_watcher(drv: &DBDriver, watcher_id: u32) -> Result<Watcher, Error> {
    let query = format!(
        "SELECT {} FROM {} WHERE {}=?1",
        ColumnName::Name,
        Table::Watchers,
        ColumnName::Id
    );
    let mut stmt = drv.conn.prepare(&query).await.map_err(|err| {
        Error::turso(
            err,
            SqlAction::PrepareStatement(StatementType::Select),
            "get_watcher",
        )
    })?;
    let mut rows = stmt.query([watcher_id]).await.map_err(|err| {
        Error::turso(
            err,
            SqlAction::ExecuteStatement(StatementType::Select),
            "get_watcher",
        )
    })?;
    let row = rows
        .next()
        .await
        .map_err(|err| Error::turso(err, SqlAction::GetNextRow, "get_watcher"))?
        .ok_or(Error::no_rows(&query))?;
    let name_val = row
        .get_value(0)
        .map_err(|err| Error::turso(err, SqlAction::GetValue("Name".to_owned()), "get_watcher"))?;
    let name =
        name_val
            .as_text()
            .ok_or(Error::cast(&Table::Watchers, &ColumnName::Name, "text"))?;
    let progress = get_watcher_shows(drv, watcher_id).await?;
    Ok(Watcher::new(watcher_id, name.as_str(), progress))
}

pub async fn get_watcher_shows(
    drv: &DBDriver,
    watcher_id: u32,
) -> Result<Vec<ShowProgress>, Error> {
    let query = format!(
        "SELECT {},{},{},{} FROM {},{} WHERE {}=?1 AND {} = {}",
        ColumnName::ShowId.qualified(&Table::ShowsWatchers),
        ColumnName::Name.qualified(&Table::Shows),
        ColumnName::SeasonNum.qualified(&Table::ShowsWatchers),
        ColumnName::NumEpisodes.qualified(&Table::ShowsWatchers),
        Table::ShowsWatchers,
        Table::Shows,
        ColumnName::WatcherId.qualified(&Table::ShowsWatchers),
        ColumnName::Id.qualified(&Table::Shows),
        ColumnName::ShowId.qualified(&Table::ShowsWatchers),
    );
    let mut stmt = drv.conn.prepare(&query).await.map_err(|err| {
        Error::turso(
            err,
            SqlAction::PrepareStatement(StatementType::Select),
            "get_watcher_shows",
        )
    })?;
    let mut rows = stmt.query([watcher_id]).await.map_err(|err| {
        Error::turso(
            err,
            SqlAction::ExecuteStatement(StatementType::Select),
            "get_watcher_shows",
        )
    })?;
    let mut progress = vec![];
    while let Some(row) = rows
        .next()
        .await
        .map_err(|err| Error::turso(err, SqlAction::GetNextRow, "get_watcher_shows"))?
    {
        let show_id = u32::try_from(
            *row.get_value(0)
                .map_err(|err| {
                    Error::turso(
                        err,
                        SqlAction::GetValue("show_id".to_owned()),
                        "get_watcher_shows",
                    )
                })?
                .as_integer()
                .ok_or(Error::cast(
                    &Table::ShowsWatchers,
                    &ColumnName::ShowId,
                    "integer",
                ))?,
        )
        .map_err(|_| Error::neg_id(&Table::ShowsWatchers, &ColumnName::ShowId))?;
        let show_name = row
            .get_value(1)
            .map_err(|err| {
                Error::turso(
                    err,
                    SqlAction::GetValue("show_name".to_owned()),
                    "get_watcher_shows",
                )
            })?
            .as_text()
            .ok_or(Error::cast(&Table::Shows, &ColumnName::Name, "text"))?
            .clone();
        let season_nr = u32::try_from(
            *row.get_value(2)
                .map_err(|err| {
                    Error::turso(
                        err,
                        SqlAction::GetValue("season_num".to_owned()),
                        "get_watcher_shows",
                    )
                })?
                .as_integer()
                .ok_or(Error::cast(
                    &Table::ShowsWatchers,
                    &ColumnName::SeasonNum,
                    "integer",
                ))?,
        )
        .map_err(|_| Error::neg_id(&Table::ShowsWatchers, &ColumnName::SeasonNum))?;
        let episode_nr = u32::try_from(
            *row.get_value(2)
                .map_err(|err| {
                    Error::turso(
                        err,
                        SqlAction::GetValue("num_episodes".to_owned()),
                        "get_watcher_shows",
                    )
                })?
                .as_integer()
                .ok_or(Error::cast(
                    &Table::ShowsWatchers,
                    &ColumnName::NumEpisodes,
                    "integer",
                ))?,
        )
        .map_err(|_| Error::neg_id(&Table::ShowsWatchers, &ColumnName::NumEpisodes))?;

        progress.push(ShowProgress::new(
            show_id, &show_name, season_nr, episode_nr,
        ))
    }
    Ok(progress)
}

pub async fn watcher_exists(drv: &DBDriver, watcher_name: &str) -> Result<bool, Error> {
    let query = format!(
        "SELECT COUNT({}) FROM {} WHERE {}=?1",
        ColumnName::Id,
        Table::Watchers,
        ColumnName::Name
    );
    let mut stmt = drv.conn.prepare(&query).await.map_err(|err| {
        Error::turso(
            err,
            SqlAction::PrepareStatement(StatementType::Select),
            "watcher_exists",
        )
    })?;
    let mut rows = stmt.query([watcher_name]).await.map_err(|err| {
        Error::turso(
            err,
            SqlAction::ExecuteStatement(StatementType::Select),
            "watcher_exists",
        )
    })?;
    let row = rows
        .next()
        .await
        .map_err(|err| Error::turso(err, SqlAction::GetNextRow, "watcher_exists"))?
        .ok_or(Error::no_rows(&query))?;
    let cnt = *row
        .get_value(0)
        .map_err(|err| {
            Error::turso(
                err,
                SqlAction::GetValue("Count".to_owned()),
                "watcher_exists",
            )
        })?
        .as_integer()
        .ok_or(Error::cast(
            &Table::Watchers,
            &ColumnName::WatcherId,
            "integer",
        ))?;
    Ok(cnt != 0)
}

pub(crate) async fn get_watcher_id(drv: &DBDriver, watcher_name: &str) -> Result<u32, Error> {
    let query = format!(
        "SELECT {} FROM {} WHERE {}=?1",
        ColumnName::Id,
        Table::Watchers,
        ColumnName::Name
    );
    let mut stmt = drv.conn.prepare(&query).await.map_err(|err| {
        Error::turso(
            err,
            SqlAction::PrepareStatement(StatementType::Select),
            "get_watcher_id",
        )
    })?;

    let mut rows = stmt.query([watcher_name]).await.map_err(|err| {
        Error::turso(
            err,
            SqlAction::ExecuteStatement(StatementType::Select),
            "get_watcher_id",
        )
    })?;

    let row = rows
        .next()
        .await
        .map_err(|err| Error::turso(err, SqlAction::GetNextRow, "get_watcher_id"))?
        .ok_or(Error::no_rows(&query))?;
    let id = u32::try_from(
        *row.get_value(0)
            .map_err(|err| Error::turso(err, SqlAction::GetNextRow, "get_watcher_id"))?
            .as_integer()
            .ok_or(Error::cast(&Table::Watchers, &ColumnName::Id, "integer"))?,
    )
    .map_err(|_| Error::neg_id(&Table::Watchers, &ColumnName::Id))?;
    Ok(id)
}

pub(crate) async fn create_watcher(drv: &DBDriver, watcher: &Watcher) -> Result<(), Error> {
    if drv.watcher_exists(&watcher.name).await? {
        return Err(Error::watcher_exists(&watcher.name));
    }
    let query = format!(
        "INSERT INTO {} ({},{}) VALUES (?1,?2)",
        Table::Watchers,
        ColumnName::Id,
        ColumnName::Name
    );
    drv.conn
        .execute(&query, [&watcher.id.to_string(), watcher.name.as_str()])
        .await
        .map_err(|err| {
            Error::turso(
                err,
                SqlAction::ExecuteStatement(StatementType::Insert),
                "create_watcher",
            )
        })?;
    for progress in watcher.watching.iter() {
        show_management::update_progress(drv, watcher.id, progress).await?;
    }
    Ok(())
}
