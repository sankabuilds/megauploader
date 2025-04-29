# MEGA Uploader CLI

A simple command-line tool to upload files and directories to MEGA cloud storage.

## Features

- Upload single files or entire directories to MEGA
- Recursive directory upload support
- Simple command-line interface
- Secure authentication

## Prerequisites

- Go 1.16 or higher
- A MEGA account

## Installation

1. Clone this repository:
```bash
git clone https://github.com/yourusername/megauploader.git
cd megauploader
```

2. Build the application:
```bash
go build
```

3. Install the application (optional):
```bash
go install
```
This will install the binary as `megauploader` to your `$GOPATH/bin` directory, making it available system-wide.

## Usage

The tool requires three mandatory parameters:
- `-email`: Your MEGA account email
- `-password`: Your MEGA account password
- `-path`: Path to the file or directory you want to upload

### Examples

Upload a single file:
```bash
megauploader -email your@email.com -password yourpassword -path /path/to/file.txt
```

Upload an entire directory:
```bash
megauploader -email your@email.com -password yourpassword -path /path/to/directory
```

Note: If you haven't installed the tool using `go install`, you'll need to use `./megauploader` instead of just `megauploader`.

## Security Note

For security reasons, it's recommended to:
- Use environment variables for sensitive information
- Never share your MEGA credentials
- Consider using a dedicated MEGA account for automated uploads

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 