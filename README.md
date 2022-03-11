# API Server

The API and rendering service responsible for serving images for mineatar.io.

## Getting Started

```bash

# Clone the repository into a `api-server` folder
git clone https://github.com/mineatar-io/api-server.git

# Navigate the working directory into the api-server folder
cd api-server

# Copy the `config.example.yml` file a new file `config.yml`
# You will need to edit the details of this file for proper functionality
cp config.example.yml config.yml

# Install the Go dependencies
go get ...

# Build the source code into a single executable
./scripts/build

# Run the executable
./bin/main
```

## Issues

If you find any issues with this API service (not the website itself), please create a [new issue](https://github.com/mineatar-io/api-server/issues) with all necessary details.

## License

[MIT License](https://github.com/mineatar-io/api-server/blob/master/LICENSE)