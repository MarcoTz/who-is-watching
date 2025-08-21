use crate::errors::Error;
use askama::Template;

#[derive(Template)]
#[template(path = "watchers.html")]
pub struct WatchersTemplate {}

impl WatchersTemplate {
    pub fn render() -> Result<String, Error> {
        Ok(WatchersTemplate {}.render()?)
    }
}
