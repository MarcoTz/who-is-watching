use axum::{Router, routing::get};
use sqlite::DBDriver;
use std::sync::Arc;
use tower_http::services::ServeDir;

mod routes;

#[tokio::main]
async fn main() {
    let assets = ServeDir::new("assets");
    let drv = Arc::new(DBDriver::connect().await.unwrap());
    let app = Router::new()
        .route("/", get(routes::index))
        .route("/shows", get(routes::shows))
        .route("/watchers", get(routes::watchers))
        .with_state(drv)
        .nest_service("/static", assets);

    let listener = tokio::net::TcpListener::bind("0.0.0.0:8000").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
