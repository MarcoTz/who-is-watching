use crate::{
    DBDriver,
    errors::{Error, SqlAction, StatementType},
    schema::{ColumnName, Table},
};
use shows::{Show, ShowProgress};

pub(crate) async fn show_exists(drv: &DBDriver, name: &str) -> Result<bool, Error> {
    let query = format!(
        "SELECT COUNT({}) FROM {} WHERE {}=?1",
        ColumnName::Id,
        Table::Shows,
        ColumnName::Name
    );
    let mut stmt = drv.conn.prepare(&query).await.map_err(|err| {
        Error::turso(
            err,
            SqlAction::PrepareStatement(StatementType::Select),
            "show_exists",
        )
    })?;
    let mut rows = stmt.query([name]).await.map_err(|err| {
        Error::turso(
            err,
            SqlAction::ExecuteStatement(StatementType::Select),
            "show_exists",
        )
    })?;
    let row = rows
        .next()
        .await
        .map_err(|err| Error::turso(err, SqlAction::GetNextRow, "show_exists"))?
        .ok_or(Error::no_rows(&query))?;
    let cnt = *row
        .get_value(0)
        .map_err(|err| Error::turso(err, SqlAction::GetValue("count".to_owned()), "show_exists"))?
        .as_integer()
        .ok_or(Error::cast(&Table::Shows, &ColumnName::Id, "integer"))?;
    Ok(cnt != 0)
}

async fn season_exists(drv: &DBDriver, show_id: u32, season_nr: u32) -> Result<bool, Error> {
    let query = format!(
        "SELECT COUNT({}) FROM {} WHERE {}=?1 AND {}=?2",
        ColumnName::ShowId,
        Table::ShowsSeasons,
        ColumnName::ShowId,
        ColumnName::SeasonNum
    );
    let mut stmt = drv.conn.prepare(&query).await.map_err(|err| {
        Error::turso(
            err,
            SqlAction::PrepareStatement(StatementType::Select),
            "season_exists",
        )
    })?;
    let mut rows = stmt.query([show_id, season_nr]).await.map_err(|err| {
        Error::turso(
            err,
            SqlAction::ExecuteStatement(StatementType::Select),
            "season_exists",
        )
    })?;
    let row = rows
        .next()
        .await
        .map_err(|err| Error::turso(err, SqlAction::GetNextRow, "season_exists"))?
        .ok_or(Error::no_rows(&query))?;
    let cnt = *row
        .get_value(0)
        .map_err(|err| {
            Error::turso(
                err,
                SqlAction::GetValue("Count".to_owned()),
                "season_exists",
            )
        })?
        .as_integer()
        .ok_or(Error::cast(
            &Table::ShowsSeasons,
            &ColumnName::Id,
            "integer",
        ))?;
    Ok(cnt != 0)
}

async fn get_show_name(drv: &DBDriver, show_id: u32) -> Result<String, Error> {
    let mut stmt = drv
        .conn
        .prepare(&format!(
            "SELECT {} FROM {} WHERE id=?1",
            ColumnName::Name,
            Table::Shows.name()
        ))
        .await
        .map_err(|err| {
            Error::turso(
                err,
                SqlAction::PrepareStatement(StatementType::Select),
                "get_show_name",
            )
        })?;
    let mut rows = stmt.query([show_id]).await.map_err(|err| {
        Error::turso(
            err,
            SqlAction::ExecuteStatement(StatementType::Select),
            "get_show_name",
        )
    })?;
    let row = rows
        .next()
        .await
        .map_err(|err| Error::turso(err, SqlAction::GetNextRow, "get_show_name"))?;
    let row = row.ok_or(Error::show_not_found(show_id))?;
    let name_val = row.get_value(0).map_err(|err| {
        Error::turso(err, SqlAction::GetValue("name".to_owned()), "get_show_name")
    })?;
    let name = name_val
        .as_text()
        .ok_or(Error::cast(&Table::Shows, &ColumnName::Name, "text"))?;
    Ok(name.clone())
}

