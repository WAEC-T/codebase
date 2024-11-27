use crate::database::models::{Messages, Users};
use crate::database::pool::DatabasePool;
use crate::database::repository::{
    create_msg, create_user, follow, get_passwd_hash, get_public_messages, get_timeline,
    get_user_by_id, get_user_by_name, get_user_timeline, is_following, unfollow,
};
use crate::frontend::flash_messages::*;
use crate::frontend::template_structs::*;
use crate::utils::datetime::format_datetime_to_message_string;
use actix_session::Session;
use actix_web::{
    get,
    http::{header, StatusCode},
    post,
    web::{self, Redirect},
    HttpRequest, HttpResponse, Responder,
};
use askama_actix::Template;
use chrono::Utc;
use md5::{Digest, Md5};
use pwhash::bcrypt;

const PAGE_MESSAGES_LIMIT: i32 = 30;

async fn get_user_id(pool: web::Data<DatabasePool>, username: &str) -> i32 {
    let mut conn = pool.get().await.unwrap();
    let user = get_user_by_name(&mut conn, username).await;
    if let Some(user) = user {
        user.user_id
    } else {
        -1
    }
}

async fn get_user_template_by_name(
    pool: web::Data<DatabasePool>,
    username: &str,
) -> Option<UserTemplate> {
    let mut conn = pool.get().await.unwrap();
    let user = get_user_by_name(&mut conn, username).await;
    if let Some(user) = user {
        Some(UserTemplate {
            user_id: user.user_id,
            username: user.username,
            email: user.email,
        })
    } else {
        None
    }
}

async fn get_user_template(pool: web::Data<DatabasePool>, user_id: i32) -> Option<UserTemplate> {
    let mut conn = pool.get().await.unwrap();
    let user = get_user_by_id(&mut conn, user_id).await;
    if let Some(user) = user {
        Some(UserTemplate {
            user_id: user.user_id,
            username: user.username,
            email: user.email,
        })
    } else {
        None
    }
}

async fn get_user(
    pool: web::Data<DatabasePool>,
    session: Session,
) -> Option<UserTemplate> {
    match session.get::<i32>("user_id") {
        Ok(Some(user_id)) => {
        get_user_template(pool, user_id).await
    }
    Ok(None) => {
        None
    }
    Err(err) => {
        eprintln!("Failed to retrieve `user_id` from session: {:?}", err);
        None
    }
    }
}

fn gravatar_url(email: &str) -> String {
    let hash = Md5::digest(email.trim().to_lowercase().as_bytes());

    let hash_str = format!("{:x}", hash);

    format!(
        "https://www.gravatar.com/avatar/{}?d=identicon&s={}",
        hash_str, 48
    )
}

fn format_messages(messages: Vec<(Messages, Users)>) -> Vec<MessageTemplate> {
    let mut messages_for_template: Vec<MessageTemplate> = Vec::new();
    for (msg, user) in messages {
        let message = MessageTemplate {
            text: msg.text,
            username: user.username,
            gravatar_url: gravatar_url(&user.email),
            pub_date: format_datetime_to_message_string(Some(msg.pub_date)),
        };
        messages_for_template.push(message)
    }
    messages_for_template
}

#[get("/")]
async fn timeline(
    pool: web::Data<DatabasePool>,
    flash: Option<FlashMessages>,
    session: Session
) -> impl Responder {
    let mut conn = pool.get().await.unwrap();
    if let Some(user) = get_user(pool.clone(), session).await {
        let messages =
            format_messages(get_timeline(&mut conn, user.user_id, PAGE_MESSAGES_LIMIT).await);

        let rendered = TimelineTemplate {
            messages,
            request_endpoint: "timeline",
            profile_user: None,
            user: Some(user),
            followed: Some(false),
            flashes: flash.unwrap_or_default().messages,
            title: String::from("My Timeline"),
        }
        .render()
        .unwrap();
        HttpResponse::Ok().body(rendered)
    } else {
        HttpResponse::TemporaryRedirect()
            .append_header((header::LOCATION, "/public"))
            .finish()
    }
}

