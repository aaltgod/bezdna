use axum::{
    http::{Error, StatusCode},
    response::IntoResponse,
    Extension, Json,
};
use serde::{Deserialize, Serialize};

use crate::AppContext;

pub async fn create_service(
    ctx: Extension<AppContext>,
    Json(input): Json<CreateService>,
) -> Result<impl IntoResponse, StatusCode> {
    tracing::info!("{:?}", input);

    let service = Service {
        name: input.name,
        address: input.address,
    };

    Ok(StatusCode::OK)
}

pub async fn get_services(ctx: Extension<AppContext>) -> Result<impl IntoResponse, StatusCode> {
    let mut tx = ctx.db.begin().await.unwrap();

    let rows = match sqlx::query("SELECT * FROM service LIMIT 1")
        .fetch_all(tx.as_mut())
        .await
    {
        Ok(rows) => rows,
        Err(_) => {
            return Err(StatusCode::NOT_FOUND);
        }
    };

    let services: Vec<Service> = vec![];

    Ok(Json(services))
}

#[derive(Debug, Serialize, Clone)]
struct User {
    id: u64,
    username: String,
}

#[derive(Deserialize, Debug)]
pub struct CreateService {
    name: String,
    address: String,
}

#[derive(Clone, Serialize, Debug)]
struct Service {
    name: String,
    address: String,
}
