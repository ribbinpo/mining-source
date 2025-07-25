# Web Scraping URL Finder

This tool uses GoColly to scrape websites and find all URL redirect paths in HTML and JavaScript.

## Features

- Finds URLs in `<a>` tags with `href` attributes
- Extracts URLs from JavaScript `onclick` handlers
- Parses URLs in `<script>` tags
- Finds form action URLs
- Extracts meta refresh URLs
- Discovers iframe, img, and link URLs
- Normalizes relative paths to absolute URLs
- Removes duplicates automatically
- **Filters by domain** - Option to include only URLs from the same domain
- **Subdomain control** - Option to include or exclude subdomains
- **Excludes images** - Option to filter out image files (jpg, png, gif, etc.)
- **Excludes resources** - Option to filter out CSS, JS, and other resource files

## Usage

```bash
# Basic usage with default URL (example.com)
go run experiment/main.go

# Scrape a specific URL
go run experiment/main.go -url https://example.com

# Filter to same domain only (including subdomains)
go run experiment/main.go -url https://example.com -same-domain=true -include-subdomains=true

# Filter to same domain only (excluding subdomains)
go run experiment/main.go -url https://example.com -same-domain=true -include-subdomains=false

# Include images in results
go run experiment/main.go -url https://golang.org -exclude-images=false

# Include resource files (CSS, JS, etc.)
go run experiment/main.go -url https://golang.org -exclude-resources=false

# Combine filters for clean page URLs only
go run experiment/main.go -url https://golang.org -same-domain=true -exclude-resources=true
```

## Command Line Options

- `-url`: Target URL to scrape (default: https://example.com)
- `-same-domain`: Filter to only include URLs from the same domain (default: false)
- `-include-subdomains`: Include URLs from subdomains (default: true)
- `-exclude-images`: Exclude image files from results (default: true)
- `-exclude-resources`: Exclude CSS, JS, and other resource files (default: true)

## Output Format

The tool outputs a simple list of URLs and paths, one per line:

```
https://example.com/page1
https://example.com/page2
path/to/resource
https://external-site.com
```

## Supported URL Sources

1. **HTML Links**: `<a href="...">`
2. **JavaScript Redirects**:
   - `window.location.href = "url"`
   - `window.location = "url"`
   - `window.open("url")`
   - `document.location.href = "url"`
   - `location.href = "url"`
   - `location.replace("url")`
3. **Form Actions**: `<form action="...">`
4. **Meta Refresh**: `<meta http-equiv="refresh" content="5;url=...">`
5. **Iframe Sources**: `<iframe src="...">`
6. **Image Sources**: `<img src="...">` (when not excluded)
7. **Link Elements**: `<link href="...">` (when not excluded)

## Filtering Features

### Domain Filtering

When `-same-domain=true` is used, only URLs from the same domain as the target URL are included. This is useful for:

- Finding internal navigation links
- Discovering site structure
- Avoiding external dependencies

### Subdomain Control

The `-include-subdomains` flag controls whether subdomains are included:

- **`-include-subdomains=true`** (default): Includes URLs from subdomains
  - Example: Scraping `example.com` will include `miniapp.example.com`
- **`-include-subdomains=false`**: Excludes URLs from subdomains
  - Example: Scraping `example.com` will exclude `miniapp.example.com`

**Smart Domain Detection**: The tool intelligently handles:

- Regular domains: `example.com` → `example.com`
- Country TLDs: `example.co.uk` → `example.co.uk`
- Subdomains: `blog.example.com` → `example.com` (base domain)

### Image Filtering

When `-exclude-images=true` (default), image files are filtered out based on file extensions:

- `.jpg`, `.jpeg`, `.png`, `.gif`, `.bmp`, `.svg`, `.webp`, `.ico`, `.tiff`, `.tif`

This helps focus on navigational and functional URLs rather than media assets.

### Resource Filtering

When `-exclude-resources=true` (default), resource files are filtered out based on file extensions:

- **Web Resources**: `.css`, `.js`, `.json`, `.xml`, `.map`, `.min`
- **Documents**: `.pdf`, `.doc`, `.docx`, `.xls`, `.xlsx`, `.ppt`, `.pptx`, `.txt`
- **Archives**: `.zip`, `.rar`, `.7z`, `.tar`, `.gz`
- **Media**: `.mp3`, `.mp4`, `.avi`, `.mov`, `.wmv`, `.flv`, `.webm`
- **Fonts**: `.woff`, `.woff2`, `.ttf`, `.eot`, `.otf`

This ensures you get only actual page URLs, not resource files.

## Examples

### Include Subdomains

```bash
go run experiment/main.go -url https://example.com -same-domain=true -include-subdomains=true
```

Output: Includes `example.com` and `miniapp.example.com`

### Exclude Subdomains

```bash
go run experiment/main.go -url https://example.com -same-domain=true -include-subdomains=false
```

Output: Only `example.com` URLs, excludes `miniapp.example.com`

### Page URLs Only (Recommended)

```bash
go run experiment/main.go -url https://golang.org -same-domain=true -exclude-resources=true
```

Output: Only internal page URLs, no CSS/JS/images

### Include All Resources

```bash
go run experiment/main.go -url https://golang.org -exclude-resources=false -exclude-images=false
```

Output: All URLs including CSS, JS, images, and other resources

### Navigation Links Only

```bash
go run experiment/main.go -url https://golang.org -same-domain=true -exclude-resources=true -exclude-images=true
```

Output: Clean list of internal navigation pages only

## Use Cases

- **Site Mapping**: Use `-same-domain=true -exclude-resources=true` to map internal page structure
- **Subdomain Discovery**: Use `-include-subdomains=true` to find all related subdomains
- **Main Site Only**: Use `-include-subdomains=false` to focus on the main domain
- **Navigation Discovery**: Find all internal links and pages
- **SEO Analysis**: Discover all crawlable pages on a site
- **Content Discovery**: Find all content pages without resource noise

## Dependencies

- `github.com/gocolly/colly/v2` - Web scraping framework
- Standard Go libraries for URL parsing and regex

## Installation

The dependencies are automatically managed by Go modules. Run:

```bash
go mod tidy
```

to ensure all dependencies are properly installed.
