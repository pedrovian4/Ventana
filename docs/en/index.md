[English](../en/index.md) | [Português](../pt/index.md)

# Ventana Documentation

## Configurations

Ventana uses a JSON configuration file to store user preferences, including language, connection history path, message directory, and other settings. Below are the main configurations and how to use them.

## Configuration File Structure
Ventana creates a default configuration file on its first run.

Example of a `ventana.json` configuration file:

```json
{
  "language": "en",
  "history_file": "/home/user/.config/ventana/history.csv",
  "message_directory": "./messages",
  "default_retry_count": 3,
  "default_retry_delay": 5,
  "config_file_path": "/home/user/.config/ventana/ventana.json",
  "scripts_directory": "/home/user/.config/ventana/scripts"
}
```

## Configuration Explanation

- language: Defines the interface language for Ventana. Supported languages are "en" (English) and "pt" (Portuguese).

- history_file: Path to the history file (history.csv). This file stores past connections, allowing quick access to previously used servers. The default path is ~/.config/ventana/history.csv.

- message_directory: Directory where interface messages are stored. Ventana uses message files to display text in the language set by the "language" configuration. The default directory is ./messages.

- default_retry_count: Defines the default number of retry attempts in case of a failed connection to a server. The default value is 3, but it can be adjusted as needed.


- default_retry_delay: Time (in seconds) to wait between each reconnection attempt. The default value is 5 seconds.


- config_file_path: Path to the configuration file itself (ventana.json). This file stores preferences and should be located at ~/.config/ventana/ventana.json.

- scripts_directory: Directory where scripts should be stored for execution. Ventana runs scripts written in Shell using bash, and these scripts should be placed in ~/.config/ventana/scripts to be recognized by the program.