#[get("/public")]
async fn public_timeline(
    pool: web::Data<DatabasePool>,
    flash_messages: Option<FlashMessages>,
    session: Session,
) -> impl Responder {
    let user = get_user(pool.clone(), session).await;
    let mut conn = pool.get().await.unwrap();
    let messages = get_public_messages(&mut conn, PAGE_MESSAGES_LIMIT).await;
    let messages_for_template = format_messages(messages);

    TimelineTemplate {
        messages: messages_for_template,
        request_endpoint: "public_timeline",
        profile_user: None,
        user,
        followed: Some(false),
        flashes: flash_messages.unwrap_or_default().messages,
        title: String::from("Public Timeline"),
    }
}

#[get("/{username}")]
async fn user_timeline(
    pool: web::Data<DatabasePool>,
    path: web::Path<String>,
    session: Session,
    flash_messages: Option<FlashMessages>,
) -> impl Responder {
    let username = path.into_inner();
    let profile_user = get_user_template_by_name(pool.clone(), &username).await;
    if let Some(profile_user) = profile_user {
        let mut followed = false;
        let user = get_user(pool.clone(), session).await;
        let mut conn = pool.get().await.unwrap();
        if let Some(user) = user.clone() {
            followed = is_following(&mut conn, profile_user.clone().user_id, user.user_id).await
        }
        let messages = format_messages(
            get_user_timeline(&mut conn, profile_user.user_id, PAGE_MESSAGES_LIMIT).await,
        );
        let profile_user_name = profile_user.clone().username;
        let rendered = TimelineTemplate {
            messages,
            request_endpoint: "user_timeline",
            profile_user: Some(profile_user.clone()),
            user,
            followed: Some(followed),
            flashes: flash_messages.unwrap_or_default().messages,
            title: format!("{}'s Timeline", profile_user_name),
        }
        .render()
        .unwrap();
        HttpResponse::Ok().body(rendered)
    } else {
        HttpResponse::NotFound().finish()
    }
}

#[get("/{username}/follow")]
async fn follow_user(
    pool: web::Data<DatabasePool>,
    path: web::Path<String>,
    _request: HttpRequest,
    session: Session,
) -> impl Responder {
    if let Ok(Some(current_user)) = session.get::<i32>("user_id") {
        let _target_username = path.clone();
        let _target_id = get_user_id(pool.clone(), &_target_username).await;
        let mut conn = pool.get().await.unwrap();
        follow(&mut conn, current_user, _target_id).await;
        let message = format!("You are now following {}", _target_username);
        add_flash(session, message.as_str());
    } else {
        return HttpResponse::Found()
            .append_header((header::LOCATION, "User not found"))
            .finish();
    }
    return HttpResponse::Found()
        .append_header((header::LOCATION, format!("/{}", path)))
        .finish();
}

#[get("/{username}/unfollow")]
async fn unfollow_user(
    pool: web::Data<DatabasePool>,
    path: web::Path<String>,
    _request: HttpRequest,
    session: Session,
) -> impl Responder {
    if let Ok(Some(current_user)) = session.get::<i32>("user_id") {
        let _target_username = path.clone();
        let _target_id = get_user_id(pool.clone(), &_target_username).await;
        let mut conn = pool.get().await.unwrap();
        unfollow(&mut conn, current_user, _target_id).await;

        let message = format!("You are no longer following {}", _target_username);
        add_flash(session, message.as_str());
    } else {
        return HttpResponse::Found()
            .append_header((header::LOCATION, "User not found"))
            .finish();
    }
    return HttpResponse::Found()
        .append_header((header::LOCATION, format!("/{}", path)))
        .finish();
}

#[post("/add_message")]
async fn add_message(
    pool: web::Data<DatabasePool>,
    msg: web::Form<MessageInfo>,
    session: Session,
) -> impl Responder {
    match session.get::<i32>("user_id") {
        Ok(Some(user_id)) => {
            let mut conn = pool.get().await.unwrap();
            let timestamp = Utc::now();
            let _ = create_msg(&mut conn, &user_id, &msg.text, timestamp, &0).await;
            add_flash(session, "Your message was recorded");
            HttpResponse::Found()
                .append_header((header::LOCATION, "/"))
                .finish()
        }
        Ok(None) => {
            HttpResponse::Unauthorized()
                .status(StatusCode::UNAUTHORIZED)
                .finish()
        }
        Err(err) => {
            eprintln!("Failed to retrieve `user_id` from session: {:?}", err);
            HttpResponse::InternalServerError().finish()
        }
    }
}

