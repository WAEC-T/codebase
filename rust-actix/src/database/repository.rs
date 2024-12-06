use super::pool::PostgresConnection;
use crate::database::models::*;
use crate::database::schema;
use chrono::DateTime;
use chrono::Utc;
use diesel::dsl::insert_into;
use diesel::{prelude::*, sql_query, sql_types::Integer};
use diesel_async::RunQueryDsl;

pub async fn get_user_id(conn: &mut PostgresConnection, name: &str) -> i32 {
    use schema::users::dsl::*;

    users
        .filter(username.eq(name))
        .select(user_id)
        .first::<i32>(conn)
        .await
        .unwrap_or(-1)
}

pub async fn create_user(
    conn: &mut PostgresConnection,
    _username: &str,
    _email: &str,
    _pw_hash: &str,
) -> Users {
    use schema::users::dsl::*;

    insert_into(users)
        .values((
            username.eq(_username),
            email.eq(_email),
            pw_hash.eq(_pw_hash),
        ))
        .returning(Users::as_returning())
        .get_result(conn)
        .await
        .expect("Error saving new post")
}

pub async fn get_public_messages(
    conn: &mut PostgresConnection,
    limit: i32,
) -> Vec<(Messages, Users)> {
    use schema::messages;
    use schema::users;

    messages::table
        .inner_join(users::table.on(messages::author_id.eq(users::user_id)))
        .order_by(messages::pub_date.desc())
        .limit(limit.into())
        .select((Messages::as_select(), Users::as_select()))
        .load(conn)
        .await
        .expect("Error loading messages and post")
}

pub async fn create_msg(
    conn: &mut PostgresConnection,
    _author_id: &i32,
    _text: &str,
    _pub_date: DateTime<Utc>,
    _flagged: &i32,
) -> Messages {
    use schema::messages::dsl::*;

    insert_into(messages)
        .values((
            author_id.eq(_author_id),
            text.eq(_text),
            pub_date.eq(&_pub_date.naive_utc()),
            flagged.eq(_flagged),
        ))
        .returning(Messages::as_select())
        .get_result(conn)
        .await
        .expect("Error creating new message")
}

pub async fn follow(conn: &mut PostgresConnection, follower_id: i32, followed_id: i32) {
    use schema::followers::dsl::*;
    insert_into(followers)
        .values((who_id.eq(follower_id), whom_id.eq(followed_id)))
        .execute(conn)
        .await
        .expect("Error creating new message");
}

pub async fn unfollow(conn: &mut PostgresConnection, follower_id: i32, followed_id: i32) {
    use schema::followers;
    let _ = diesel::delete(
        followers::table.filter(
            followers::who_id
                .eq(follower_id)
                .and(followers::whom_id.eq(followed_id)),
        ),
    )
    .execute(conn)
    .await;
}

pub async fn get_followers(conn: &mut PostgresConnection, user_id: i32, limit: i32) -> Vec<Users> {
    use schema::followers;
    use schema::users;

    users::table
        .inner_join(followers::table.on(users::user_id.eq(followers::whom_id)))
        .filter(followers::who_id.eq(user_id))
        .select(Users::as_select())
        .limit(limit.into())
        .load(conn)
        .await
        .expect("Couldn't get followers")
}

pub async fn get_user_by_id(conn: &mut PostgresConnection, user_id: i32) -> Option<Users> {
    use schema::users;

    users::table
        .find(user_id)
        .select(Users::as_select())
        .first(conn)
        .await
        .optional()
        .expect("Error fetching user by id")
}

pub async fn get_user_by_name(conn: &mut PostgresConnection, username: &str) -> Option<Users> {
    use schema::users;

    users::table
        .filter(users::username.eq(username))
        .select(Users::as_select())
        .first(conn)
        .await
        .optional()
        .expect("Error fetching user by name")
}

pub async fn get_user_timeline(
    conn: &mut PostgresConnection,
    id: i32,
    limit: i32,
) -> Vec<(Messages, Users)> {
    use schema::messages;
    use schema::users;

    messages::table
        .inner_join(users::table.on(messages::author_id.eq(users::user_id)))
        .filter(messages::flagged.eq(0))
        .filter(users::user_id.eq(id))
        .order_by(messages::pub_date.desc())
        .limit(limit.into())
        .select((Messages::as_select(), Users::as_select()))
        .load(conn)
        .await
        .expect("Error loading messages and post")
}

pub async fn get_timeline(
    conn: &mut PostgresConnection,
    id: i32,
    limit: i32,
) -> Vec<(Messages, Users)> {
    let query = "((SELECT users.user_id, users.username, users.email, users.pw_hash, 
        messages.message_id, messages.author_id, messages.text, messages.pub_date, messages.flagged 
        FROM followers
        INNER JOIN messages ON followers.whom_id = messages.author_id
        INNER JOIN users ON messages.author_id = users.user_id
        WHERE followers.who_id = $1)
        UNION
        (SELECT users.user_id, users.username, users.email, users.pw_hash, 
        messages.message_id, messages.author_id, messages.text, messages.pub_date, messages.flagged 
        FROM messages
        INNER JOIN users ON messages.author_id = users.user_id
        WHERE users.user_id = $1))
        ORDER BY pub_date DESC
        LIMIT $2;
        ";

    sql_query(query)
        .bind::<Integer, _>(id)
        .bind::<Integer, _>(limit)
        .load::<(Messages, Users)>(conn)
        .await
        .expect("")
}

pub async fn is_following(
    conn: &mut PostgresConnection,
    followed_id: i32,
    follower_id: i32,
) -> bool {
    use schema::followers;

    let result: Result<Option<i32>, diesel::result::Error> = followers::table
        .find((follower_id, followed_id))
        .select(followers::who_id)
        .first(conn)
        .await
        .optional();

    result.unwrap().is_some()
}

pub async fn get_latest(conn: &mut PostgresConnection) -> i32 {
    use schema::latest;

    latest::table
        .find(1)
        .select(latest::value)
        .first(conn)
        .await
        .expect("Get latest failed")
}

pub async fn set_latest(conn: &mut PostgresConnection, latest: &i32) {
    use schema::latest;

    let _: usize = diesel::update(latest::table.filter(latest::id.eq(1)))
        .set(latest::value.eq(latest))
        .execute(conn)
        .await
        .expect("Set latest failed");
}
