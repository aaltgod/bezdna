use anyhow::anyhow;
use crate::handler::types;
use axum::{Extension, Json};
use regex::bytes;

use crate::domain;
use crate::handler::types::{AppError, AppResponse};
use crate::sniffer::external_types::PORTS_TO_SNIFF;

pub async fn create_service(
    ctx: Extension<types::AppContext>,
    Json(req): Json<types::Service>,
) -> Result<AppResponse, AppError> {
    let name = req.name;
    let port = req.port as u16;
    let flag_regexp = bytes::Regex::new(req.flag_regexp.as_str()).map_err(|e| {
        AppError::BadRequest("Невалидное регулярное выражение".to_string(), anyhow!(e.to_string()))
    })?;

    ctx.services_repo
        .upsert_service(domain::Service {
            id: 0,
            name,
            port,
            flag_regexp: flag_regexp.clone(),
        })
        .await
        .map_err(|e| {
            AppError::InternalServerError(e)
        })?;

    PORTS_TO_SNIFF.lock().await.insert(port, flag_regexp);

    Ok(AppResponse::Created)
}
