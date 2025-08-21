use axum::{extract::State, response::Html};
use pages::{index::IndexTemplate, shows::ShowsTemplate, watchers::WatchersTemplate};
use sqlite::DBDriver;
use std::sync::Arc;

pub async fn index() -> Html<String> {
    Html(IndexTemplate::render().unwrap())
}

pub async fn shows(State(state): State<Arc<DBDriver>>) -> Html<String> {
    Html(ShowsTemplate::render(&state).await.unwrap())
}

pub async fn watchers(State(state): State<Arc<DBDriver>>) -> Html<String> {
    Html(WatchersTemplate::render(&state).await.unwrap())
}
