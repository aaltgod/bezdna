use sqlx::postgres::{PgHasArrayType, PgTypeInfo};

use crate::domain;

#[derive(Debug, sqlx::Type)]
#[sqlx(type_name = "flag_direction")]
pub enum FlagDirection {
    IN,
    OUT,
}

impl From<domain::FlagDirection> for FlagDirection {
    fn from(value: domain::FlagDirection) -> Self {
        match value {
            domain::FlagDirection::IN => FlagDirection::IN,
            domain::FlagDirection::OUT => FlagDirection::OUT,
        }
    }
}

impl Into<domain::FlagDirection> for FlagDirection {
    fn into(self) -> domain::FlagDirection {
        match self {
            FlagDirection::IN => domain::FlagDirection::IN,
            FlagDirection::OUT => domain::FlagDirection::OUT,
        }
    }
}

impl PgHasArrayType for FlagDirection {
    fn array_type_info() -> PgTypeInfo {
        PgTypeInfo::with_name("_flag_direction")
    }
}
