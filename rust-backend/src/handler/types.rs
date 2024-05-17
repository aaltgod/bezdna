use serde::{Deserialize, Serialize};

use crate::domain;
use crate::repository::db::postgres::services as services_repo;

#[derive(Clone)]
pub struct AppContext {
    pub services_repo: services_repo::Repository,
}

#[derive(Debug, Serialize, Clone)]
pub struct User {
    pub id: u64,
    pub username: String,
}

#[derive(Deserialize, Debug)]
pub struct CreateService {
    pub name: String,
    pub port: i32,
    pub flag_regexp: String,
}

#[derive(Clone, Serialize, Debug)]
struct Service {
    pub id: i64,
    pub name: String,
    pub port: i32,
    pub flag_regexp: String,
}

impl From<domain::Service> for Service {
    fn from(value: domain::Service) -> Self {
        Service {
            id: value.id as i64,
            name: value.name,
            port: value.port as i32,
            flag_regexp: value.flag_regexp.to_string(),
        }
    }
}