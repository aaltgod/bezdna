#[macro_use]
extern crate lazy_static;
#[macro_use]
extern crate log;

use std::{sync::mpsc::sync_channel, time, time::Duration};

use axum::{
    body::{Body, HttpBody},
    Extension,
    http::{Request, StatusCode},
    middleware::{self, Next},
    response::IntoResponse,
    Router, routing::{get, post, Route},
};
use sqlx::{PgPool, postgres::PgPoolOptions};
use tokio::sync::{mpsc, mpsc::Sender};

use handler::{create_service, get_services};
use sniffer::Sniffer;

pub mod domain;
pub mod handler;
pub mod sniffer;

#[derive(Clone)]
pub struct AppContext {
    pub db: PgPool,
    pub tx: Sender<u16>,
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

    tx.send(3123).await.unwrap();

    let app = Router::new()
        .route("/get-services", get(get_services))
        .route("/create-service", post(create_service))
        .layer(middleware::from_fn(info_middleware))
        .layer(Extension(AppContext {
            db: pool.clone(),
            tx,
        }));

    futures_util::future::join_all(vec![
        tokio::spawn(async move {
            Sniffer::new(pool, "lo", chrono::Duration::seconds(10))
                .run(rx)
                .await
                .expect("sniffer")
        }),
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
