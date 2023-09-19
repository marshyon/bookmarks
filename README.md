
# Go Pinboard Bookmarks Project

The Go Pinboard Bookmarks Project is a simple Go application that reads https://pinboard.in/export/format:json/ JSON data, processes it, and upserts the data into an SQLite database using GORM. This project works therefore with JSON, time parsing, GORM for database interactions, and upsert operations.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)

## Installation

1. Clone the repository to your local machine:

   ```bash
   git clone https://github.com/marshyon/go-bookmarks.git
   ```

2. Change to the project directory:

   ```bash
   cd go-bookmarks
   ```

3. Install the required dependencies (GORM):

   ```bash
   go get -u gorm.io/gorm
   ```

4. Create the SQLite database:

   ```bash
   touch bookmarks.db
   ```

5. Build the project:

   ```bash
   go build
   ```

## Usage

1. Prepare your JSON data in a file named `bookmarks.json` in the project directory. The JSON data should have the same structure as the provided sample data.

2. Run the application:

   ```bash
   ./go-bookmarks
   ```

   This will read the JSON data, process it, and upsert it into the SQLite database.

3. You can then interact with the SQLite database using your preferred SQLite client or utility.

## Project Structure

The project structure is organized as follows:

- `main.go`: The main application logic for reading JSON data, processing it, and upserting it into the database.
- `bookmarks.json`: Sample JSON data with bookmarks.
- `bookmarks.db`: The SQLite database where the converted bookmarks are stored.
- `README.md`: This README file.

## Contributing

Contributions to this project are welcome! If you'd like to contribute, please follow these guidelines:

1. Fork the repository.
2. Create a new branch for your feature or bug fix: `git checkout -b feature-name`.
3. Make your changes and commit them: `git commit -m 'Add new feature'`.
4. Push your branch to your forked repository: `git push origin feature-name`.
5. Create a pull request from your forked repository to this repository.

Please make sure your code follows the existing coding style and includes appropriate tests. We'll review your pull request as soon as possible.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
