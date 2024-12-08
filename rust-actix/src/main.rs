mod server;
//use env_logger::Env;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    const PORT: u16 = 5000;
    //env_logger::init_from_env(Env::default().default_filter_or("info"));
    server::start(PORT).await
}
