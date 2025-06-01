use std::sync::Arc;

use crate::database::models::{Messages, Users};
use crate::database::pool::DatabasePool;
use crate::database::repository::{
    create_msg, create_user, follow, get_public_messages, get_timeline, get_user_by_id,
    get_user_by_name, get_user_id, get_user_timeline, is_following, unfollow,
};
use crate::database::PostgresConnection;
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

const PAGE_MESSAGES_LIMIT: i32 = 30;

async fn get_user_template_by_name<'a>(
    conn: &mut PostgresConnection,
    username: &str,
) -> Option<UserTemplate> {
    if let Some(user) = get_user_by_name(conn, username).await {
        Some(UserTemplate {
            user_id: user.user_id,
            username: user.username,
            email: user.email,
        })
    } else {
        None
    }
}

async fn get_user_template(conn: &mut PostgresConnection, user_id: i32) -> Option<UserTemplate> {
    let user = get_user_by_id(conn, user_id).await;
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

async fn get_user(conn: &mut PostgresConnection, session: Session) -> Option<UserTemplate> {
    match session.get::<i32>("user_id") {
        Ok(Some(user_id)) => get_user_template(conn, user_id).await,
        Ok(None) => None,
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
        let gravatar_url_string = gravatar_url(&user.email); // Store the full URL in a variable.
        let pub_date_string = format_datetime_to_message_string(Some(msg.pub_date));

        let message = MessageTemplate {
            text: msg.text,
            username: user.username,
            gravatar_url: gravatar_url_string,
            pub_date: pub_date_string,
        };
        messages_for_template.push(message)
    }
    messages_for_template
}

#[get("/")]
async fn timeline(
    pool: web::Data<Arc<DatabasePool>>,
    flash: Option<FlashMessages>,
    session: Session,
) -> impl Responder {
    let mut conn = pool.get().await.unwrap();
    if let Some(user) = get_user(&mut conn, session).await {
        let messages =
            format_messages(get_timeline(&mut conn, user.user_id, PAGE_MESSAGES_LIMIT).await);

        let rendered = TimelineTemplate {
            messages,
            request_endpoint: "timeline",
            profile_user: None,
            user: Some(user),
            followed: Some(false),
            flashes: flash.unwrap_or_default().messages,
            title: "My Timeline",
            error: "",
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
    pool: web::Data<Arc<DatabasePool>>,
    flash_messages: Option<FlashMessages>,
    session: Session,
) -> impl Responder {
    let mut conn = pool.get().await.unwrap();
    let user = get_user(&mut conn, session).await;
    let messages = get_public_messages(&mut conn, PAGE_MESSAGES_LIMIT).await;
    let messages_for_template = format_messages(messages);

    TimelineTemplate {
        messages: messages_for_template,
        request_endpoint: "public_timeline",
        profile_user: None,
        user,
        followed: Some(false),
        flashes: flash_messages.unwrap_or_default().messages,
        title: "Public Timeline",
        error: "",
    }
}

#[get("/user/{username}")]
async fn user_timeline(
    pool: web::Data<Arc<DatabasePool>>,
    path: web::Path<String>,
    session: Session,
    flash_messages: Option<FlashMessages>,
) -> impl Responder {
    let mut conn = pool.get().await.unwrap();
    let username = path.into_inner();
    let profile_user = get_user_template_by_name(&mut conn, &username).await;
    if let Some(profile_user) = profile_user {
        let mut followed = false;
        let user = get_user(&mut conn, session).await;
        if let Some(ref user) = user {
            followed = is_following(&mut conn, profile_user.user_id, user.user_id).await
        }
        let messages = format_messages(
            get_user_timeline(&mut conn, profile_user.user_id, PAGE_MESSAGES_LIMIT).await,
        );
        let profile_user_name = &profile_user.username;
        let rendered = TimelineTemplate {
            messages,
            request_endpoint: "user_timeline",
            profile_user: Some(&profile_user),
            user,
            followed: Some(followed),
            flashes: flash_messages.unwrap_or_default().messages,
            title: &(profile_user_name.to_string() + "'s Timeline"),
            error: "",
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
    pool: web::Data<Arc<DatabasePool>>,
    path: web::Path<String>,
    _request: HttpRequest,
    session: Session,
) -> impl Responder {
    if let Ok(Some(current_user)) = session.get::<i32>("user_id") {
        let _target_username = &path;
        let mut conn = pool.get().await.unwrap();
        let _target_id = get_user_id(&mut conn, _target_username).await;
        follow(&mut conn, current_user, _target_id).await;
        let message = format!("You are now following {}", _target_username);
        add_flash(&session, message.as_str());
    } else {
        return HttpResponse::Found()
            .append_header((header::LOCATION, "User not found"))
            .finish();
    }
    return HttpResponse::Found()
        .append_header((header::LOCATION, format!("/user/{}", path)))
        .finish();
}

#[get("/{username}/unfollow")]
async fn unfollow_user(
    pool: web::Data<Arc<DatabasePool>>,
    path: web::Path<String>,
    _request: HttpRequest,
    session: Session,
) -> impl Responder {
    if let Ok(Some(current_user)) = session.get::<i32>("user_id") {
        let mut conn = pool.get().await.unwrap();
        let _target_username = &path;
        let _target_id = get_user_id(&mut conn, _target_username).await;
        unfollow(&mut conn, current_user, _target_id).await;

        let message = format!("You are no longer following {}", _target_username);
        add_flash(&session, message.as_str());
    } else {
        return HttpResponse::Found()
            .append_header((header::LOCATION, "User not found"))
            .finish();
    }
    return HttpResponse::Found()
        .append_header((header::LOCATION, format!("/user/{}", path)))
        .finish();
}

#[post("/add_message")]
async fn add_message(
    pool: web::Data<Arc<DatabasePool>>,
    msg: web::Form<MessageInfo>,
    session: Session,
    req: HttpRequest,  // Add HttpRequest to access headers
) -> impl Responder {
    let referer = req.headers().get("referer").and_then(|val| val.to_str().ok());

    match session.get::<i32>("user_id") {
        Ok(Some(user_id)) => {
            let mut conn = pool.get().await.unwrap();
            let timestamp = Utc::now();

            // Get the previous page URL from referer header
            let previous_page: Option<String> = referer.map(|s| s.to_string());

            if msg.text.is_empty() {
                if let Some(user) = get_user(&mut conn, session).await {
                    let messages = format_messages(
                        get_timeline(&mut conn, user.user_id, PAGE_MESSAGES_LIMIT).await,
                    );

                    let context = TimelineTemplate {
                        messages,
                        request_endpoint: "timeline",
                        profile_user: None,
                        user: Some(user),
                        followed: Some(false),
                        flashes: Vec::new(),
                        title: "My Timeline",
                        error: "Message cannot be empty!",
                    }
                    .render()
                    .unwrap();

                    return HttpResponse::Ok().body(context);
                } else {
                    return HttpResponse::TemporaryRedirect()
                        .append_header((header::LOCATION, "/public"))
                        .finish();
                }
            }

            let _ = create_msg(&mut conn, &user_id, &msg.text, timestamp, &0).await;
            add_flash(&session, "Your message was recorded");

            // Redirect to previous page or default to user timeline
            let redirect_url = previous_page.unwrap_or("/".to_string());

            HttpResponse::Found()
                .append_header((header::LOCATION, redirect_url))
                .finish()
        }
        Ok(None) => HttpResponse::Unauthorized()
            .status(StatusCode::UNAUTHORIZED)
            .finish(),
        Err(err) => {
            eprintln!("Failed to retrieve `user_id` from session: {:?}", err);
            HttpResponse::InternalServerError().finish()
        }
    }
}

#[get("/login")]
async fn login(flash_messages: Option<FlashMessages>, session: Session) -> impl Responder {
    match session.get::<i32>("user_id") {
        Ok(Some(_)) => {
            add_flash(&session, "You are already logged in");
            HttpResponse::TemporaryRedirect()
                .append_header((header::LOCATION, "/"))
                .finish()
        }
        Ok(None) => {
            let rendered = LoginTemplate {
                user: None,
                flashes: flash_messages.unwrap_or_default().messages,
                error: "",
                username: "",
            }
            .render()
            .unwrap();
            HttpResponse::Ok().body(rendered)
        }
        Err(err) => {
            eprintln!("Failed to retrieve `user_id` from session: {:?}", err);
            add_flash(&session, "An error occurred while processing your session.");
            HttpResponse::InternalServerError()
                .body("An error occurred while processing your request.")
        }
    }
}

#[post("/login")]
async fn post_login(
    pool: web::Data<Arc<DatabasePool>>,
    info: web::Form<LoginInfo>,
    _request: HttpRequest,
    session: Session,
) -> impl Responder {
    let mut conn = pool.get().await.unwrap();
    let mut error_message: &str = "";

    let result = get_user_by_name(&mut conn, &info.username).await;

    if result.is_none() {
        error_message = "Invalid username";
    } else if let Some(user) = result {
        if !&info.password.eq(&user.pw_hash) {
            error_message = "Invalid password";
        } else {
            session
                .insert("user_id", user.user_id)
                .unwrap_or_else(|err| {
                    eprintln!("Failed to insert user_id into session: {:?}", err);
                });
            add_flash(&session, "You were logged in");

            return HttpResponse::Found()
                .append_header((header::LOCATION, "/"))
                .finish();
        }
    }

    let context = LoginTemplate {
        user: None,
        error: error_message,
        username: &info.username,
        flashes: Vec::new(),
    }
    .render()
    .unwrap();

    HttpResponse::Ok().body(context)
}

#[get("/register")]
async fn register(flash_messages: Option<FlashMessages>) -> impl Responder {
    RegisterTemplate {
        flashes: flash_messages.unwrap_or_default().messages,
        error: "",
        email: "",
        username: "",
        password: "",
        user: None,
    }
}

#[post("/register")]
async fn post_register<'a>(
    pool: web::Data<Arc<DatabasePool>>,
    info: web::Form<RegisterInfo>,
    session: Session,
) -> impl Responder {
    let mut conn = pool.get().await.unwrap();
    let mut error_message = None;

    if info.username.is_empty() {
        error_message = Some("You have to enter a username".to_string());
    } else if info.email.is_empty() || !info.email.contains('@') {
        error_message = Some("You have to enter a valid email address".to_string());
    } else if info.password.is_empty() {
        error_message = Some("You have to enter a password".to_string());
    } else if info.password != info.password2 {
        error_message = Some("The two passwords do not match".to_string());
    } else if get_user_id(&mut conn, &info.username).await != -1 {
        error_message = Some("The username is already taken".to_string());
    }

    if let Some(error) = error_message {
        let context = RegisterTemplate {
            user: None,
            email: &info.email,
            username: &info.username,
            password: &info.password,
            flashes: Vec::new(),
            error: &error,
        }
        .render()
        .unwrap();

        HttpResponse::Ok().body(context)
    } else {
        let _ = create_user(&mut conn, &info.username, &info.email, &info.password).await;

        add_flash(
            &session,
            "You were successfully registered and can login now",
        );

        return HttpResponse::SeeOther()
            .append_header(("Location", "/login"))
            .finish();
    }
}

#[get("/logout")]
async fn logout(session: Session) -> impl Responder {
    add_flash(&session, "You were logged out");
    session.remove("user_id");
    Redirect::to("/public").see_other()
}
