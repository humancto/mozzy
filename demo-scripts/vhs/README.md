# VHS Demo Recordings

This directory contains VHS tape files for generating animated GIF demos of mozzy.

## What is VHS?

VHS is a tool for generating terminal GIFs from script files. It's like writing a screenplay for your terminal - you define what to type, when to press Enter, and how long to wait.

## Installation

```bash
brew install vhs
```

VHS requires `ffmpeg` and `ttyd`:
```bash
brew install ffmpeg ttyd
```

## Available Tapes

1. **01-github-api.tape** - GitHub API with JQ filtering
2. **02-collections.tape** - Saving and executing collections
3. **03-jq-filtering.tape** - Multiple JQ filtering examples
4. **04-verbose-mode.tape** - Verbose output with headers and timing
5. **05-post-request.tape** - POST request with JSON body

## Generating GIFs

### Generate a single demo:

```bash
cd demo-scripts/vhs
vhs 01-github-api.tape
```

This will create: `../assets/demo-01-github-api.gif`

### Generate all demos at once:

```bash
cd demo-scripts/vhs
vhs 01-github-api.tape &
vhs 02-collections.tape &
vhs 03-jq-filtering.tape &
vhs 04-verbose-mode.tape &
vhs 05-post-request.tape &
wait
echo "✅ All GIFs generated!"
```

Or use this helper script:

```bash
#!/bin/bash
cd demo-scripts/vhs
for tape in *.tape; do
  echo "🎬 Recording $tape..."
  vhs "$tape"
done
echo "✅ All demos generated in assets/ directory"
```

## Output Location

All GIFs are generated in the `assets/` directory:
- `assets/demo-01-github-api.gif`
- `assets/demo-02-collections.gif`
- `assets/demo-03-jq-filtering.gif`
- `assets/demo-04-verbose.gif`
- `assets/demo-05-post.gif`

## Customizing Tapes

Edit the `.tape` files to customize:

### Terminal Appearance
```tape
Set FontSize 16              # Font size (default: 16)
Set Width 1200               # Terminal width in pixels
Set Height 600               # Terminal height in pixels
Set Padding 20               # Padding around terminal
Set Theme "Monokai"          # Color theme
```

Available themes: `Monokai`, `Dracula`, `Nord`, `One Dark`, `Solarized Dark`, etc.

### Timing
```tape
Sleep 1s                     # Wait 1 second
Sleep 500ms                  # Wait 500 milliseconds
Type@500ms "text"            # Type with 500ms between chars
```

### Typing Commands
```tape
Type "mozzy GET /api"        # Type command (don't execute)
Enter                        # Press Enter key
Type "command" Enter         # Type and execute
```

### Output Format
```tape
Output demo.gif              # Generate GIF (default)
Output demo.mp4              # Generate MP4 video
Output demo.webm             # Generate WebM video
```

## Tips for Better Recordings

1. **Keep it short** - 10-30 seconds is ideal
2. **Add context** - Start with a comment explaining what you're doing
3. **Show results** - Wait long enough to see the full response
4. **Use colors** - Always add `--color` flag to mozzy commands
5. **Clean output** - Make sure mozzy outputs are readable

## Troubleshooting

### "mozzy: command not found"
Make sure mozzy is installed and in your PATH:
```bash
which mozzy
```

### Recording takes too long
VHS records in real-time. Reduce sleep times in the tape file.

### GIF file size too large
1. Reduce terminal width/height
2. Shorten the recording
3. Use GIF optimization tools:
```bash
brew install gifsicle
gifsicle -O3 --lossy=80 input.gif -o output.gif
```

### Colors not showing
Make sure your tape uses:
```tape
Set Theme "Monokai"          # Or another theme
```
And mozzy commands use `--color` flag.

## Advanced: Custom Theme

Create a custom theme JSON:

```json
{
  "name": "Custom",
  "black": "#1e1e1e",
  "red": "#f44747",
  "green": "#4ec9b0",
  "yellow": "#ffcc66",
  "blue": "#3b8eea",
  "purple": "#c586c0",
  "cyan": "#4fc1ff",
  "white": "#d4d4d4"
}
```

Use it:
```tape
Set Theme path/to/theme.json
```

## Resources

- VHS Documentation: https://github.com/charmbracelet/vhs
- Available Themes: https://github.com/charmbracelet/vhs#themes
- Examples: https://github.com/charmbracelet/vhs/tree/main/examples
