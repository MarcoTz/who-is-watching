use crate::errors::Error;
use askama::Template;
use sqlite::DBDriver;
use watchers::Watcher;

#[derive(Template)]
#[template(path = "watchers.html")]
pub struct WatchersTemplate<'a> {
    watchers: &'a [Watcher],
}

impl<'a> WatchersTemplate<'a> {
    pub async fn render(drv: &DBDriver) -> Result<String, Error> {
        let watchers = drv.get_watchers().await?;
        Ok(WatchersTemplate {
            watchers: &watchers,
        }
        .render()?)
    }
}
