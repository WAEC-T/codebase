mod server;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    const PORT: u16 = 5000;
    server::start(PORT).await
}
