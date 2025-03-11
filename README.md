# Cushion ISA Investment Platform

A lightweight, event-driven investment system with a basic web interface.

## Features

- Basic UI with responsive design
- Real-time investment submission
- Asynchronous transaction processing
- Input validation and error handling
- Queue-based architecture (easily extensible to AWS SQS)

## Prerequisites

- Go 1.16 or later
- Git

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd cushion-isa
```

2. Install dependencies:
```bash
go mod tidy
```

## Running the Application

1. Start the server:
```bash
go run main.go
```

2. Open your browser and navigate to:
```
http://localhost:8080
```

## Architecture

- **Frontend**: HTML5, TailwindCSS, and vanilla JavaScript
- **Backend**: Go with Gin framework
- **Queue System**: In-memory Go channel (easily replaceable with AWS SQS)
- **Processing**: Asynchronous worker service

## API Endpoints

### POST /api/invest
Submit a new investment request

**Request Body:**
```json
{
    "name": "John Doe",
    "email": "john@example.com",
    "fund": "sustainable-growth",
    "amount": 1000.00
}
```

**Response:**
```json
{
    "message": "Investment request received and being processed",
    "status": "success"
}
```

## Future Enhancements

1. Integration with AWS SQS for production-grade message queue
2. PostgreSQL database integration for transaction storage
3. User authentication and authorization
4. Transaction history view
5. Email notifications for investment status updates

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 
>>>>>>> 7719b45 (Initial commit)
