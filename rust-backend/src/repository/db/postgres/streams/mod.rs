use std::collections::BTreeMap;
use anyhow::anyhow;
use sqlx::PgPool;

use crate::domain;
use crate::repository::db::postgres::packets::types;

#[derive(Clone)]
pub struct Repository {
    db: PgPool,
}

impl Repository {
    pub fn new(db: PgPool) -> Self {
        Repository { db }
    }

    pub async fn create_streams(
        &self,
        streams: Vec<domain::Stream>,
    ) -> Result<Vec<u64>, anyhow::Error> {
        let qty = streams.len();
        let mut service_ports: Vec<i32> = Vec::with_capacity(qty);

        streams.into_iter().for_each(|stream| {
            service_ports.push(stream.service_port as u32 as i32);
        });

        let records = sqlx::query!(
            r#"
        INSERT INTO streams
		    (
			    service_port
		    )
		SELECT
		      UNNEST($1::integer[]) AS service_port
        RETURNING id
        "#,
            service_ports.as_slice(),
        )
            .fetch_all(&self.db)
            .await
            .map_err(|e| anyhow!(e.to_string()))?;

        Ok(records.into_iter().map(|record| record.id as u64).collect())
    }

    pub async fn get_packets_by_stream_id(
        &self,
        service_ids: Vec<i64>,
        stream_id: i64,
        limit: i64,
    ) -> Result<BTreeMap<i64, Vec<domain::Packet>>, anyhow::Error> {
        let records = match sqlx::query!(
            r#"
        SELECT
            stream_ids.id AS stream_id,
            packets.direction AS "packet_direction: types::PacketDirection",
            packets.payload,
            packets.at FROM (
                        WITH ports AS (SELECT port
                                       FROM services
                                       WHERE id = ANY($1::bigint[])
                                       )
                        SELECT id FROM streams
                        WHERE
                            service_port = ANY(SELECT port FROM ports)
                        AND
                            id > $2
                        ORDER BY id DESC
                        LIMIT $3
                    ) AS stream_ids
                        INNER JOIN packets ON stream_ids.id = packets.stream_id
        ORDER BY stream_ids.id DESC, packets.at
        "#,
            service_ids as _,
            stream_id,
            limit
        )
            .fetch_all(&self.db)
            .await
        {
            Ok(res) => res,
            Err(e) => {
                return match e {
                    sqlx::Error::RowNotFound => Ok(BTreeMap::default()),
                    _ => Err(anyhow!(e.to_string())),
                };
            }
        };


        let mut result: BTreeMap<i64, Vec<domain::Packet>> = BTreeMap::new();

        records
            .into_iter()
            .for_each(|record| {
                let packet = domain::Packet {
                    id: 0,
                    direction: record.packet_direction.into(),
                    payload: record.payload,
                    stream_id: record.stream_id as u64,
                    at: record.at,
                };

                result
                    .entry(record.stream_id)
                    .and_modify(|packets| packets.push(packet.clone()))
                    .or_insert(vec![packet]);
            });

        Ok(result)
    }
}
