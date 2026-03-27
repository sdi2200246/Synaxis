CREATE TABLE "user" (
  "id" uuid PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "password_hash" varchar NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "phone" varchar,
  "address" varchar,
  "city" varchar,
  "country" varchar,
  "tax_id" varchar,
  "role" varchar NOT NULL DEFAULT 'user',
  "status" varchar NOT NULL DEFAULT 'pending',
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "venue" (
  "id" uuid PRIMARY KEY,
  "name" varchar NOT NULL,
  "address" varchar NOT NULL,
  "city" varchar NOT NULL,
  "country" varchar NOT NULL,
  "latitude" decimal,
  "longitude" decimal
);

CREATE TABLE "event" (
  "id" uuid PRIMARY KEY,
  "organizer_id" uuid NOT NULL,
  "venue_id" uuid NOT NULL,
  "title" varchar NOT NULL,
  "event_type" varchar NOT NULL,
  "status" varchar NOT NULL DEFAULT 'DRAFT',
  "description" text NOT NULL,
  "capacity" integer NOT NULL,
  "start_datetime" timestamp NOT NULL,
  "end_datetime" timestamp NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "category" (
  "id" uuid PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "parent_id" uuid
);

CREATE TABLE "tickettype" (
  "id" uuid PRIMARY KEY,
  "event_id" uuid NOT NULL,
  "name" varchar NOT NULL,
  "price" decimal NOT NULL,
  "quantity" integer NOT NULL,
  "available" integer NOT NULL
);

CREATE TABLE "booking" (
  "id" uuid PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "ticket_type_id" uuid NOT NULL,
  "number_of_tickets" integer NOT NULL,
  "total_cost" decimal NOT NULL,
  "status" varchar NOT NULL DEFAULT 'ACTIVE',
  "booked_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "conversation" (
  "id" uuid PRIMARY KEY,
  "booking_id" uuid NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "message" (
  "id" uuid PRIMARY KEY,
  "conversation_id" uuid NOT NULL,
  "sender_id" uuid NOT NULL,
  "content" text NOT NULL,
  "is_read" boolean NOT NULL DEFAULT false,
  "sent_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "visit" (
  "id" uuid PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "event_id" uuid NOT NULL,
  "visited_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "media" (
  "id" uuid PRIMARY KEY,
  "event_id" uuid NOT NULL,
  "filename" varchar NOT NULL,
  "uploaded_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "eventcategory" (
  "event_id" uuid NOT NULL,
  "category_id" uuid NOT NULL,
  PRIMARY KEY ("event_id", "category_id")
);

-- Foreign keys
ALTER TABLE "event" ADD CONSTRAINT event_organizer_fkey FOREIGN KEY ("organizer_id") REFERENCES "user" ("id");
ALTER TABLE "event" ADD CONSTRAINT event_venue_fkey FOREIGN KEY ("venue_id") REFERENCES "venue" ("id");
ALTER TABLE "tickettype" ADD CONSTRAINT tickettype_event_fkey FOREIGN KEY ("event_id") REFERENCES "event" ("id") ON DELETE CASCADE;
ALTER TABLE "booking" ADD CONSTRAINT booking_user_fkey FOREIGN KEY ("user_id") REFERENCES "user" ("id");
ALTER TABLE "booking" ADD CONSTRAINT booking_tickettype_fkey FOREIGN KEY ("ticket_type_id") REFERENCES "tickettype" ("id");
ALTER TABLE "conversation" ADD CONSTRAINT conversation_booking_fkey FOREIGN KEY ("booking_id") REFERENCES "booking" ("id");
ALTER TABLE "message" ADD CONSTRAINT message_conversation_fkey FOREIGN KEY ("conversation_id") REFERENCES "conversation" ("id") ON DELETE CASCADE;
ALTER TABLE "message" ADD CONSTRAINT message_sender_fkey FOREIGN KEY ("sender_id") REFERENCES "user" ("id");
ALTER TABLE "category" ADD CONSTRAINT category_parent_fkey FOREIGN KEY ("parent_id") REFERENCES "category" ("id");
ALTER TABLE "visit" ADD CONSTRAINT visit_user_fkey FOREIGN KEY ("user_id") REFERENCES "user" ("id");
ALTER TABLE "visit" ADD CONSTRAINT visit_event_fkey FOREIGN KEY ("event_id") REFERENCES "event" ("id") ON DELETE CASCADE;
ALTER TABLE "media" ADD CONSTRAINT media_event_fkey FOREIGN KEY ("event_id") REFERENCES "event" ("id") ON DELETE CASCADE;
ALTER TABLE "eventcategory" ADD CONSTRAINT eventcategory_event_fkey FOREIGN KEY ("event_id") REFERENCES "event" ("id") ON DELETE CASCADE;
ALTER TABLE "eventcategory" ADD CONSTRAINT eventcategory_category_fkey FOREIGN KEY ("category_id") REFERENCES "category" ("id");

-- CHECK constraints
ALTER TABLE "user" ADD CONSTRAINT user_role_check CHECK (role IN ('admin','user'));
ALTER TABLE "user" ADD CONSTRAINT user_status_check CHECK (status IN ('pending','approved','rejected'));
ALTER TABLE "event" ADD CONSTRAINT event_status_check CHECK (status IN ('DRAFT','PUBLISHED','COMPLETED','CANCELLED'));
ALTER TABLE "event" ADD CONSTRAINT event_capacity_check CHECK (capacity > 0);
ALTER TABLE "event" ADD CONSTRAINT event_datetime_check CHECK (end_datetime > start_datetime);
ALTER TABLE "tickettype" ADD CONSTRAINT tickettype_price_check CHECK (price >= 0);
ALTER TABLE "tickettype" ADD CONSTRAINT tickettype_quantity_check CHECK (quantity > 0);
ALTER TABLE "tickettype" ADD CONSTRAINT tickettype_available_check CHECK (available >= 0);
ALTER TABLE "booking" ADD CONSTRAINT booking_tickets_check CHECK (number_of_tickets > 0);
ALTER TABLE "booking" ADD CONSTRAINT booking_cost_check CHECK (total_cost >= 0);
ALTER TABLE "booking" ADD CONSTRAINT booking_status_check CHECK (status IN ('ACTIVE','COMPLETED','CANCELLED'));

-- UNIQUE on venue coordinates
ALTER TABLE "venue" ADD CONSTRAINT venue_coordinates_unique UNIQUE (latitude, longitude);