async fn get_max_season(drv: &DBDriver, show_id: u32) -> Result<u32, Error> {
    let show_name = get_show_name(drv, show_id).await?;
    let mut stmt = drv
        .conn
        .prepare(&format!(
            "SELECT MAX({}) FROM {} WHERE {}=?1;",
            ColumnName::SeasonNum,
            Table::ShowsSeasons.name(),
            ColumnName::ShowId
        ))
        .await
        .map_err(|err| {
            Error::turso(
                err,
                SqlAction::PrepareStatement(StatementType::Select),
                "get_max_season",
            )
        })?;
    let mut rows = stmt.query([show_id]).await.map_err(|err| {
        Error::turso(
            err,
            SqlAction::ExecuteStatement(StatementType::Select),
            "get_max_season",
        )
    })?;
    let row = rows
        .next()
        .await
        .map_err(|err| Error::turso(err, SqlAction::GetNextRow, "get_max_season"))?;
    let row = row.ok_or(Error::no_seasons(&show_name))?;
    let season_val = row.get_value(0).map_err(|err| {
        Error::turso(
            err,
            SqlAction::GetValue("season_num".to_owned()),
            "get_max_season",
        )
    })?;
    let season = u32::try_from(
        *season_val
            .as_integer()
            .ok_or(Error::no_seasons(&show_name))?,
    )
    .map_err(|_| Error::neg_id(&Table::ShowsSeasons, &ColumnName::SeasonNum))?;
    Ok(season)
}

async fn get_show(drv: &DBDriver, show_id: u32) -> Result<Show, Error> {
    let max_season = get_max_season(drv, show_id).await?;
    let mut stmt = drv
        .conn
        .prepare(&format!(
            "
                SELECT {}.{},se.{} 
                FROM {},{} as se 
                WHERE se.{}={}.{} AND se.{} = ?1;",
            Table::Shows,
            ColumnName::Name,
            ColumnName::NumEpisodes,
            Table::Shows.name(),
            Table::ShowsSeasons.name(),
            ColumnName::ShowId,
            Table::Shows,
            ColumnName::Id,
            ColumnName::SeasonNum
        ))
        .await
        .map_err(|err| {
            Error::turso(
                err,
                SqlAction::PrepareStatement(StatementType::Select),
                "get_show",
            )
        })?;
    let mut rows = stmt.query([max_season]).await.map_err(|err| {
        Error::turso(
            err,
            SqlAction::ExecuteStatement(StatementType::Select),
            "get_show",
        )
    })?;
    let row = rows
        .next()
        .await
        .map_err(|err| Error::turso(err, SqlAction::GetNextRow, "get_show"))?
        .ok_or(Error::show_not_found(show_id))?;
    let name_val = row
        .get_value(0)
        .map_err(|err| Error::turso(err, SqlAction::GetValue("Name".to_owned()), "get_show"))?;
    let name = name_val
        .as_text()
        .ok_or(Error::cast(&Table::Shows, &ColumnName::Name, "text"))?;
    let episodes = u32::try_from(
        *row.get_value(1)
            .map_err(|err| {
                Error::turso(
                    err,
                    SqlAction::GetValue("num_episodes".to_owned()),
                    "get_show",
                )
            })?
            .as_integer()
            .ok_or(Error::cast(
                &Table::Shows,
                &ColumnName::NumEpisodes,
                "integer",
            ))?,
    )
    .map_err(|_| Error::neg_id(&Table::ShowsSeasons, &ColumnName::NumEpisodes))?;
    Ok(Show::new(show_id, name.as_str(), max_season, episodes))
}

async fn get_show_ids(drv: &DBDriver) -> Result<Vec<u32>, Error> {
    let mut stmt = drv
        .conn
        .prepare(&format!(
            "SELECT {} FROM {};",
            ColumnName::Id,
            Table::Shows.name()
        ))
        .await
        .map_err(|err| {
            Error::turso(
                err,
                SqlAction::PrepareStatement(StatementType::Select),
                "get_show_ids",
            )
        })?;
    let mut rows = stmt.query(()).await.map_err(|err| {
        Error::turso(
            err,
            SqlAction::ExecuteStatement(StatementType::Select),
            "get_show_ids",
        )
    })?;
    let mut ids = vec![];
    while let Some(row) = rows
        .next()
        .await
        .map_err(|err| Error::turso(err, SqlAction::GetNextRow, "get_show_ids"))?
    {
        let id = u32::try_from(
            *row.get_value(0)
                .map_err(|err| {
                    Error::turso(err, SqlAction::GetValue("id".to_owned()), "get_show_ids")
                })?
                .as_integer()
                .ok_or(Error::cast(&Table::Shows, &ColumnName::Id, "integer"))?,
        )
        .map_err(|_| Error::neg_id(&Table::Shows, &ColumnName::Id))?;
        ids.push(id);
    }
    Ok(ids)
}

pub(crate) async fn get_shows(drv: &DBDriver) -> Result<Vec<Show>, Error> {
    let mut shows = vec![];
    let ids = get_show_ids(drv).await?;
    for show_id in ids {
        shows.push(get_show(drv, show_id).await?);
    }
    Ok(shows)
}

