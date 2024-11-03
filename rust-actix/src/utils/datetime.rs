use chrono::{NaiveDateTime, DateTime, Utc};

pub fn convert_naive_to_utc(naive: NaiveDateTime) -> DateTime<Utc> {
    DateTime::from_naive_utc_and_offset(naive, Utc)
}