#[get("/login")]
async fn login(
    flash_messages: Option<FlashMessages>,
    session: Session,
) -> impl Responder {
    match session.get::<i32>("user_id") {
        Ok(Some(_)) => {
        add_flash(session, "You are already logged in");
        HttpResponse::TemporaryRedirect()
            .append_header((header::LOCATION, "/"))
            .finish()
        }
        Ok(None) => {
            let rendered = LoginTemplate {
                user: None,
                flashes: flash_messages.unwrap_or_default().messages,
                error: String::from(""),
                username: String::from(""),
            }
            .render()
            .unwrap();
            HttpResponse::Ok().body(rendered)
        }
        Err(err) => {
            eprintln!("Failed to retrieve `user_id` from session: {:?}", err);
            add_flash(session, "An error occurred while processing your session.");
            HttpResponse::InternalServerError()
                .body("An error occurred while processing your request.")
        }
    }
}

#[post("/login")]
async fn post_login(
    pool: web::Data<DatabasePool>,
    info: web::Form<LoginInfo>,
    _request: HttpRequest,
    session: Session,
) -> impl Responder {
    let mut conn = pool.get().await.unwrap();
    let result = get_passwd_hash(&mut conn, &info.username).await;
    if result.is_none() {
        add_flash(session, "Invalid username");
        return HttpResponse::Found()
            .append_header((header::LOCATION, "/login"))
            .finish();
    }

    if let Some(stored_hash) = result {
        if bcrypt::verify(info.password.clone(), &stored_hash) {
            // Successful login
            let user_id = get_user_id(pool.clone(), &info.username).await;
            session.insert("user_id", user_id).unwrap_or_else(|err| {
                eprintln!("Failed to insert user_id into session: {:?}", err);
            });
            add_flash(session, "You were logged in");

            return HttpResponse::Found()
                .append_header((header::LOCATION, "/"))
                .finish();
        }
    }

    // Password incorrect
    add_flash(session, "Invalid password");
    return HttpResponse::Found()
        .append_header((header::LOCATION, "/login"))
        .finish();
}

#[get("/register")]
async fn register(flash_messages: Option<FlashMessages>) -> impl Responder {
    RegisterTemplate {
        flashes: flash_messages.unwrap_or_default().messages,
        error: String::from(""),
        email: String::from(""),
        username: String::from(""),
        password: String::from(""),
        user: None,
    }
}

#[post("/register")]
async fn post_register(
    pool: web::Data<DatabasePool>,
    info: web::Form<RegisterInfo>,
    session: Session,
) -> impl Responder {
    if info.username.is_empty() {
        add_flash(session, "You have to enter a username");
        return Redirect::to("/register").see_other();
    } else if info.email.is_empty() || !info.email.contains('@') {
        add_flash(session, "You have to enter a valid email address");
        return Redirect::to("/register").see_other();
    } else if info.password.is_empty() {
        add_flash(session, "You have to enter a password");
        return Redirect::to("/register").see_other();
    } else if info.password != info.password2 {
        add_flash(session, "The two passwords do not match");
        return Redirect::to("/register").see_other();
    } else if get_user_id(pool.clone(), &info.username).await != -1 {
        add_flash(session, "The username is already taken");
        return Redirect::to("/register").see_other();
    }

    let hash = bcrypt::hash(info.password.clone()).unwrap();

    let mut conn = pool.get().await.unwrap();
    let _ = create_user(&mut conn, &info.username, &info.email, &hash).await;

    add_flash(
        session,
        "You were successfully registered and can login now",
    );
    Redirect::to("/login").see_other()
}

#[get("/logout")]
async fn logout(session: Session) -> impl Responder {
    add_flash(session.clone(), "You were logged out.");
    session.remove("user_id");
    Redirect::to("/public").see_other()
}
