pub mod domain;
pub mod handler;
pub mod sniffer;

use std::{sync::mpsc::sync_channel, time::Duration};

use axum::{
    body::{Body, HttpBody},
    http::{Request, StatusCode},
    middleware::{self, Next},
    response::IntoResponse,
    routing::{get, post, Route},
    Extension, Router,
};
use handler::{create_service, get_services};
use sniffer::Sniffer;

use sqlx::{postgres::PgPoolOptions, PgPool};
use tokio::sync::{mpsc, mpsc::Sender};

#[macro_use]
extern crate log;

#[derive(Clone)]
pub struct AppContext {
    pub db: PgPool,
    pub tx: Sender<i32>,
}

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();

    let (tx, rx) = mpsc::channel(1);

    let pool = PgPoolOptions::new()
        .max_connections(5)
        .acquire_timeout(Duration::from_secs(3))
        .connect("postgresql://localhost/bezdna?user=user&password=1234&port=5433")
        .await
        .unwrap();

    let app = Router::new()
        .route("/get-services", get(get_services))
        .route("/create-service", post(create_service))
        .layer(middleware::from_fn(info_middleware))
        .layer(Extension(AppContext {
            db: pool.clone(),
            tx,
        }));

    futures_util::future::join_all(vec![
        tokio::spawn(async move { Sniffer::new(pool, "lo").run(rx).await.expect("sniffer") }),
        tokio::spawn(async move {
            axum::Server::bind(&"0.0.0.0:3123".parse().unwrap())
                .serve(app.into_make_service())
                .await
                .expect("server")
        }),
    ])
    .await;
}

async fn info_middleware(
    req: Request<Body>,
    next: Next<Body>,
) -> Result<impl IntoResponse, StatusCode> {
    tracing::info!("{:?}", req.uri());

    Ok(next.run(req).await)
}
