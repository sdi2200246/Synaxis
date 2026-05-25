# Synaxis

> A full-stack event management and booking platform built with Go, React, and PostgreSQL.

[![Demo Video](https://img.shields.io/badge/Demo-Watch%20Video-red?style=for-the-badge&logo=youtube)](<!-- DEMO_LINK -->)
![Go](https://img.shields.io/badge/Go-1.22-00ADD8?style=for-the-badge&logo=go)
![React](https://img.shields.io/badge/React-18-61DAFB?style=for-the-badge&logo=react)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?style=for-the-badge&logo=postgresql)
![Python](https://img.shields.io/badge/Python-3.11-3776AB?style=for-the-badge&logo=python)

---

Synaxis is a platform where organizers can create and manage events, define ticket types with real availability tracking, and communicate directly with attendees — while attendees can discover events, make bookings, and receive automatic notifications when something changes. Guest users can browse and explore the full event catalog without creating an account, with contextual prompts guiding them to register only when they want to take action.
 
---

## Engineering Goals

This project was built with explicit goals beyond just making things work:

**Clean Architecture**
> The backend is strictly layered: entities → repositories → services → handlers. Each layer depends only inward. Service interfaces are defined at the boundary so implementations are swappable without touching business logic.

**Correct REST API**
> Endpoints use appropriate HTTP semantics — status codes, content negotiation (`Accept: application/xml` vs `application/json`), paginated list responses with `has_more`, and explicit action endpoints for domain operations like publishing or cancelling an event rather than raw status mutations from the client.

**DDD**
> Business rules live on the entities themselves, not scattered across handlers. Methods like `ApproveCancellation()` and `ApproveCreate()` encode invariants directly on the domain type — the service layer orchestrates, but the entity decides whether an operation is valid.

**Normalized Database Design**
> The relational schema follows normalization rules — no redundant data, clear foreign key relationships, and junction tables where needed (e.g. event↔category). Every table has a single well-defined responsibility, making queries predictable and updates safe.
 
**Data Integrity and Consistency**
> Anywhere two operations need to succeed or fail together — ticket quantity updates, booking creation under concurrency, conversation + first message on cancellation — the implementation uses transactions and atomic SQL patterns. Capacity invariants are enforced at the service layer before any write reaches the database.

**Consistent, Enterprise-grade UI**
> The frontend follows a strict design system — shared CSS variables for color, spacing, and typography, a reusable component and utility class library, and consistent patterns for cards, dialogs, forms, and states across every page. The result is a UI that feels coherent end-to-end rather than assembled from parts.


## Tech Stack
 
| Layer | Technology |
|---|---|
| **Backend** | Go 1.22, Gin |
| **Database** | PostgreSQL 16, pgx/pgxpool |
| **Migrations** | golang-migrate |
| **Query building** | tqla (dynamic SQL templates) |
| **Authentication** | JWT (golang-jwt/v5), bcrypt |
| **Frontend** | React 18, TypeScript, Vite |
| **Routing** | React Router v6 |
| **HTTP client** | Axios |
| **Maps** | react-leaflet, OpenStreetMap |
| **ML Pipeline** | Python 3.11, NumPy |
| **Deployment** | Docker(native development)|
| **CI** | GitHub Actions |



## Architecture Overview
 
The backend follows a strict layered architecture where each layer depends only inward. No layer is skipped and no concerns leak across boundaries.

 <p align="center">
     <img height="600" alt="synaxis_architecture_fixed" src="https://github.com/user-attachments/assets/ff14702c-60f9-4fc1-98b1-4341ced5cd46" />
</p>
 
**Entities and interfaces** run vertically through all backend layers. Entities carry domain validation methods (`ApproveCreate()`, `ApproveCancellation()`, `HasCapacityFor()`). Interfaces are consumer-defined and narrow — no layer takes a fat dependency on another layer's concrete type.
 
### Key architectural decisions
 
**EventBus for decoupled side effects**
> `EventService` changes event status and publishes an `EventCancelled` domain event. `CancelEventService` subscribes independently and handles attendee notification via messaging — neither service knows about the other. The bus is in-memory, topic-based, with goroutine delivery and `sync.WaitGroup` guarantees. Accepted tradeoff: events are not durable across crashes, which is acceptable at this scale.
 
**Service-of-services avoided via narrow interfaces**
> Cross-aggregate reads go directly to repos through consumer-defined interfaces (`EventCapacityProvider`, `VenueOwnershipChecker`) rather than calling another service. This prevents exploiting gaps in a called service's auth logic and keeps every service describable in a single sentence with no "and".
 
**DTOs confined to the controller layer**
> Services return bare Go structs with no JSON or XML tags. Marshalling and response shaping live exclusively in controllers. The XML/JSON export endpoint reuses the same service methods as the regular API — transport format is invisible to the service layer.

**Repository pattern with interface-driven dependencies**
> Every data access concern lives behind a repository interface defined in `interfaces/interfaces.go`. Services depend on the interface, never the concrete repo. This keeps business logic fully testable without a live database and makes the data layer swappable. Repos own their transactions internally — no `tx` objects are passed between repos or up to the service layer.

**ML pipeline as a fully isolated layer**
> The recommendation engine is a completely separate Python process with no runtime coupling to the Go backend. It reads interaction data (visits and bookings) directly from PostgreSQL, trains the Biased Matrix Factorization model, and writes scored recommendations back to the `recommendation` table. The Go backend reads from that table at query time — the two sides share only the database, never a function call or network request. The pipeline is re-run manually after new interaction data accumulates, making the recommendation layer independently deployable and updatable without touching the application server.


## REST API Design Decisions 
<details>
<summary><strong>Plural nouns for every resource</strong></summary>
<blockquote>All endpoints use plural nouns: <code>/events</code>, <code>/bookings</code>, <code>/conversations</code>, <code>/venues</code>, <code>/categories</code>, <code>/users</code>, <code>/tickets</code>. No singular forms, no verb-based routes. The HTTP method communicates the operation; the URL identifies the resource.</blockquote>
</details>
<details>
<summary><strong>Correct HTTP method semantics</strong></summary>
<blockquote><code>POST</code> creates, <code>GET</code> reads, <code>PATCH</code> partially updates, <code>DELETE</code> removes. No <code>PUT</code> anywhere — every update endpoint accepts a sparse payload and touches only the fields present. This eliminates read-modify-write cycles and accidental overwrites.</blockquote>
</details>
<details>
<summary><strong>Meaningful status codes on every response path</strong></summary>
<blockquote><code>201 Created</code> on successful creation. <code>409 Conflict</code> when a capacity constraint or uniqueness rule is violated. <code>403 Forbidden</code> for ownership failures. <code>401 Unauthorized</code> for missing or expired tokens. <code>400 Bad Request</code> for malformed input. Every known failure maps to a named <code>apperr</code> sentinel with a specific code — <code>500</code> only surfaces for genuinely unexpected errors.</blockquote>
</details>
<details>
<summary><strong>Nested routes when parent context is required, flat when the child ID is sufficient</strong></summary>
<blockquote>Child resource creation and listing nest under the parent because the parent ID is required: <code>POST /events/:id/tickets</code>, <code>GET /events/:id/bookings</code>, <code>POST /conversations/:id/messages</code>. Individual child operations use flat routes when the child ID alone is enough: <code>GET /tickets/:id</code>, <code>PATCH /messages/:id</code>, <code>GET /bookings</code>. This avoids forcing the client to supply a redundant parent ID on operations where it adds no value.</blockquote>
</details>
<details>
<summary><strong>Path parameters for identity, query parameters for filtering</strong></summary>
<blockquote>Path segments identify a specific resource: <code>/events/:id</code>, <code>/tickets/:ticket_id</code>, <code>/users/:id</code>. Query parameters filter or modify collections: <code>GET /events?title=jazz&category_id=...&limit=20</code>, <code>GET /venues?name=...&capacity=...</code>, <code>GET /admin/events?organizer_id=...</code>. The two are never mixed — no resource identifier ever appears as a query parameter, and no filter ever appears in the path.</blockquote>
</details>
<details>
<summary><strong>Server-side hydration — one request, complete response</strong></summary>
<blockquote>Event responses include nested <code>venue</code>, <code>categories[]</code>, and <code>media[]</code> objects hydrated at the service layer. The client never needs a second request to render an event card. Related data that belongs together structurally is composed server-side rather than requiring client-side joins.</blockquote>
</details>
<details>
<summary><strong>Pagination via limit/offset with <code>has_more</code></strong></summary>
<blockquote>All list endpoints return <code>{ items, has_more }</code> instead of a total count. <code>has_more</code> is computed by fetching <code>limit + 1</code> rows — no <code>COUNT(*)</code> query needed. The frontend uses <code>has_more</code> to drive infinite scroll via <code>IntersectionObserver</code>.</blockquote>
</details>
<details>
<summary><strong>Content negotiation for export</strong></summary>
<blockquote>The admin export endpoint returns XML or JSON based on the <code>Accept</code> header rather than exposing separate endpoints per format. <code>c.NegotiateFormat()</code> handles dispatch; DTOs carry both <code>json</code> and <code>xml</code> struct tags. The serializer is selected at the controller layer only — the service returns bare structs.</blockquote>
</details>
<details>
<summary><strong>OptionalAuth for dual public/authenticated behavior</strong></summary>
<blockquote><code>GET /events</code> is served through <code>OptionalAuth</code> middleware — the JWT is parsed when present but not required. Guest users see published events; authenticated users get richer responses. One route, no duplication.</blockquote>
</details>
<details>
<summary><strong>Identity from token, not from the URL</strong></summary>
<blockquote>Early versions had <code>GET /my-events</code> and <code>GET /my-bookings</code>. These were retired. <code>GET /bookings</code> reads the user ID from the JWT to return that user's bookings. The API surface stays minimal and avoids the <code>my-*</code> anti-pattern.</blockquote>
</details>
<details>
<summary><strong>Atomic booking creation</strong></summary>
<blockquote><code>POST /events/:id/bookings</code> runs inside a single database transaction that decrements <code>available</code> on the ticket type and inserts the booking row together. If either fails, the entire transaction rolls back — no partial state, no overselling under concurrent requests.</blockquote>
</details>
<details>
<summary><strong>Admin namespace for elevated operations</strong></summary>
<blockquote>Admin-only routes live under <code>/admin/</code> with dedicated middleware (<code>AdminOnly()</code>). Approve and reject are modeled as <code>POST /admin/users/:id/approve</code> and <code>POST /admin/users/:id/reject</code> — explicit action sub-resources rather than a generic <code>PATCH</code> with a status field, because these are irreversible administrative decisions, not data edits.</blockquote>
</details>
<details>
<summary><strong>Full API reference</strong></summary>
<br>
Authentication uses `Authorization: Bearer <token>`.
 
### Auth
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| `POST` | `/users` | Public | Register a new user |
| `POST` | `/auth/login` | Public | Login, returns JWT |
 
### Events
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| `GET` | `/events` | Public (OptionalAuth) | Browse published events (`title`, `category_id`, `venue_id`, `event_type`, `start_after`, `start_before`, `limit`, `offset`) |
| `GET` | `/events/:id` | Public | Get single event |
| `GET` | `/events/:id/categories` | Public | Get categories for an event |
| `GET` | `/events/recommendations` | Authenticated | Personalized recommended events (`limit`, `offset`) |
| `POST` | `/events` | Authenticated | Create a new event |
| `PATCH` | `/events/:id` | Authenticated | Update fields or transition status |
| `DELETE` | `/events/:id` | Authenticated | Delete a draft event |
 
### Ticket Types
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| `POST` | `/events/:id/tickets` | Authenticated | Create a ticket type |
| `GET` | `/events/:id/tickets` | Authenticated | List ticket types for an event |
| `PATCH` | `/events/:id/tickets/:ticket_id` | Authenticated | Update a ticket type |
| `GET` | `/tickets/:id` | Authenticated | Get a single ticket type |
 
### Bookings
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| `POST` | `/events/:id/bookings` | Authenticated | Create a booking |
| `GET` | `/events/:id/bookings` | Authenticated | List bookings for an event |
| `GET` | `/bookings` | Authenticated | List own bookings |
 
### Media
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| `POST` | `/events/:id/media` | Authenticated | Upload a photo |
| `DELETE` | `/events/:id/media/:media_id` | Authenticated | Delete a photo |
| `GET` | `/media/events/:id/:filename` | Public | Serve photo (static, bypasses handler) |
 
### Visits
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| `POST` | `/events/:id/visits` | Authenticated | Record a view interaction |
 
### Messaging
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| `POST` | `/conversations` | Authenticated | Start a conversation on a booking |
| `GET` | `/conversations` | Authenticated | List own conversations |
| `PATCH` | `/conversations/:id/read` | Authenticated | Mark conversation as read |
| `POST` | `/conversations/:id/messages` | Authenticated | Send a message |
| `GET` | `/conversations/:id/messages` | Authenticated | Get messages in a conversation |
| `PATCH` | `/messages/:id` | Authenticated | Edit or soft-delete a message |
 
### Reference Data
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| `GET` | `/venues` | Public | List venues (`name`, `capacity` filters) |
| `GET` | `/venues/:id` | Public | Get a single venue |
| `GET` | `/categories` | Public | List all categories |
| `GET` | `/users/:id` | Authenticated | Get user by ID |
 
### Admin
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| `GET` | `/admin/users` | Admin | List all users |
| `POST` | `/admin/users/:id/approve` | Admin | Approve a user |
| `POST` | `/admin/users/:id/reject` | Admin | Reject a user |
| `GET` | `/admin/events?organizer_id=` | Admin | Export organizer events (`Accept: application/json` or `application/xml`) |
 
</details>

## Database Design
 
13 tables, fully normalized (3NF), managed through `golang-migrate` migrations.

<p align="center">
  <a href="https://github.com/user-attachments/assets/1fde81b6-b82b-46af-aa47-a94c041c885d">
    <img src="https://github.com/user-attachments/assets/1fde81b6-b82b-46af-aa47-a94c041c885d" alt="Database Schema" width="100%">
  </a>
</p>

### Design decisions
 
<details>
<summary><strong>UUID primary keys on all tables</strong></summary>
<blockquote>All entities use UUID surrogate keys. Natural keys (username, email) are user-controlled and can change — UUID never changes, keeping all FK references stable across the schema. Uniqueness on natural keys is enforced via <code>UNIQUE</code> constraints instead.<br><br>Exception: junction tables (<code>eventcategory</code>, <code>conversation_participant</code>, <code>recommendation</code>) use composite PKs because they are pure relationship tables with no external references pointing at them.</blockquote>
</details>
<details>
<summary><strong>Conversation 3NF fix — removed transitive dependencies</strong></summary>
<blockquote><code>attendee_id</code> and <code>organizer_id</code> were initially columns on <code>conversation</code>. Both are transitively dependent via <code>booking_id</code>:<br><br><code>conversation → booking_id → user_id</code> (attendee)<br><code>conversation → booking_id → ticket_type_id → event_id → organizer_id</code><br><br>Fix: both columns were removed. Participants are now stored in a separate <code>conversation_participant</code> table with a composite PK <code>(conversation_id, user_id)</code> and a role column. Participant identity is derived via joins with proper indexes — three joins run in microseconds.</blockquote>
</details>
<details>
<summary><strong>total_cost — intentional denormalization</strong></summary>
<blockquote><code>Booking.total_cost</code> is technically derivable from <code>number_of_tickets × TicketType.price</code>. It is kept as a deliberate denormalization because:<br><br>1. The spec DTD explicitly requires <code>TotalCost</code> as a Booking element<br>2. It is a price snapshot — if the organizer changes the ticket price later, existing bookings must reflect what was actually paid<br>3. It is treated as immutable — never recalculated after creation<br><br>This is standard practice in all booking and e-commerce systems.</blockquote>
</details>
<details>
<summary><strong>Message.is_read instead of a MessageRead junction</strong></summary>
<blockquote>A <code>message_read</code> junction table was considered for per-user read tracking. Rejected because conversations always have exactly two participants. With two participants, a boolean <code>is_read</code> on <code>message</code> is sufficient: the sender always considers their own message read, and the receiver marks it when they open the conversation. A junction table would only be necessary for group chats with 3+ participants.</blockquote>
</details>
<details>
<summary><strong>CHECK constraints as database-level safety nets</strong></summary>
<blockquote>Business invariants are enforced at two levels: the service layer validates before writing, and the database rejects anything that slips through. <code>CHECK</code> constraints cover: <code>capacity > 0</code>, <code>price >= 0</code>, <code>available >= 0</code>, <code>number_of_tickets > 0</code>, <code>total_cost >= 0</code>, <code>end_datetime > start_datetime</code>, and enum-like status columns restricted to their valid values via <code>IN (...)</code> checks.</blockquote>
</details>
<details>
<summary><strong>UNIQUE constraints preventing real-world conflicts</strong></summary>
<blockquote><code>venue(latitude, longitude)</code> prevents two venues at identical coordinates. <code>event(venue_id, start_datetime)</code> prevents double-booking a venue at the same time. <code>user(username)</code> and <code>user(email)</code> enforce identity uniqueness. <code>category(name)</code> prevents duplicate category entries. These constraints act as the final line of defense — the application validates first, but the database guarantees correctness.</blockquote>
</details>
<details>
<summary><strong>CASCADE deletes scoped to owned children</strong></summary>
<blockquote><code>ON DELETE CASCADE</code> is applied only to tables that are fully owned by their parent and have no meaning without it: <code>tickettype</code>, <code>media</code>, <code>visit</code>, <code>eventcategory</code>, <code>conversation_participant</code>, and <code>message</code> cascade from their parent. <code>booking</code> does not cascade — bookings are retained for audit and history even if the referenced ticket type changes.</blockquote>
</details>
<details>
<summary><strong>Self-referencing category hierarchy</strong></summary>
<blockquote><code>category.parent_id</code> references <code>category.id</code>, enabling a tree of subcategories without a separate table. Nullable — top-level categories have <code>parent_id = NULL</code>. This supports arbitrary nesting depth while keeping the schema flat.</blockquote>
</details>
<details>
<summary><strong>Normalization summary</strong></summary>
<blockquote><strong>1NF:</strong> All 13 tables pass. No multi-valued columns anywhere. The many-to-many relationship between events and categories is handled through the <code>eventcategory</code> junction table — no comma-separated lists.<br><br><strong>2NF:</strong> All tables pass. Single-UUID-PK tables satisfy 2NF automatically. Junction tables (<code>eventcategory</code>, <code>conversation_participant</code>, <code>recommendation</code>) have no non-key columns that depend on only part of their composite key.<br><br><strong>3NF:</strong> All tables pass. The original transitive dependency in <code>conversation</code> (attendee/organizer derivable via booking) was identified and fixed. <code>total_cost</code> on <code>booking</code> is a documented intentional denormalization for price snapshot purposes.</blockquote>
</details>


# Backend Layers
 
### Entities
 
Entities are pure data structs mapped 1:1 to database tables. They carry no JSON tags — serialization belongs to the controller layer. What they do carry is domain validation logic as named methods:
 
| Method | Entity | Guards |
|--------|--------|--------|
| `ApproveCreate()` | Event | Required fields, future start date, end > start |
| `ApproveCancellation()` | Event | Must be in PUBLISHED status |
| `ApproveDeletion()` | Event | Must not be CANCELLED |
| `IsBookingAvailable()` | Event | Must be PUBLISHED |
| `AllowsTicketModification()` | Event | Must be DRAFT or PUBLISHED |
| `HasCapacityFor(current, add)` | Event | Sum of ticket quantities vs event capacity |
| `HasAvailability(requested)` | TicketType | Available seats vs requested quantity |
| `ApproveCreate()` | Media | 5MB cap, extension whitelist |
| `CanEditContent()` | Message | Only sender, not deleted |
| `CanTransitionTo(status)` | Message | Valid status transitions, no reversal |
| `ValidateContent()` | Message | No empty or whitespace-only messages |
 
The service layer orchestrates; the entity decides whether an operation is valid. If a business rule can be expressed as a boolean on the entity's own fields, it lives on the entity — not in the service.
 
### Services
 
Each service is describable in a single sentence with no "and". If "and" appears, the service gets split.
 
<details>
 
<summary><strong>All services and responsibilities</strong></summary>

<br>
 
| Service | One-sentence responsibility |
|---------|---------------------------|
| `AuthService` | Handles registration, login, password hashing, and JWT token generation. |
| `UserService` | Manages user retrieval, filtering, and admin approve/reject workflows. |
| `EventService` | Owns event CRUD, status transitions (publish/cancel), and emits domain events via the EventBus. |
| `CancelEventService` | Subscribes to `EventCancelled` events and notifies all affected attendees via messaging. |
| `TicketTypeService` | Owns ticket type creation and updates with capacity validation against the parent event. |
| `BookingService` | Creates bookings with atomic availability decrement inside a single transaction. |
| `MessageService` | Manages conversations, messages, read state, and enforces messaging domain rules. |
| `MediaService` | Validates and tracks photo uploads/deletes with ownership checks (file I/O stays in the handler). |
| `VenueService` | Lists and filters venues by name and capacity. |
| `VisitService` | Records event view interactions with in-memory 10-second cooldown deduplication. |
| `ExportService` | Assembles full event data (with tickets and bookings) for admin XML/JSON export. |
 
</details>

## Service layer decisions
 
<details>
<summary><strong>Services never call other services</strong></summary>
<br>
<blockquote>Early versions had services calling other services through interfaces like <code>EventsProvider</code>. This was removed — it risks exploiting gaps in the called service's authorization logic and creates hidden coupling. Cross-aggregate reads go directly to repos through narrow, consumer-defined interfaces (<code>EventCapacityProvider</code>, <code>VenueOwnershipChecker</code>). Each service depends only on the repos it needs.</blockquote>
</details>
<details>
<summary><strong>TicketTypeService split from BookingService</strong></summary>
<br>
<blockquote>BookingService originally owned ticket types, bookings, and capacity validation — three responsibilities. It was split into <code>BookingService</code> (booking creation only) and <code>TicketTypeService</code> (ticket lifecycle and capacity invariant). The capacity check lives in <code>TicketTypeService</code> because ticket quantity is the invariant being protected, even though bookings consume the availability.</blockquote>
</details>
<details>
<summary><strong>File I/O kept out of services</strong></summary>
<br>
<blockquote><code>MediaService</code> handles ownership verification, conflict checks, and domain validation (<code>ApproveCreate()</code>). The actual file write to disk and <code>os.Remove</code> happen in the handler. File bytes never cross the service boundary — the service works only with metadata. This keeps services testable without a filesystem.</blockquote>
</details>
<details>
<summary><strong>Ownership checks at the service layer, role checks at middleware</strong></summary>
<br>
<blockquote>Middleware handles role-based access (<code>AdminOnly()</code>, <code>AuthMiddleware()</code>). The service layer handles ownership — "is this user the organizer of this event?" — because only the service has access to the repos needed to answer that question. Controllers never make authorization decisions.</blockquote>
</details>
<details>
<summary><strong>Handlers orchestrate cross-domain concerns</strong></summary>
<br>
<blockquote>When a booking requires verifying that an event is <code>PUBLISHED</code>, the handler fetches the event and checks status before delegating to <code>BookingService</code>. The handler is the orchestrator — it coordinates between domains without containing business logic itself. This avoids making <code>BookingService</code> depend on <code>EventService</code>.</blockquote>
</details>
<details>
<summary><strong>ExportService bypasses service-layer auth</strong></summary>
<br>
<blockquote>The admin export endpoint needs data across all organizers. <code>ExportService</code> calls repos directly rather than going through domain services, which would reject the request due to ownership checks. This is intentional — the admin middleware has already verified the caller is an admin before the handler is reached.</blockquote>
</details>

### Controllers
 
Controllers are thin — they parse requests, call the appropriate service, and map the result to a response DTO. Three rules:
 
1. **All DTOs and marshalling tags live here** — services return bare Go structs, controllers add `json` and `xml` tags via separate response types
2. **Mapper functions** (`ToEventResponse`, `ToBookingListResponse`) convert service output to the transport shape
3. **No business logic** — if a controller needs an `if` that isn't about request parsing or response format, the logic belongs in a service or entity method

## ML Recommendation Pipeline
 
The recommendation engine uses **Biased Matrix Factorization** to predict user-event affinity scores from two interaction signals: visits (weak interest) and bookings (strong interest).
 
The pipeline is a standalone Python process — no runtime coupling to the Go backend. It reads interaction data from PostgreSQL, trains the model, and writes scored recommendations back to the `recommendation` table. The Go backend reads from that table at query time.
 
### DataLoader interface
 
The model accepts any data source through a `DataLoader` interface with three implementations:
 
| Implementation | Purpose |
|----------------|---------|
| `MockDataLoader` | Synthetic data for unit testing and rapid iteration |
| `FileDataLoader` | JSON datasets for reproducible offline evaluation |
| `DatabaseDataLoader` | Live PostgreSQL data for production training |
 
This made it possible to develop and tune the model entirely offline, then swap to the real database for final training without changing the model code.
 
### Evaluation
 
Training includes a train/test split with the following metrics tracked per run: RMSE, MAE, Precision@K, Recall@K, and NDCG@K.
 


## Frontend Architecture
 
### Decisions
 
<details>
<summary><strong>StaticDataProvider — fetch once, share globally</strong></summary>
<br>
<blockquote>Venues and categories are reference data that rarely changes. <code>StaticDataProvider</code> fetches both once at app mount and exposes them via <code>useStaticData()</code>. Every form that needs a venue dropdown or category picker consumes the same cached data — no per-component <code>useEffect</code> fetches, no duplicate requests, no stale copies.</blockquote>
</details>
<details>
<summary><strong>MessagesContext — centralized polling and unread state</strong></summary>
<br>
<blockquote>Before this context existed, the sidebar, the messages page, and the layout all independently polled for conversations and tracked unread state — three components, three fetch loops, three sources of truth. <code>MessagesProvider</code> replaced all of them with a single 15-second polling loop that exposes <code>conversations</code>, <code>hasUnread</code>, <code>refresh()</code>, and <code>markAsRead()</code>. Polling is skipped for guests and admins. Read-state updates propagate instantly across every component that consumes the context.</blockquote>
</details>
<details>
<summary><strong>GuestGate — inline prompts instead of hard redirects</strong></summary>
<br>
<blockquote>Unauthenticated users are never forcefully redirected to <code>/login</code>. Instead, a reusable <code>GuestGate</code> component renders in-place with a contextual prompt ("Log in to book tickets", "Register to message the organizer") and direct login/register CTAs. This lets guests browse the full platform naturally — restricted actions guide them toward authentication rather than blocking them. Three variants exist: full-page, compact (inside dialogs), and embedded (inline within cards).</blockquote>
</details>
<details>
<summary><strong>Client-side hydration for bare API responses</strong></summary>
<br>
<blockquote>The backend returns bare responses with foreign key IDs. The frontend API layer hydrates the full shape via parallel fan-out: <code>getUserBookings</code> fetches bookings, then batch-fetches ticket types and events in parallel, then hydrates each event with its venue. The result matches the component contract without the backend needing to compose multi-table responses. Trade-off accepted: <code>1 + 2N</code> round trips per list — acceptable at this scale, batch endpoints can follow if needed.</blockquote>
</details>
<details>
<summary><strong>CSS architecture — tokens, utilities, component-scoped files</strong></summary>
<br>
<blockquote>No inline styles anywhere. CSS is organized into four layers: <code>tokens.css</code> (color, spacing, typography variables), <code>base.css</code> (resets and element defaults), <code>utilities.css</code> (reusable single-purpose classes), and component-scoped files (<code>Events.css</code>, <code>Tickets.css</code>, <code>Bookings.css</code>, <code>Messages.css</code>). Class names follow BEM-style conventions with component prefixes. Page-level layout lives in <code>layout.css</code>. No bare element selectors in shared files — every rule is scoped to a class.</blockquote>
</details>
 
 
