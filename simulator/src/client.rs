use reqwest::Client;
use serde_json::Value;
use std::sync::Arc;
use tokio::sync::RwLock;

use crate::config::{ApiConfig, MeterConfig};
use crate::generator::MeterReading;

pub struct ApiClient {
    client: Client,
    base_url: String,
    username: String,
    password: String,
    token: Arc<RwLock<String>>,
}

impl ApiClient {
    pub fn new(config: &ApiConfig) -> Result<Self, Box<dyn std::error::Error>> {
        let client = Client::builder()
            .timeout(std::time::Duration::from_secs(30))
            .build()?;

        Ok(Self {
            client,
            base_url: config.base_url.clone(),
            username: config.username.clone(),
            password: config.password.clone(),
            token: Arc::new(RwLock::new(String::new())),
        })
    }

    pub async fn login(&self) -> Result<(), Box<dyn std::error::Error>> {
        let url = format!("{}/auth/login", self.base_url);
        let body = serde_json::json!({
            "username": self.username,
            "password": self.password,
        });

        let resp = self.client.post(&url).json(&body).send().await?;
        let data: Value = resp.json().await?;

        if let Some(token) = data["data"]["token"].as_str() {
            let mut t = self.token.write().await;
            *t = token.to_string();
            Ok(())
        } else {
            Err(format!("Login failed: {:?}", data).into())
        }
    }

    pub async fn register_meter(
        &self,
        meter: &MeterConfig,
    ) -> Result<(), Box<dyn std::error::Error>> {
        let url = format!("{}/meters", self.base_url);
        let token = self.token.read().await;

        let body = serde_json::json!({
            "meter_id": meter.id,
            "consumer_id": "consumer1",
            "meter_type": meter.meter_type,
            "location": "Simulated",
            "status": "active",
        });

        let resp = self
            .client
            .post(&url)
            .bearer_auth(&*token)
            .json(&body)
            .send()
            .await?;

        if resp.status().is_success() {
            Ok(())
        } else {
            let text = resp.text().await?;
            Err(format!("Register failed: {}", text).into())
        }
    }

    pub async fn submit_reading(
        &self,
        reading: &MeterReading,
    ) -> Result<String, Box<dyn std::error::Error>> {
        let url = format!("{}/meters/readings", self.base_url);
        let token = self.token.read().await;

        let resp = self
            .client
            .post(&url)
            .bearer_auth(&*token)
            .json(reading)
            .send()
            .await;

        match resp {
            Ok(r) if r.status().is_success() => {
                let data: Value = r.json().await?;
                let tx_id = data["tx_id"]
                    .as_str()
                    .unwrap_or("unknown")
                    .to_string();
                Ok(tx_id)
            }
            Ok(r) => {
                let status = r.status();
                let text = r.text().await.unwrap_or_default();

                // Retry with re-login on 401
                if status.as_u16() == 401 {
                    drop(token);
                    self.login().await?;
                    return self.submit_reading_inner(reading).await;
                }
                Err(format!("HTTP {}: {}", status, text).into())
            }
            Err(e) => Err(e.into()),
        }
    }

    async fn submit_reading_inner(
        &self,
        reading: &MeterReading,
    ) -> Result<String, Box<dyn std::error::Error>> {
        let url = format!("{}/meters/readings", self.base_url);
        let token = self.token.read().await;

        let resp = self
            .client
            .post(&url)
            .bearer_auth(&*token)
            .json(reading)
            .send()
            .await?;

        let data: Value = resp.json().await?;
        let tx_id = data["tx_id"]
            .as_str()
            .unwrap_or("unknown")
            .to_string();
        Ok(tx_id)
    }
}
