use axum::{Extension, Json};
use axum::http::StatusCode;
use axum::response::IntoResponse;

use crate::handler::types::AppContext;

pub async fn get_services(ctx: Extension<AppContext>) -> Result<impl IntoResponse, StatusCode> {
    let services = ctx.services_repo.get_all_services().await.
        map_err(|e| {
            error!("{}", e.to_string());
            StatusCode::INTERNAL_SERVER_ERROR
        })?;


    Ok(Json("services"))
}



