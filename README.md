# opencode-localproviders

Auto-populate opencode.ai provider configurations from live OpenAI-compatible APIs.

## Usage

```bash
./opencode-localproviders --base-url http://localhost:1234/ --provider lmstudio
```

### Flags

- `--base-url` / `-u`: OpenAI-compatible API base URL (e.g., `http://localhost:1234/v1` or `http://localhost:1234/`)
- `--provider` / `-p`: Provider key name (e.g., `lmstudio`, `ollama`)
- `--name` / `-n`: Display name (default: titlecase provider + " (local)")
- `--config` / `-c`: Path to `opencode.json` (default: `~/.config/opencode/opencode.json`)
- `--dry-run`: Print generated config block without modifying file

## Examples

### LM Studio (local)
```bash
./opencode-localproviders --base-url http://localhost:1234/ --provider lmstudio
```

### Ollama
```bash
./opencode-localproviders --base-url http://localhost:11434/ --provider ollama
```

### Dry run (preview changes)
```bash
./opencode-localproviders --base-url http://localhost:1234/ --provider lmstudio --dry-run
```

## How it works

1. Normalizes base URL to OpenAI-compatible format (`/v1`)
2. Fetches available models from `GET {baseURL}/models`
3. Generates provider config block with all models
4. Reads `~/.config/opencode/opencode.json`
5. Updates or creates `provider.<name>` entry
6. Writes config back with indentation preserved

## Config format

Generated config:
```json
{
  "npm": "@ai-sdk/openai-compatible",
  "name": "LMStudio (local)",
  "options": {
    "baseURL": "http://localhost:1234/v1"
  },
  "models": {
    "model-id": {
      "name": "model-id"
    }
  }
}
```

## Build

```bash
go build .
```

No external dependencies — uses only Go stdlib.
