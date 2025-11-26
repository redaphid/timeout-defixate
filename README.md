# timeout-defixate

Break enforcer for the [Time Out](https://www.dejal.com/timeout/) app. Helps manage ADHD hyperfocus by escalating consequences when you skip or postpone too many breaks.

## How It Works

1. Monitors macOS unified logs for Time Out's explicit log messages
2. Detects `"Postponed the break"` and `"Skipped the break"` events
3. After 5 skips/postpones → **lock loop** (re-locks screen every 5s for 2 minutes)
4. After 10 skips/postpones → **shuts down your computer**
5. Counter resets when you complete a break (`"finished: yes"`)

## Requirements

- **macOS** (uses unified logging and private screen lock API)
- **Time Out app** with diagnostic logging enabled

### Enabling Time Out Console Logging

1. Open Time Out preferences
2. Go to **Advanced** in the sidebar
3. Scroll down to **Diagnostic options**
4. Check **Output scheduler logging**
5. Check **Only include significant changes**

## Build

```bash
go build -o timeout-defixate .
```

## Usage

```bash
# Default: 5 skips = lock, 10 = shutdown
./timeout-defixate

# Custom limits
./timeout-defixate --lock-limit=3 --shutdown-limit=5
```

## Installation as LaunchAgent

Copy this plist to run at login:

```bash
cat > ~/Library/LaunchAgents/com.hypnodroid.timeout-defixate.plist << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.hypnodroid.timeout-defixate</string>
    <key>ProgramArguments</key>
    <array>
        <string>/Users/hypnodroid/Projects/timeout-defixate/timeout-defixate</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/tmp/timeout-defixate.log</string>
    <key>StandardErrorPath</key>
    <string>/tmp/timeout-defixate.err</string>
</dict>
</plist>
EOF
```

Load the agent:

```bash
launchctl load ~/Library/LaunchAgents/com.hypnodroid.timeout-defixate.plist
```

### View logs

```bash
tail -f /tmp/timeout-defixate.log
```

### Stop the service

```bash
launchctl unload ~/Library/LaunchAgents/com.hypnodroid.timeout-defixate.plist
```
