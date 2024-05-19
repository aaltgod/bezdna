use axum::{Extension, Json};
use serde::{Deserialize, Serialize};
use crate::handler::types;
use crate::handler::types::{AppError, Packet};

pub async fn get_streams_by_service_ids(
    ctx: Extension<types::AppContext>,
    Json(req): Json<GetStreamsByServiceIDsRequest>,
) -> Result<Json<GetStreamsByServiceIDsResponse>, AppError> {
    let packets_by_stream_id = ctx.streams_repo
        // TODO: сделать limit настраиваемым.
        .get_packets_by_stream_id(req.service_ids, req.last_stream_id, 50)
        .await
        .map_err(AppError::InternalServerError)?;

    let mut resp = GetStreamsByServiceIDsResponse { stream_ids_with_packets: Vec::with_capacity(packets_by_stream_id.len()) };

    for (stream_id, packets) in packets_by_stream_id {
        resp.stream_ids_with_packets.push(StreamIDWithPackets {
            stream_id,
            packets: packets
                .into_iter()
                .map(|packet| Packet {
                    direction: packet.direction.to_string(),
                    payload: packet.payload,
                    at: packet.at.to_string(),
                })
                .collect(),
        })
    }
    Ok(Json(resp))
}

#[derive(Clone, Debug, Deserialize)]
pub struct GetStreamsByServiceIDsRequest {
    pub service_ids: Vec<i64>,
    pub last_stream_id: i64,
}

#[derive(Clone, Debug, Serialize)]
pub struct GetStreamsByServiceIDsResponse {
    pub(self) stream_ids_with_packets: Vec<StreamIDWithPackets>,
}

#[derive(Clone, Debug, Serialize)]
struct StreamIDWithPackets {
    pub stream_id: i64,
    pub packets: Vec<Packet>,
}