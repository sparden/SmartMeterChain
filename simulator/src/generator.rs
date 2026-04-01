use chrono::Utc;
use rand::Rng;
use serde::Serialize;

use crate::config::MeterConfig;

#[derive(Debug, Serialize)]
pub struct MeterReading {
    pub meter_id: String,
    pub reading: f64,
    pub timestamp: i64,
    pub temperature: f64,
    pub voltage: f64,
    #[serde(skip)]
    pub is_anomaly: bool,
}

pub struct MeterGenerator {
    config: MeterConfig,
    cumulative_reading: f64,
    anomaly_chance: f64,
}

impl MeterGenerator {
    pub fn new(config: MeterConfig, anomaly_chance: f64) -> Self {
        Self {
            config,
            cumulative_reading: 0.0,
            anomaly_chance,
        }
    }

    pub fn next_reading(&mut self) -> MeterReading {
        let mut rng = rand::thread_rng();
        let now = Utc::now();
        let hour = now.format("%H").to_string().parse::<u32>().unwrap_or(12);

        // Time-of-day consumption pattern
        let time_multiplier = match self.config.meter_type.as_str() {
            "domestic" => match hour {
                6..=9 => 1.5,   // Morning peak
                10..=16 => 0.6, // Away at work
                17..=22 => self.config.peak_multiplier, // Evening peak
                _ => 0.3,       // Night (low)
            },
            "commercial" => match hour {
                9..=18 => self.config.peak_multiplier, // Business hours
                _ => 0.2,
            },
            "industrial" => match hour {
                6..=22 => self.config.peak_multiplier, // Operating hours
                _ => 0.5, // Maintenance load
            },
            _ => 1.0,
        };

        // Base consumption + time pattern + noise
        let noise = 1.0 + rng.gen_range(-self.config.noise_factor..self.config.noise_factor);
        let mut consumption = self.config.base_load_kwh * time_multiplier * noise;

        // Inject anomaly
        let is_anomaly = rng.gen::<f64>() < self.anomaly_chance;
        if is_anomaly {
            consumption *= rng.gen_range(5.0..10.0); // Spike
        }

        self.cumulative_reading += consumption;

        // Simulate temperature (Indian climate: 20-45°C)
        let temperature = 25.0 + rng.gen_range(-5.0..20.0);

        // Simulate voltage (standard 230V ± 10%)
        let voltage = 230.0 + rng.gen_range(-23.0..23.0);

        MeterReading {
            meter_id: self.config.id.clone(),
            reading: (self.cumulative_reading * 100.0).round() / 100.0,
            timestamp: now.timestamp(),
            temperature: (temperature * 10.0).round() / 10.0,
            voltage: (voltage * 10.0).round() / 10.0,
            is_anomaly,
        }
    }
}
