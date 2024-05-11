use anyhow::anyhow;
use regex::bytes;
use sqlx::PgPool;

use crate::domain;

pub struct Repository {
    db: PgPool,
}

impl Repository {
    pub fn new(db: PgPool) -> Self {
        Repository { db }
    }

    pub async fn upsert_service(&self, service: domain::Service) -> Result<(), anyhow::Error> {
        sqlx::query!(
            r#"
        INSERT INTO services(name, port, flag_regexp)
            VALUES($1, $2, $3)
                ON CONFLICT ON CONSTRAINT services_port_key
	                DO UPDATE SET
				        name=EXCLUDED.name,
				        flag_regexp=EXCLUDED.flag_regexp
        "#,
            service.name,
            service.port as u32 as i32,
            service.flag_regexp.to_string(),
        )
        .execute(&self.db)
        .await
        .map_err(|e| anyhow!(e.to_string()))?;

        Ok(())
    }

    pub async fn get_service_by_port(
        &self,
        port: u16,
    ) -> Result<Option<domain::Service>, anyhow::Error> {
        let record = match sqlx::query!(
            r#"
        SELECT id, name, port, flag_regexp
        FROM services
            WHERE port = $1
        "#,
            port as u32 as i32
        )
        .fetch_one(&self.db)
        .await
        {
            Ok(res) => res,
            Err(e) => {
                return match e {
                    sqlx::Error::RowNotFound => Ok(None),
                    _ => Err(anyhow!(e.to_string())),
                }
            }
        };

        let service = domain::Service {
            id: record.id as u64,
            name: record.name,
            port: record.port as u16,
            flag_regexp: bytes::Regex::new(record.flag_regexp.as_str())
                .map_err(|e| anyhow!(e.to_string()))?,
        };

        Ok(Some(service))
    }

    pub async fn get_all_services(&self) -> Result<Vec<domain::Service>, anyhow::Error> {
        let records = match sqlx::query!(
            r#"
        SELECT id, name, port, flag_regexp
        FROM services
        "#
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

        let mut services: Vec<domain::Service> = Vec::with_capacity(records.len());

        for record in records.into_iter() {
            services.push(domain::Service {
                id: record.id as u64,
                name: record.name,
                port: record.port as u16,
                flag_regexp: bytes::Regex::new(record.flag_regexp.as_str())
                    .map_err(|e| anyhow!(e.to_string()))?,
            });
        }

        Ok(services)
    }
}
