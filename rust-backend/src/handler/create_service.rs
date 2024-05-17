use anyhow::anyhow;
use axum::{Extension, Json};
use axum::http::{Error, StatusCode};
use axum::response::{ErrorResponse, IntoResponse};
use regex::bytes;

use crate::domain;
use crate::handler::types::{AppContext, CreateService};
use crate::sniffer::external_types::PORTS_TO_SNIFF;

pub async fn create_service(
    ctx: Extension<AppContext>,
    Json(req): Json<CreateService>,
) -> Result<impl IntoResponse, StatusCode> {
    tracing::info!("{:?}", req);

    let name = req.name;
    let port = req.port as u16;
    let flag_regexp = bytes::Regex::new(req.flag_regexp.as_str())
        .map_err(|e|
            {
                warn!("{}", e.to_string());
                StatusCode::BAD_REQUEST
            })?;

    ctx.services_repo.upsert_service(domain::Service {
        id: 0,
        name,
        port,
        flag_regexp: flag_regexp.clone(),
    }).await.map_err(|e|
        {
            warn!("{}", e.to_string());
            StatusCode::INTERNAL_SERVER_ERROR
        })?;

    PORTS_TO_SNIFF.lock().await.insert(
        port,
        flag_regexp,
    );

    Ok(StatusCode::ACCEPTED)
}
