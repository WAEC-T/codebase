use actix_web::{Scope, web};
use crate::api::routes::*;

pub fn api_services() -> Scope {
    web::scope("/api")
        .service(get_latest)
        .service(post_register)
        .service(get_messages)
        .service(get_messages_per_user)
        .service(post_messages_per_user)
        .service(get_followers)
        .service(post_followers)      
}

// #[actix_web::main]
// pub async fn start() -> std::io::Result<()> {
//     let mut labels = HashMap::new();
//     labels.insert("label1".to_string(), "value1".to_string());

//     HttpServer::new(move || {
//         App::new()
//             .wrap(Logger::default())
//             .service(get_latest)
//             .service(post_register)
//             .service(get_messages)
//             .service(get_messages_per_user)
//             .service(post_messages_per_user)
//             .service(get_followers)
//             .service(post_followers)
//     })
//     .bind(("0.0.0.0", 5001))?
//     .run()
//     .await
// }

