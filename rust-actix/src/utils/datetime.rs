use chrono::{DateTime, NaiveDateTime, Utc};

pub fn convert_naive_to_utc(naive: NaiveDateTime) -> DateTime<Utc> {
    DateTime::from_naive_utc_and_offset(naive, Utc)
}

pub fn format_datetime_to_message_string(timestamp: Option<NaiveDateTime>) -> String {
    match timestamp {
        None => "Unknown date".to_string(),
        Some(t) => t.format("%Y-%m-%d @ %H:%M").to_string(),
    }
}
