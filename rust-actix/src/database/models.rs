use crate::database::schema::{followers, messages, users};
use chrono::NaiveDateTime;
use diesel::prelude::*;

#[derive(Queryable, Selectable)]
#[diesel(table_name = crate::database::schema::users)]
#[diesel(check_for_backend(diesel::pg::Pg))]
#[derive(QueryableByName)]
pub struct Users {
    pub user_id: i32,
    pub username: String,
    pub email: String,
    pub pw_hash: String,
}

#[derive(Insertable)]
#[diesel(table_name = users)]
pub struct NewUser<'a> {
    pub username: &'a str,
    pub email: &'a str,
    pub pw_hash: &'a str,
}

#[derive(Queryable, Selectable, Insertable)]
#[diesel(table_name = crate::database::schema::followers)]
#[diesel(check_for_backend(diesel::pg::Pg))]
#[derive(QueryableByName)]
pub struct Followers {
    pub who_id: i32,
    pub whom_id: i32,
}

#[derive(Insertable)]
#[diesel(table_name = followers)]
pub struct NewFollower<'a> {
    pub who_id: &'a i32,
    pub whom_id: &'a i32,
}

#[derive(Queryable, QueryableByName, Selectable, Insertable)]
#[diesel(table_name = crate::database::schema::messages)]
#[diesel(check_for_backend(diesel::pg::Pg))]
pub struct Messages {
    pub message_id: i32,
    pub author_id: i32,
    pub text: String,
    pub pub_date: NaiveDateTime,
    pub flagged: i32,
}

#[derive(Insertable)]
#[diesel(table_name = messages)]
pub struct NewMessage<'a> {
    pub author_id: &'a i32,
    pub text: &'a str,
    pub pub_date: &'a NaiveDateTime,
    pub flagged: &'a i32,
}

#[derive(Queryable)]
#[diesel(table_name = latest)]
pub struct Latest {
    pub id: i32,
    pub value: i32,
}