pub(crate) async fn create_show(drv: &DBDriver, show: &Show) -> Result<(), Error> {
    if show_exists(drv, &show.name).await? {
        return Err(Error::show_exists(&show.name));
    }
    let shows_query = format!(
        "INSERT INTO {} ({}, {}) VALUES (?1,?2);",
        Table::Shows,
        ColumnName::Id,
        ColumnName::Name
    );
    drv.conn
        .execute(
            &shows_query,
            [show.id.to_string().as_str(), show.name.as_str()],
        )
        .await
        .map_err(|err| {
            Error::turso(
                err,
                SqlAction::ExecuteStatement(StatementType::Insert),
                "create_show",
            )
        })?;

    add_season(
        drv,
        show.id,
        show.latest_episode.season_nr,
        show.latest_episode.episode_nr,
    )
    .await?;
    Ok(())
}

async fn add_season(
    drv: &DBDriver,
    show_id: u32,
    season_nr: u32,
    num_episodes: u32,
) -> Result<(), Error> {
    if season_exists(drv, show_id, season_nr).await? {
        let show_name = get_show_name(drv, show_id).await?;
        return Err(Error::season_exists(&show_name, season_nr));
    }
    let seasons_query = format!(
        "INSERT INTO {} ({}, {}, {}) VALUES (?1,?2,?3);",
        Table::ShowsSeasons,
        ColumnName::ShowId,
        ColumnName::SeasonNum,
        ColumnName::NumEpisodes
    );
    drv.conn
        .execute(&seasons_query, [show_id, season_nr, num_episodes])
        .await
        .map_err(|err| {
            Error::turso(
                err,
                SqlAction::ExecuteStatement(StatementType::Insert),
                "add_season",
            )
        })?;
    Ok(())
}

async fn progress_exists(drv: &DBDriver, watcher_id: u32, show_id: u32) -> Result<bool, Error> {
    let query = format!(
        "SELECT COUNT(*) FROM {} WHERE {}=?1 AND {}=?2",
        Table::ShowsWatchers,
        ColumnName::WatcherId,
        ColumnName::ShowId
    );
    let mut stmt = drv.conn.prepare(&query).await.map_err(|err| {
        Error::turso(
            err,
            SqlAction::PrepareStatement(StatementType::Select),
            "progress_exists",
        )
    })?;
    let mut rows = stmt.query([watcher_id, show_id]).await.map_err(|err| {
        Error::turso(
            err,
            SqlAction::ExecuteStatement(StatementType::Select),
            "progress_exists",
        )
    })?;
    let row = rows
        .next()
        .await
        .map_err(|err| Error::turso(err, SqlAction::GetNextRow, "progress_exists"))?
        .ok_or(Error::no_rows(&query))?;
    let cnt = *row
        .get_value(0)
        .map_err(|err| {
            Error::turso(
                err,
                SqlAction::GetValue("Count".to_owned()),
                "progress_exists",
            )
        })?
        .as_integer()
        .ok_or(Error::cast(
            &Table::ShowsWatchers,
            &ColumnName::WatcherId,
            "integer",
        ))?;
    Ok(cnt != 0)
}

pub(crate) async fn update_progress(
    drv: &DBDriver,
    watcher_id: u32,
    progress: &ShowProgress,
) -> Result<(), Error> {
    if !progress_exists(drv, watcher_id, progress.show_id).await? {
        return create_progress(drv, watcher_id, progress).await;
    }

    let query = format!(
        "UPDATE {} SET {}=?1, {}=?2 WHERE {}=?3 AND {}=?4",
        Table::ShowsWatchers,
        ColumnName::SeasonNum,
        ColumnName::NumEpisodes,
        ColumnName::WatcherId,
        ColumnName::ShowId
    );
    drv.conn
        .execute(
            &query,
            [
                progress.last_watched.season_nr,
                progress.last_watched.episode_nr,
                watcher_id,
                progress.show_id,
            ],
        )
        .await
        .map_err(|err| {
            Error::turso(
                err,
                SqlAction::ExecuteStatement(StatementType::Insert),
                "update_progress",
            )
        })?;
    Ok(())
}

async fn create_progress(
    drv: &DBDriver,
    watcher_id: u32,
    progress: &ShowProgress,
) -> Result<(), Error> {
    let query = format!(
        "INSERT INTO {} ({},{},{},{}) VALUES (?1,?2,?3,?4)",
        Table::ShowsWatchers,
        ColumnName::WatcherId,
        ColumnName::ShowId,
        ColumnName::SeasonNum,
        ColumnName::NumEpisodes
    );
    drv.conn
        .execute(
            &query,
            [
                watcher_id,
                progress.show_id,
                progress.last_watched.season_nr,
                progress.last_watched.episode_nr,
            ],
        )
        .await
        .map_err(|err| {
            Error::turso(
                err,
                SqlAction::ExecuteStatement(StatementType::Insert),
                "create_progress",
            )
        })?;
    Ok(())
}
