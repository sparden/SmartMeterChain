mod config;
mod generator;
mod client;

use clap::Parser;
use std::time::Duration;
use tracing::{info, error};

#[derive(Parser)]
#[command(name = "smartmeter-simulator")]
#[command(about = "Simulates smart meter readings for SmartMeterChain")]
struct Cli {
    /// Path to config file
    #[arg(short, long, default_value = "config.yaml")]
    config: String,

    /// Override interval (seconds)
    #[arg(short, long)]
    interval: Option<u64>,

    /// Register meters before starting simulation
    #[arg(long)]
    register: bool,

    /// Run only once (single reading per meter)
    #[arg(long)]
    once: bool,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    tracing_subscriber::fmt()
        .with_env_filter(
            tracing_subscriber::EnvFilter::try_from_default_env()
                .unwrap_or_else(|_| "info".into()),
        )
        .init();

    let cli = Cli::parse();
    let cfg = config::SimConfig::load(&cli.config)?;
    let interval = Duration::from_secs(cli.interval.unwrap_or(cfg.simulation.interval_seconds));

    info!("SmartMeterChain Simulator starting");
    info!("API: {}", cfg.api.base_url);
    info!("Meters: {}", cfg.meters.len());
    info!("Interval: {}s", interval.as_secs());

    let api_client = client::ApiClient::new(&cfg.api)?;

    // Authenticate
    info!("Authenticating...");
    api_client.login().await?;
    info!("Authenticated successfully");

    // Register meters if requested
    if cli.register {
        for meter in &cfg.meters {
            match api_client.register_meter(meter).await {
                Ok(_) => info!("Registered meter: {}", meter.id),
                Err(e) => error!("Failed to register {}: {}", meter.id, e),
            }
        }
    }

    // Initialize generators
    let mut generators: Vec<generator::MeterGenerator> = cfg
        .meters
        .iter()
        .map(|m| generator::MeterGenerator::new(m.clone(), cfg.simulation.inject_anomaly_chance))
        .collect();

    info!("Starting simulation loop...");

    loop {
        for gen in &mut generators {
            let reading = gen.next_reading();
            info!(
                "[{}] Reading: {:.2} kWh (anomaly: {})",
                reading.meter_id, reading.reading, reading.is_anomaly
            );

            match api_client.submit_reading(&reading).await {
                Ok(tx_id) => info!("[{}] Submitted -> tx: {}", reading.meter_id, tx_id),
                Err(e) => error!("[{}] Submit failed: {}", reading.meter_id, e),
            }
        }

        if cli.once {
            info!("Single run complete. Exiting.");
            break;
        }

        tokio::time::sleep(interval).await;
    }

    Ok(())
}
