use axum::{Extension, Json};

use crate::handler::types::{AppError, AppContext, Services};

pub async fn get_services(ctx: Extension<AppContext>) -> Result<Json<Services>, AppError> {
    let services = ctx.services_repo.get_all_services().await.map_err(|e| {
        AppError::InternalServerError(e)
    })?;

    Ok(Json(Services::from(services)))
}
