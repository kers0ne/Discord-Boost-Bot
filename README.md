# Discord Boost Bot

A Discord bot that tracks and manages server boosts with notifications and statistics.

## Features

- Real-time boost notifications
- Boost statistics and leaderboards
- Server boost level tracking
- Command system for boost info

## Commands

- `!help` - Show all commands
- `!boosts` - Display server boost status
- `!boostcount` - Show total boost count
- `!topboosters` - List top boosters
- `!ping` - Check bot latency

## Setup

1. Create a Discord application at https://discord.com/developers/applications
2. Copy your bot token to `.env` file
3. Invite bot with Administrator permissions
4. Run: `go mod tidy && go run main.go`

## Environment Variables

```
DISCORD_TOKEN=your_bot_token_here
COMMAND_PREFIX=!
BOOST_CHANNEL_ID=channel_id_for_notifications
```

## Requirements

- Go programming language (version 1.16 or later)
- Discord bot token with the necessary permissions to join servers and boost them
- Invite links of the servers you want to join
- List of Discord bot tokens (stored in a `tokens.txt` file)
- Operating system: Windows, macOS, or Linux

## Installation

1. Clone the repository to your local machine:

   ```bash
   git clone https://github.com/Spinayy/Discord-Boost-Bot.git
   ```

2. Navigate to the project directory:

   ```bash
   cd discord-server-booster
   ```

3. Compile the Go code:

   ```bash
   go build
   ```

4. Run the executable:

   ```bash
   ./discord-server-booster
   ```

## Usage

1. Choose the desired operation from the menu:

   - **Join Server**: Enter the invite code of the server you want to join.
   - **Boost Server**: Enter the ID of the server you want to boost.

2. Follow the prompts to provide the required information (invite code or server ID).

3. Sit back and let the tool automatically join servers or boost them using the provided bot tokens.

## Configuration

Before running the tool, make sure to:

- Populate the `tokens.txt` file with valid Discord bot tokens, each on a new line.
- Ensure that your Discord bot has the necessary permissions to join servers and boost them.

## Disclaimer

This tool is for educational purposes only. The developers are not responsible for any misuse of this tool.

## Contributing

Contributions are welcome! If you have any suggestions, bug fixes, or new features to add, please open an issue or create a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
