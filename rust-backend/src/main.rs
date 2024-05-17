#[macro_use]
extern crate lazy_static;
#[macro_use]
extern crate log;

use std::time::Duration;

use axum::{
    body::Body,
    Extension,
    http::{Request, StatusCode},
    middleware::{self, Next},
    response::IntoResponse,
    Router, routing::{get, post},
};
use sqlx::postgres::PgPoolOptions;

use handler::{create_service::create_service, get_services::get_services};
use repository::db::postgres::packets as packets_repo;
use repository::db::postgres::services as services_repo;
use repository::db::postgres::streams as streams_repo;
use sniffer::Sniffer;

use crate::handler::types::AppContext;

pub mod domain;
pub mod handler;
pub mod repository;
pub mod sniffer;


#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();

    let pool = PgPoolOptions::new()
        .max_connections(5)
        .acquire_timeout(Duration::from_secs(3))
        .connect("postgresql://localhost/bezdna?user=user&password=1234&port=5433")
        .await
        .unwrap();

    let streams_repo = streams_repo::Repository::new(pool.clone());
    let packets_repo = packets_repo::Repository::new(pool.clone());
    let services_repo = services_repo::Repository::new(pool.clone());

    let app = Router::new()
        .route("/get-services", get(get_services))
        .route("/create-service", post(create_service))
        .layer(middleware::from_fn(info_middleware))
        .layer(Extension(AppContext { services_repo }));

    futures_util::future::join_all(vec![
        tokio::spawn(async move {
            Sniffer::new(
                streams_repo,
                packets_repo,
                "lo0",
                chrono::Duration::seconds(10),
                chrono::Duration::seconds(20),
            )
                .run()
                .await
                .expect("run sniffer")
        }),
        tokio::spawn(async move {
            axum::Server::bind(&"0.0.0.0:3124".parse().unwrap())
                .serve(app.into_make_service())
                .await
                .expect("run server")
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
