Cushon Retail ISA Investment Solution

**Overview**

This solution enables Cushon to offer ISA investments to retail customers who are not associated with an employer, allowing them to select a single fund from available options and specify their investment amount. This is with the assumption that Cushon has an existing customer database, payment processing is handled separately, and fund information is provided by another system.

**Architecture**

A three-tier architecture is implemented:

Frontend: Web-based user interface (HTML, TailwindCSS, JavaScript)

**Backend API**: Developed in Golang 

**Database:** PostgreSQL for structured data storage

The retail ISA functionality is separate from the employer-based offering, utilising distinct API endpoints and database tables.

**Features and Design Approach**

Fund Selection

Display a list of available funds.

Allow customers to select a single fund.

Fund descriptions and past performance overview.

**Investment Amount**

Users can enter the amount they want to invest.

Investment limits.

Transaction Confirmation

Securely store investment details.

Display confirmation to the user after a successful transaction.

**Security and Compliance**

Implements data encryption and secure API handling.

**Database Schema**

Customers Table – Stores customer information.

Funds Table – Maintains details of available funds.

Investments Table – Records customer investments.

**Implementation Approach**

Frontend:

Developed in HTML, TailwindCSS, and JavaScript.

Users can select funds and submit investments through a simple form.

Backend:

Golang (Gin framework) manages business logic and API endpoints.

Asynchronous processing using an in-memory queue (Kafka-ready for production).

Database:

PostgreSQL used for structured data management.

Automated database migrations are implemented for schema updates.

Queue System (Asynchronous Processing)

In-memory queue for development.

Kafka-ready configuration for production scalability.

Workers process investment transactions asynchronously to avoid API delays.

**Technical Considerations**

Authentication & Security

API rate limiting to prevent brute force attacks.

Secure request validation using middleware.

S**calability & Deployment**

Local development setup with potential AWS deployment.

Using Docker for containerisation and streamlined deployment.

AWS EC2 or ECS Fargate with auto-scaling for production.

CI/CD pipeline for automated testing and deployments.

Error Handling & Logging

Graceful API error handling with meaningful status codes.

Implements request validation using middleware.

Logs transaction processing and errors.

**Future Enhancements**

Multi-Fund Selection: Extend the system to allow investment across multiple funds.

Real-Time Investment Tracking: Provide users with insights on their investment performance.

Mobile Application Support: Develop native mobile apps for iOS and Android.

Automated Notifications: Implement email/SMS notifications for transaction updates.

Production Queue System: Replace in-memory queue with AWS SQS or Kafka.

Summary

This solution is designed for scalability & security. The event-driven approach ensures seamless investment processing, and the queue system prepares it for real-world scaling. The current implementation is functional and can be expanded with advanced features in future iterations.

