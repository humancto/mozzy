# Mozzy Demo Scripts

These scripts demonstrate mozzy's key features and can be used to create GIF demos for the README.

## Creating GIFs/Screenshots

### Option 1: macOS Built-in Screen Recording (Easiest)

1. Open Terminal and navigate to this directory
2. Press `Cmd + Shift + 5` to open Screenshot toolbar
3. Select "Record Selected Portion"
4. Select your terminal window
5. Run the script: `./01-github-api.sh`
6. Press `Cmd + Ctrl + Esc` to stop recording
7. The video will be saved to your Desktop

### Option 2: Using asciinema (Terminal Recordings)

```bash
# Install asciinema
brew install asciinema

# Record a demo
asciinema rec demo1.cast -c "./01-github-api.sh"

# Upload to asciinema.org (optional)
asciinema upload demo1.cast

# Convert to GIF (requires agg)
brew install agg
agg demo1.cast demo1.gif
```

### Option 3: Using VHS (Automated GIF Generation)

```bash
# Install vhs
brew install vhs

# Create a tape file (example: demo.tape)
# Then generate GIF
vhs demo.tape
```

## Available Demo Scripts

1. **01-github-api.sh** - GitHub API exploration with colored output and jq filtering
2. **02-collections.sh** - Saving and executing request collections
3. **03-verbose-mode.sh** - Verbose output with headers and timing breakdown
4. **04-jq-filtering.sh** - Multiple examples of inline JQ filtering
5. **05-post-request.sh** - POST request with JSON body

## Running the Scripts

Make them executable:
```bash
chmod +x *.sh
```

Run individually:
```bash
./01-github-api.sh
./02-collections.sh
./03-verbose-mode.sh
./04-jq-filtering.sh
./05-post-request.sh
```

## Tips for Good Recordings

1. **Terminal Size**: Set to 80x24 or 100x30 for better visibility
2. **Font Size**: Increase font size (Cmd + +) for readability
3. **Color Scheme**: Use a theme with good contrast (Solarized Dark, Tomorrow Night, etc.)
4. **Timing**: The scripts have built-in `sleep` delays - adjust as needed
5. **Clean**: Clear terminal before recording (`clear` command)

## Converting to GIF

After recording, you may need to:
1. Trim the video (use QuickTime or ffmpeg)
2. Convert to GIF (use ffmpeg, gifski, or online converters)
3. Optimize GIF size (use gifsicle or online optimizers)

### Using ffmpeg to convert:

```bash
# Install ffmpeg
brew install ffmpeg

# Convert video to GIF
ffmpeg -i recording.mov -vf "fps=10,scale=800:-1:flags=lanczos" -c:v gif output.gif

# Optimize GIF size
brew install gifsicle
gifsicle -O3 --lossy=80 output.gif -o optimized.gif
```

## Hosting the GIFs

Upload the GIFs to:
- GitHub repo in `assets/` or `demos/` folder
- GitHub Releases
- CDN like Cloudinary, Imgur, or imgbb

Then update the README.md image links from placeholders to actual GIF URLs.
