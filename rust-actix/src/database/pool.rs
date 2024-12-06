use std::error::Error;
use std::{env, time::Duration};

use diesel::{ConnectionError, ConnectionResult};
use diesel_async::{
    pooled_connection::{bb8, AsyncDieselConnectionManager, ManagerConfig},
    AsyncPgConnection,
};

use futures_util::FutureExt;
use tokio_postgres::NoTls;

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
    let mut config = ManagerConfig::default();
    config.custom_setup = Box::new(|url| establish_connection(url).boxed());
    let manager =
        AsyncDieselConnectionManager::<AsyncPgConnection>::new_with_config(database_url, config);
    let pool = bb8::Pool::builder()
        .connection_timeout(Duration::from_secs(15))
        .max_lifetime(Some(Duration::from_secs(60 * 60 * 24)))
        .idle_timeout(Some(Duration::from_secs(60 * 2)))
        .max_size(40)
        .min_idle(Some(20))
        .build(manager)
        .await?;
    println!("Pool status: {:?}", pool.state());
    Ok(DatabasePool(pool))
}

async fn establish_connection(database_url: &str) -> ConnectionResult<AsyncPgConnection> {
    if database_url.contains("localhost")
        || database_url.contains("database")
        || database_url.contains("host.docker.internal")
    {
        let (client, connection) =
            tokio_postgres::connect(database_url, NoTls)
                .await
                .map_err(|e| {
                    ConnectionError::BadConnection(format!(
                        "Error connecting to {}: {}",
                        database_url, e
                    ))
                })?;
        AsyncPgConnection::try_from_client_and_connection(client, connection).await
    } else {
        let tls_connector = native_tls::TlsConnector::builder()
            .danger_accept_invalid_certs(true)
            .build()
            .map_err(|e| {
                ConnectionError::BadConnection(format!("Error building TLS connector: {}", e))
            })?;
        let postgres_tls = postgres_native_tls::MakeTlsConnector::new(tls_connector);
        let (client, connection) = tokio_postgres::connect(database_url, postgres_tls)
            .await
            .map_err(|e| {
                ConnectionError::BadConnection(format!(
                    "Error connecting to {}: {}",
                    database_url, e
                ))
            })?;
        AsyncPgConnection::try_from_client_and_connection(client, connection).await
    }
}
