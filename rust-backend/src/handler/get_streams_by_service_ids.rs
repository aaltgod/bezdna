use axum::{Extension, Json};
use serde::{Deserialize, Serialize};
use crate::domain;
use crate::handler::types;
use crate::handler::types::{AppError, Packet, StreamIDWithPackets};

pub async fn get_streams_by_service_ids(
    ctx: Extension<types::AppContext>,
    Json(req): Json<GetStreamsByServiceIDsRequest>,
) -> Result<Json<GetStreamsByServiceIDsResponse>, AppError> {
    let services = ctx.services_repo
        .get_services_by_ids(req.service_ids)
        .await.
        map_err(AppError::InternalServerError)?;

    let services_slice = services.as_slice();

    let packets_by_stream = ctx.streams_repo
        .get_packets_by_stream(
            services_slice
                .iter()
                .map(|service| service.port)
                .collect(),
            req.last_stream_id,
            // TODO: сделать limit настраиваемым.
            20,
        )
        .await
        .map_err(AppError::InternalServerError)?;

    let mut resp = GetStreamsByServiceIDsResponse { stream_ids_with_packets: Vec::with_capacity(packets_by_stream.len()) };

    for (stream, packets) in packets_by_stream {
        resp.stream_ids_with_packets.push(StreamIDWithPackets {
            stream_id: stream.id,
            packets: packets
                .into_iter()
                .map(|packet| Packet {
                    direction: packet.direction.to_string(),
                    payload: packet.payload,
                    at: packet.at.to_string(),
                    flag_regexp: services_slice
                        .iter()
                        .find(|service| service.port.eq(&stream.service_port))
                        .map_or("".to_string(), |s| s.flag_regexp.to_string()),
                    color: match packet.direction {
                        domain::PacketDirection::IN => "#33FF46".to_string(),
                        domain::PacketDirection::OUT => "#FF3333".to_string(),
                    },
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
    pub stream_ids_with_packets: Vec<StreamIDWithPackets>,
}