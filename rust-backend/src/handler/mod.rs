use crate::AppContext;
use axum::{http::StatusCode, response::IntoResponse, Extension, Json};
use serde::{Deserialize, Serialize};

pub async fn create_service(
    ctx: Extension<AppContext>,
    Json(input): Json<CreateService>,
) -> Result<impl IntoResponse, StatusCode> {
    tracing::info!("{:?}", input);

    match sqlx::query_as!(
        Service,
        "
    SELECT id, name, port FROM services
    WHERE name = $1 AND port = $2
    ",
        input.name,
        input.port,
    )
    .fetch_optional(&ctx.db)
    .await
    {
        Ok(res) => match res {
            Some(_) => return Err(StatusCode::CONFLICT),
            None => (),
        },
        Err(err) => {
            tracing::error!("{}", err);
            return Err(StatusCode::INTERNAL_SERVER_ERROR);
        }
    };

    let service = sqlx::query_as!(
        Service,
        "
    INSERT INTO services(name, port)
        VALUES ($1, $2)
    RETURNING id, name, port
    ",
        input.name,
        input.port
    )
    .fetch_one(&ctx.db)
    .await
    .map_err(|e| {
        tracing::error!("{}", e);
        return StatusCode::INTERNAL_SERVER_ERROR;
    })?;

    ctx.tx.send(service.port as u16).await.unwrap();

    Ok(Json(service))
}

pub async fn get_services(ctx: Extension<AppContext>) -> Result<impl IntoResponse, StatusCode> {
    let services = sqlx::query_as!(Service, "SELECT id, name, port FROM services")
        .fetch_all(&ctx.db)
        .await
        .map_err(|e| {
            tracing::error!("{}", e);
            return StatusCode::INTERNAL_SERVER_ERROR;
        })?;

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
    port: i32,
}

#[derive(Clone, Serialize, Debug)]
struct Service {
    id: i64,
    name: String,
    port: i32,
}
