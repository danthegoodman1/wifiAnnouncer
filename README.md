# WifiAnnouncer

WifiAnnouncer scans your local WiFi network and monitors known hostnames for their presence. It will announce when they connect to the network, and when they leave. This is (usually) indicative of someone leaving or arriving at a location.

## Configuration

The following is an example config file:

```yaml
# The name of the GCP voice to use
voiceName: en-US-Wavenet-D

# The country code of the voice to use
languageCode: en-US

# Which interface to scan (any IP in the /24), in this example, we will scan 192.168.86.0/24
interface: "192.168.86.1"

# What to suffix the `name` with when someone arrives
arrivedSuffix: has arrived

# What to suffix the `name` with when someone leaves
leftSuffix: has left

# What to prefix the `name` with when someone arrives
arrivedPrefix: Look!

# What to prefix the `name` with when someone leaves
leftPrefix: Goodbye,

# The devices to speak about
registeredDevices:
  # The name to speak
  - name: Dan
    # Their network identifier
    hostname: dans-iphone-x.lan.
    # Whether when they are added, they are known as here or away (initial state)
    defaultState: away
```

With this configuration, when my iPhone X disconnects from the wifi, the voice will say _Look! Dan has arrived_, and when it disconnects it will say _Goodbye, Dan has left_.

To find all voice and language code options visit: https://cloud.google.com/text-to-speech#section-2 and play with the options!

## How it Works

### Scanning for Devices

### Registering Devices

### DB Corruption

If you experience errors regarding `people.db`, you can reset the potentially corrupted SQLite3 file by deleting it. WifiAnnouncer will recreate it automatically, but will lose its state (whether people are connected or away)
