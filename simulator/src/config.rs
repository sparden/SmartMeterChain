use serde::Deserialize;
use std::fs;

#[derive(Debug, Deserialize, Clone)]
pub struct SimConfig {
    pub api: ApiConfig,
    pub meters: Vec<MeterConfig>,
    pub simulation: SimulationConfig,
}

#[derive(Debug, Deserialize, Clone)]
pub struct ApiConfig {
    pub base_url: String,
    pub username: String,
    pub password: String,
}

#[derive(Debug, Deserialize, Clone)]
pub struct MeterConfig {
    pub id: String,
    #[serde(rename = "type")]
    pub meter_type: String,
    pub base_load_kwh: f64,
    pub peak_multiplier: f64,
    pub noise_factor: f64,
}

#[derive(Debug, Deserialize, Clone)]
pub struct SimulationConfig {
    pub interval_seconds: u64,
    pub cumulative: bool,
    pub inject_anomaly_chance: f64,
}

impl SimConfig {
    pub fn load(path: &str) -> Result<Self, Box<dyn std::error::Error>> {
        let content = fs::read_to_string(path)?;
        let config: SimConfig = serde_yaml::from_str(&content)?;
        Ok(config)
    }
}
