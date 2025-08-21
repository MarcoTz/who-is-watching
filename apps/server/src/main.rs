use axum::{Router, response::Html, routing::get};
use pages::{index::IndexTemplate, shows::ShowsTemplate, watchers::WatchersTemplate};
use tower_http::services::ServeDir;

#[tokio::main]
async fn main() {
    let assets = ServeDir::new("assets");
    // build our application with a single route
    let app = Router::new()
        .route(
            "/",
            get(|| async { Html(IndexTemplate::render().unwrap()) }),
        )
        .route(
            "/shows",
            get(|| async { Html(ShowsTemplate::render().unwrap()) }),
        )
        .route(
            "/watchers",
            get(|| async { Html(WatchersTemplate::render().unwrap()) }),
        )
        .nest_service("/static", assets);

    // run our app with hyper, listening globally on port 3000
    let listener = tokio::net::TcpListener::bind("0.0.0.0:8000").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
