# WifiAnnouncer

WifiAnnouncer scans your local WiFi network and monitors known hostnames for their presence. It will announce when they connect to the network, and when they leave. This is (usually) indicative of someone leaving or arriving at a location.

## Configuration

## How it Works

### Scanning for Devices

### Registering Devices

### DB Corruption

If you experience errors regarding `people.db`, you can reset the potentially corrupted SQLite3 file by deleting it. WifiAnnouncer will recreate it automatically, but will lose its state (whether people are connected or away)
