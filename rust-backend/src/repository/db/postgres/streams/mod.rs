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

    pub async fn create_streams(
        &self,
        streams: Vec<domain::Stream>,
    ) -> Result<Vec<i64>, anyhow::Error> {
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
		SELECT * FROM
		        unnest($1::integer[])   AS service_port
        RETURNING id
        "#,
            service_ports.as_slice(),
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
            service_port
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
            })
            .collect())
    }
}
