# WifiAnnouncer

WifiAnnouncer scans your local WiFi network and monitors known hostnames for their presence. It will announce when they connect to the network, and when they leave. This is (usually) indicative of someone leaving or arriving at a location.

## Requirements

1. Golang must be installed
2. You must create a GCP account & project, and create a service account file that has permission of

- Go to GCP dashboard `->` IAM & Admin `->` Service Accounts on the left nav `->` Create Service Account on the top nav `->` Input a name `->` Create and Continue `->` Select a role `->` Cloud Speech Editor `->` Continue `->` Done
- Then click on your service account in the list `->` Keys in the top nav `->` Add Key `->` Create new key `->` JSON `->` Create
- That newly downloaded file will be what you set the `GOOGLE_APPLICATION_CREDENTIALS` variable to

## Building and Running Locally

1. Clone the repo: `git clone https://github.com/danthegoodman1/wifiAnnouncer`
2. Build it: `go build`
3. Set your GCP service account file path environment variable: `export GOOGLE_APPLICATION_CREDENTIALS=/some/path/project-210111-910eb110cabd.json`
4. Run it: `./wifiAnnouncer`

## Building and Running with balena

You don't want to run this on your laptop, so let's put it on something we can stick in the corner.

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

### Debug logging

Executing

### DB Corruption

If you experience errors regarding `people.db`, you can reset the potentially corrupted SQLite3 file by deleting it. WifiAnnouncer will recreate it automatically, but will lose its state (whether people are connected or away)
