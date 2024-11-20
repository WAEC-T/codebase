use std::future::{ready, Ready};

use actix_web::{
    dev::{forward_ready, Service, ServiceRequest, ServiceResponse, Transform},
    http::header,
    Error,
};
use futures_util::future::LocalBoxFuture;

pub struct AuthMiddleware;

impl<S, B> Transform<S, ServiceRequest> for AuthMiddleware
where
    S: Service<ServiceRequest, Response = ServiceResponse<B>, Error = Error>,
    S::Future: 'static,
    B: 'static,
{
    type Response = ServiceResponse<B>;
    type Error = Error;
    type InitError = ();
    type Transform = VerifyAuthToken<S>;
    type Future = Ready<Result<Self::Transform, Self::InitError>>;

    fn new_transform(&self, service: S) -> Self::Future {
        ready(Ok(VerifyAuthToken { service }))
    }
}
pub struct VerifyAuthToken<S> {
    service: S,
}

impl<S, B> Service<ServiceRequest> for VerifyAuthToken<S>
where
    S: Service<ServiceRequest, Response = ServiceResponse<B>, Error = Error>,
    S::Future: 'static,
    B: 'static,
{
    type Response = ServiceResponse<B>;
    type Error = Error;
    type Future = LocalBoxFuture<'static, Result<Self::Response, Self::Error>>;

    forward_ready!(service);

    fn call(&self, req: ServiceRequest) -> Self::Future {
        if req.path().starts_with("/api")
            && req
                .headers()
                .get(header::AUTHORIZATION)
                .map_or(true, |auth_header| {
                    auth_header != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh"
                })
        {
            return Box::pin(async { Err(actix_web::error::ErrorUnauthorized("Not Authorized")) });
        }

        let fut: <S as Service<ServiceRequest>>::Future = self.service.call(req);

        Box::pin(async move {
            let res = fut.await?;
            Ok(res)
        })
    }
}
