use anyhow::anyhow;
use sqlx::PgPool;

use crate::domain;

pub struct Repository {
    db: PgPool,
}

impl Repository {
    pub fn new(db: PgPool) -> Self {
        Repository { db }
    }

    pub async fn upsert_streams(
        &self,
        streams: Vec<domain::Stream>,
    ) -> Result<Vec<i64>, anyhow::Error> {
        let qty = streams.len();
        let mut service_ports: Vec<i32> = Vec::with_capacity(qty);
        let mut payloads: Vec<String> = Vec::with_capacity(qty);
        let mut started_at: Vec<chrono::NaiveDateTime> = Vec::with_capacity(qty);
        let mut ended_at: Vec<chrono::NaiveDateTime> = Vec::with_capacity(qty);

        streams.into_iter().for_each(|stream| {
            service_ports.push(stream.service_port as u32 as i32);
            payloads.push(stream.payload);
            started_at.push(stream.started_at.naive_utc());
            ended_at.push(stream.ended_at.naive_utc());
        });

        let records = sqlx::query!(
            r#"
        INSERT INTO streams
		(
			service_port,
			payload,
			started_at,
			ended_at
		)
		SELECT * FROM
		        unnest($1::integer[])   AS service_port,
		        unnest($2::text[])      AS payload,
		        unnest($3::timestamp[]) AS started_at,
		        unnest($4::timestamp[]) AS ended_at
        RETURNING id
        "#,
            service_ports.as_slice(),
            payloads.as_slice(),
            started_at.as_slice(),
            ended_at.as_slice()
        )
        .fetch_all(&self.db)
        .await
        .map_err(|e| anyhow!(e.to_string()))?;

        Ok(records.into_iter().map(|record| record.id).collect())
    }

    pub async fn get_last_streams(&self, limit: u64) -> Result<Vec<domain::Stream>, anyhow::Error> {
        let records = match sqlx::query!(
            r#"
        SELECT
            id,
            service_port,
            payload AS "payload!",
            started_at AS "started_at!",
            ended_at AS "ended_at!"
        FROM streams
        ORDER BY id DESC
        LIMIT $1
        "#,
            limit as i64
        )
        .fetch_all(&self.db)
        .await
        {
            Ok(res) => res,
            Err(e) => {
                return match e {
                    sqlx::Error::RowNotFound => Ok(vec![]),
                    _ => Err(anyhow!(e.to_string())),
                }
            }
        };

        Ok(records
            .into_iter()
            .map(|record| domain::Stream {
                id: record.id as u64,
                service_port: record.service_port as u16,
                payload: record.payload,
                started_at: record.started_at,
                ended_at: record.ended_at,
            })
            .collect())
    }
}
