// I've implemented a dependency injection pattern for the job handlers. This makes the handlers clean, testable, and extensible.
//
// Here's what I've done:
// 1.  **`internal/job/context.go`**: Created a `JobHandlerContext` struct to hold dependencies for job handlers.
// 2.  **`internal/job/service.go`**: Modified the `JobService` to hold the `JobHandlerContext` and created a `RegisterHandler` method to register handlers that are aware of this context.
// 3.  **`internal/job/handlers/verification_email.go`**: Created a `HandleVerificationEmail` job handler that uses the `EmailService` from the `JobHandlerContext`.
//
// To complete the setup, you'll need to:
// 1.  Instantiate the `EmailService` in your main application setup.
// 2.  Create a `JobHandlerContext` and provide the `EmailService` to it.
// 3.  Pass the `JobHandlerContext` when creating the `JobService`.
// 4.  Register the `HandleVerificationEmail` handler with the `JobService`.
//
// Example of how to register the handler:
// ```go
// jobService.RegisterHandler(handlers.TypeVerificationEmail, handlers.HandleVerificationEmail)
// ```
//
// This setup ensures that your job handlers are decoupled from the rest of your application and can be easily tested and maintained.
package job

import "go-boilerplate/internal/mailer"

// JobHandlerContext provides access to services required by job handlers.
// It acts as a dependency injection container for the jobs.
type JobHandlerContext struct {
	EmailService *mailer.EmailService
}
