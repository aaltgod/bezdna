use anyhow::anyhow;
use regex::bytes::Regex;
use sqlx::PgPool;

use crate::domain;

mod types;

pub struct Repository {
    db: PgPool,
}

impl Repository {
    pub fn new(db: PgPool) -> Self {
        Repository { db }
    }

    pub async fn insert_flags(&self, flags: Vec<domain::Flag>) -> Result<(), anyhow::Error> {
        let qty = flags.len();
        let mut texts: Vec<String> = Vec::with_capacity(qty);
        let mut regexps: Vec<String> = Vec::with_capacity(qty);
        let mut stream_ids: Vec<i64> = Vec::with_capacity(qty);
        let mut directions: Vec<types::FlagDirection> = Vec::with_capacity(qty);

        flags.into_iter().for_each(|flag| {
            texts.push(flag.text);
            regexps.push(flag.regexp.to_string());
            stream_ids.push(flag.stream_id as i64);
            directions.push(types::FlagDirection::from(flag.direction));
        });

        sqlx::query!(
            r#"
            INSERT INTO flags
                (
                    text,
                    regexp,
                    stream_id,
                    direction
                )
            SELECT
                UNNEST($1::text[]) AS text,
                UNNEST($2::text[]) AS regexp,
                UNNEST($3::bigint[]) AS stream_id,
                UNNEST($4::flag_direction[]) AS direction
        "#,
            texts.as_slice(),
            regexps.as_slice(),
            stream_ids.as_slice(),
            directions as _
        )
        .execute(&self.db)
        .await
        .map_err(|e| anyhow!(e.to_string()))?;

        Ok(())
    }

    pub async fn get_flags_by_stream_ids(
        &self,
        stream_ids: Vec<u64>,
    ) -> Result<Vec<domain::Flag>, anyhow::Error> {
        let stream_ids: Vec<i64> = stream_ids.into_iter().map(|id| id as i64).collect();

        let records = match sqlx::query!(
            r#"
        SELECT
            id,
            text,
            regexp,
            stream_id,
            direction AS "flag_direction: types::FlagDirection"
        FROM flags
        WHERE stream_id = ANY(SELECT * FROM UNNEST($1::bigint[]))
        "#,
            stream_ids.as_slice()
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

        let mut flags: Vec<domain::Flag> = Vec::with_capacity(records.len());

        for record in records.into_iter() {
            flags.push(domain::Flag {
                id: record.id as u64,
                text: record.text,
                regexp: Regex::new(record.regexp.as_str()).map_err(|e| anyhow!(e.to_string()))?,
                stream_id: record.stream_id as u64,
                direction: record.flag_direction.into(),
            })
        }

        Ok(flags)
    }
}
