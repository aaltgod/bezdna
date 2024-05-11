use axum::{Extension, http::StatusCode, Json, response::IntoResponse};
use regex::bytes;
use serde::{Deserialize, Serialize};

use crate::AppContext;
use crate::sniffer::external_types::PORTS_TO_SNIFF;

pub async fn create_service(
    ctx: Extension<AppContext>,
    Json(input): Json<CreateService>,
) -> Result<impl IntoResponse, StatusCode> {
    tracing::info!("{:?}", input);

    PORTS_TO_SNIFF.lock().await.insert(
        input.port as u16,
        bytes::Regex::new(input.flag_regexp.as_str()).unwrap(),
    );

    // match sqlx::query_as!(
    //     Service,
    //     "
    // SELECT id, name, port FROM services
    // WHERE name = $1 AND port = $2
    // ",
    //     input.name,
    //     input.port,
    // )
    // .fetch_optional(&ctx.db)
    // .await
    // {
    //     Ok(res) => match res {
    //         Some(_) => return Err(StatusCode::CONFLICT),
    //         None => (),
    //     },
    //     Err(err) => {
    //         tracing::error!("{}", err);
    //         return Err(StatusCode::INTERNAL_SERVER_ERROR);
    //     }
    // };
    //
    // let service = sqlx::query_as!(
    //     Service,
    //     "
    // INSERT INTO services(name, port)
    //     VALUES ($1, $2)
    // RETURNING id, name, port
    // ",
    //     input.name,
    //     input.port
    // )
    // .fetch_one(&ctx.db)
    // .await
    // .map_err(|e| {
    //     tracing::error!("{}", e);
    //     return StatusCode::INTERNAL_SERVER_ERROR;
    // })?;
    //

    Ok(Json(""))
}

pub async fn get_services(ctx: Extension<AppContext>) -> Result<impl IntoResponse, StatusCode> {
    // let services = sqlx::query_as!(Service, "SELECT id, name, port FROM services")
    //     .fetch_all(&ctx.db)
    //     .await
    //     .map_err(|e| {
    //         tracing::error!("{}", e);
    //         return StatusCode::INTERNAL_SERVER_ERROR;
    //     })?;

    Ok(Json(""))
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
    flag_regexp: String,
}

#[derive(Clone, Serialize, Debug)]
struct Service {
    id: i64,
    name: String,
    port: i32,
}
