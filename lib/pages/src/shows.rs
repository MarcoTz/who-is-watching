use crate::errors::Error;
use askama::Template;
use shows::Show;
use sqlite::DBDriver;

#[derive(Template)]
#[template(path = "shows.html")]
pub struct ShowsTemplate<'a> {
    shows: &'a [Show],
}

impl<'a> ShowsTemplate<'a> {
    pub async fn render(drv: &DBDriver) -> Result<String, Error> {
        let shows = drv.get_shows().await?;
        Ok(ShowsTemplate { shows: &shows }.render()?)
    }
}
