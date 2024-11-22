use std::env;
use std::error::Error;

use diesel_async::{AsyncPgConnection, pooled_connection::{AsyncDieselConnectionManager,bb8}};

pub type PostgresConnection = AsyncPgConnection;
pub type PostgresPooledConnection<'a> = bb8::PooledConnection<'a, PostgresConnection>;
pub type PostgresPool = bb8::Pool<PostgresConnection>;

#[derive(Clone)]
pub struct DatabasePool(pub PostgresPool);

impl DatabasePool {
    pub async fn get(&self) -> Result<PostgresPooledConnection, bb8::RunError> {
        self.0.get().await
    }
}

fn load_db_url_from_env() -> Result<String, Box<dyn Error>> {
    let database_url = env::var("DATABASE_URL")?;
    Ok(database_url)
}

pub async fn establish_pool() -> Result<DatabasePool, bb8::RunError> {
    let database_url = load_db_url_from_env().expect("Failed to load database URL");
    let manager = AsyncDieselConnectionManager::<AsyncPgConnection>::new(database_url);
    let pool = bb8::Pool::builder().build(manager).await?;
    Ok(DatabasePool(pool))
}