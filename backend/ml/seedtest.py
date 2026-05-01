import argparse
import os
import random
import uuid
from collections import defaultdict
from datetime import datetime, timedelta
from pathlib import Path

import bcrypt
import psycopg2
from dotenv import load_dotenv
from psycopg2.extras import execute_values


BASE_DIR = Path(__file__).resolve().parent
load_dotenv(dotenv_path=BASE_DIR.parent / ".env")

REC_DB = (
    os.getenv("REC_DATABASE_URL")
    or os.getenv("DATABASE_URL")
    or "postgresql://postgres:postgres@localhost:5433/synaxis_rec?sslmode=disable"
)

DEFAULT_USERS = 3000
DEFAULT_EVENTS = 859
SEED = 42

random.seed(SEED)

CATEGORY_NAMES = ["Music", "Sports", "Theatre", "Conference", "Workshop", "Festival"]

CATEGORY_IDS = {
    name: str(uuid.uuid5(uuid.NAMESPACE_URL, f"synaxis:category:{name}"))
    for name in CATEGORY_NAMES
}

VENUES = [
    (str(uuid.uuid5(uuid.NAMESPACE_URL, "synaxis:venue:herod")), "Ωδείο Ηρώδου Αττικού", "Διονυσίου Αρεοπαγίτου", "Athens", "Greece", 37.9704, 23.7245, 5000),
    (str(uuid.uuid5(uuid.NAMESPACE_URL, "synaxis:venue:pallas")), "Θέατρο Παλλάς", "Βουκουρεστίου 5", "Athens", "Greece", 37.9792, 23.7351, 1200),
    (str(uuid.uuid5(uuid.NAMESPACE_URL, "synaxis:venue:oaka")), "ΟΑΚΑ", "Κηφισίας 37", "Maroussi", "Greece", 38.0368, 23.7875, 70000),
    (str(uuid.uuid5(uuid.NAMESPACE_URL, "synaxis:venue:megaro")), "Μέγαρο Μουσικής", "Βασιλίσσης Σοφίας", "Athens", "Greece", 37.9756, 23.7492, 1960),
    (str(uuid.uuid5(uuid.NAMESPACE_URL, "synaxis:venue:technopolis")), "Τεχνόπολη", "Πειραιώς 100", "Athens", "Greece", 37.9779, 23.7114, 3000),
    (str(uuid.uuid5(uuid.NAMESPACE_URL, "synaxis:venue:sef")), "ΣΕΦ", "Εθνάρχου Μακαρίου", "Piraeus", "Greece", 37.9400, 23.6650, 14000),
]

PROFILES = {
    "music_lover":   {"Music": 0.80, "Festival": 0.60, "Theatre": 0.10, "Sports": 0.05, "Conference": 0.05, "Workshop": 0.05},
    "sports_fan":    {"Sports": 0.85, "Festival": 0.30, "Music": 0.10, "Theatre": 0.05, "Conference": 0.05, "Workshop": 0.05},
    "tech_person":   {"Conference": 0.80, "Workshop": 0.75, "Music": 0.10, "Sports": 0.05, "Theatre": 0.05, "Festival": 0.10},
    "theatre_goer":  {"Theatre": 0.85, "Music": 0.40, "Festival": 0.20, "Sports": 0.05, "Conference": 0.05, "Workshop": 0.10},
    "festival_goer": {"Festival": 0.85, "Music": 0.60, "Sports": 0.30, "Theatre": 0.20, "Conference": 0.10, "Workshop": 0.10},
    "generalist":    {"Music": 0.30, "Sports": 0.30, "Theatre": 0.30, "Conference": 0.30, "Festival": 0.30, "Workshop": 0.30},
}
PROFILE_NAMES = list(PROFILES.keys())

EVENT_TEMPLATES = {
    "Music":      ["Rock Night", "Jazz Evening", "Classical Recital", "Electronic Fest", "Piano Concert"],
    "Sports":     ["Marathon Camp", "Basketball Final", "Football Match", "Tennis Open", "Swimming Gala"],
    "Theatre":    ["Hamlet", "Medea", "Modern Drama", "Comedy Night", "Opera Gala"],
    "Conference": ["Tech Summit", "AI Conference", "Startup Expo", "Dev Forum", "Data Summit"],
    "Workshop":   ["React Workshop", "Go Bootcamp", "ML Intro", "UX Design", "DevOps Day"],
    "Festival":   ["Street Food Fest", "Cultural Fest", "Film Festival", "Art Fest", "Wine Fest"],
}

EVENT_TYPE_BY_CATEGORY = {
    "Music": "Concert",
    "Sports": "Match",
    "Theatre": "Play",
    "Conference": "Conference",
    "Workshop": "Workshop",
    "Festival": "Festival",
}

CATEGORY_BASE_PRICE = {
    "Music": 25.0,
    "Sports": 20.0,
    "Theatre": 30.0,
    "Conference": 40.0,
    "Workshop": 15.0,
    "Festival": 20.0,
}

CATEGORY_HOURS = {
    "Music": [19, 20, 21, 22],
    "Sports": [18, 19, 20, 21],
    "Theatre": [19, 20, 21],
    "Conference": [9, 10, 11, 12, 13],
    "Workshop": [10, 12, 14, 16],
    "Festival": [12, 14, 16, 18, 20],
}

CATEGORY_DURATIONS = {
    "Music": 3,
    "Sports": 3,
    "Theatre": 3,
    "Conference": 8,
    "Workshop": 6,
    "Festival": 5,
}

CATEGORY_MONTHS = {
    "Music": (5, 12),
    "Sports": (5, 11),
    "Theatre": (5, 12),
    "Conference": (9, 11),
    "Workshop": (5, 12),
    "Festival": (6, 9),
}

PASSWORD_HASH = bcrypt.hashpw(b"password123", bcrypt.gensalt()).decode()

FIXED_USERS = [
    {
        "id": str(uuid.uuid5(uuid.NAMESPACE_URL, "synaxis:user:organizer")),
        "username": "organizer",
        "first_name": "Org",
        "last_name": "User",
        "email": "organizer@rec.test",
        "phone": "0000000001",
        "address": "Org St 1",
        "city": "Athens",
        "country": "Greece",
        "tax_id": "000000000",
        "role": "admin",
        "status": "approved",
        "profile": "generalist",
    },
    {
        "id": str(uuid.uuid5(uuid.NAMESPACE_URL, "synaxis:user:john")),
        "username": "john",
        "first_name": "John",
        "last_name": "Doe",
        "email": "john@rec.test",
        "phone": "0000000002",
        "address": "Ermou 10",
        "city": "Athens",
        "country": "Greece",
        "tax_id": "123456789",
        "role": "user",
        "status": "approved",
        "profile": "music_lover",
    },
    {
        "id": str(uuid.uuid5(uuid.NAMESPACE_URL, "synaxis:user:maria")),
        "username": "maria",
        "first_name": "Maria",
        "last_name": "Papadopoulou",
        "email": "maria@rec.test",
        "phone": "0000000003",
        "address": "Stadiou 15",
        "city": "Athens",
        "country": "Greece",
        "tax_id": "987654321",
        "role": "user",
        "status": "pending",
        "profile": "theatre_goer",
    },
    {
        "id": str(uuid.uuid5(uuid.NAMESPACE_URL, "synaxis:user:jake")),
        "username": "jake",
        "first_name": "Jake",
        "last_name": "Doe",
        "email": "jake@rec.test",
        "phone": "0000000004",
        "address": "Ermou 9",
        "city": "Athens",
        "country": "Greece",
        "tax_id": "123456785",
        "role": "user",
        "status": "approved",
        "profile": "sports_fan",
    },
    {
        "id": str(uuid.uuid5(uuid.NAMESPACE_URL, "synaxis:user:alice")),
        "username": "alice",
        "first_name": "Alice",
        "last_name": "Music",
        "email": "alice@rec.test",
        "phone": "0000000005",
        "address": "Kifisias 1",
        "city": "Athens",
        "country": "Greece",
        "tax_id": "111111111",
        "role": "user",
        "status": "approved",
        "profile": "music_lover",
    },
    {
        "id": str(uuid.uuid5(uuid.NAMESPACE_URL, "synaxis:user:carol")),
        "username": "carol",
        "first_name": "Carol",
        "last_name": "Tech",
        "email": "carol@rec.test",
        "phone": "0000000006",
        "address": "Patision 5",
        "city": "Athens",
        "country": "Greece",
        "tax_id": "222222222",
        "role": "user",
        "status": "approved",
        "profile": "tech_person",
    },
]

TIER_RULES = {
    "power":  {"bookings": 20, "clicks": 80, "booking_split": (14, 4, 2), "click_split": (84, 24, 12)},
    "active": {"bookings": 10, "clicks": 60,  "booking_split": (7, 2, 1),   "click_split": (42, 12, 6)},
    "casual": {"bookings": 4,  "clicks": 25,  "booking_split": (3, 1, 0),   "click_split": (18, 5, 2)},
}


def bulk_insert(cur, sql: str, rows: list[tuple], page_size: int = 1000) -> None:
    if rows:
        execute_values(cur, sql, rows, page_size=page_size)


def reset_db(cur) -> None:
    cur.execute("""
        TRUNCATE TABLE
            "recommendation",
            "media",
            "message",
            conversation_participant,
            "conversation",
            "booking",
            "tickettype",
            "eventcategory",
            "visit",
            "event",
            "category",
            "venue",
            "user"
        RESTART IDENTITY CASCADE;
    """)


def build_users(n_users: int) -> list[dict]:
    users = FIXED_USERS[:]
    if n_users <= len(users):
        return users[:n_users]

    remaining = n_users - len(users)

    for i in range(remaining):
        username = f"user_{i:05d}"
        profile = PROFILE_NAMES[i % len(PROFILE_NAMES)]

        mod = i % 20
        if mod < 16:
            status = "approved"
        elif mod < 18:
            status = "pending"
        else:
            status = "rejected"

        users.append({
            "id": str(uuid.uuid5(uuid.NAMESPACE_URL, f"synaxis:user:{username}")),
            "username": username,
            "first_name": f"First{i:05d}",
            "last_name": f"Last{i:05d}",
            "email": f"{username}@rec.test",
            "phone": f"{3000000000 + i}",
            "address": "Synthetic Street 1",
            "city": "Athens",
            "country": "Greece",
            "tax_id": f"{100000000 + i % 900000000}",
            "role": "user",
            "status": status,
            "profile": profile,
        })

    return users


def build_events(n_events: int) -> list[dict]:
    events: list[dict] = []
    used_slots: set[tuple[str, str]] = set()

    counts = [
        n_events // len(CATEGORY_NAMES) + (1 if i < (n_events % len(CATEGORY_NAMES)) else 0)
        for i in range(len(CATEGORY_NAMES))
    ]

    global_idx = 1
    for cat_idx, (cat, count) in enumerate(zip(CATEGORY_NAMES, counts)):
        templates = EVENT_TEMPLATES[cat]
        month_min, month_max = CATEGORY_MONTHS[cat]
        hours = CATEGORY_HOURS[cat]

        for i in range(count):
            venue = VENUES[(global_idx - 1) % len(VENUES)]
            venue_id = venue[0]

            month = month_min + ((global_idx - 1) % (month_max - month_min + 1))
            day = 1 + ((global_idx - 1) % 28)
            hour = hours[i % len(hours)]

            start_dt = datetime(2026, month, day, hour, 0, 0)
            while (venue_id, start_dt.isoformat(sep=" ")) in used_slots:
                start_dt += timedelta(hours=1)

            used_slots.add((venue_id, start_dt.isoformat(sep=" ")))
            end_dt = start_dt + timedelta(hours=CATEGORY_DURATIONS[cat])

            status_mod = global_idx % 20
            if status_mod == 0:
                status = "DRAFT"
            elif status_mod == 1:
                status = "CANCELLED"
            else:
                status = "PUBLISHED"

            capacity = min(venue[7], 200 + ((global_idx - 1) % 12) * 150)

            events.append({
                "id": str(uuid.uuid5(uuid.NAMESPACE_URL, f"synaxis:event:{cat}:{global_idx}")),
                "organizer_id": str(uuid.uuid5(uuid.NAMESPACE_URL, "synaxis:user:organizer")),
                "venue_id": venue_id,
                "title": f"{templates[i % len(templates)]} #{global_idx}",
                "event_type": EVENT_TYPE_BY_CATEGORY[cat],
                "status": status,
                "description": f"Auto-generated {cat.lower()} event for recommender testing.",
                "capacity": capacity,
                "start_datetime": start_dt,
                "end_datetime": end_dt,
                "category": cat,
                "category_id": CATEGORY_IDS[cat],
            })
            global_idx += 1

    return events


def _flatten_events(events_by_cat: dict[str, list[dict]], cats: list[str]) -> list[dict]:
    flat: list[dict] = []
    seen = set()
    for cat in cats:
        for ev in events_by_cat.get(cat, []):
            if ev["id"] not in seen:
                flat.append(ev)
                seen.add(ev["id"])
    return flat


def _take_deterministic(pool: list[dict], n: int, key: str, avoid_ids: set[str] | None = None) -> list[dict]:
    if n <= 0 or not pool:
        return []

    avoid_ids = avoid_ids or set()
    filtered = [ev for ev in pool if ev["id"] not in avoid_ids]
    if not filtered:
        return []

    offset = uuid.uuid5(uuid.NAMESPACE_URL, key).int % len(filtered)
    ordered = filtered[offset:] + filtered[:offset]
    return ordered[:n]


def build_interactions(users, events):
    published = sorted(
        [e for e in events if e["status"] == "PUBLISHED"],
        key=lambda e: (e["start_datetime"], e["id"])
    )

    events_by_cat: dict[str, list[dict]] = defaultdict(list)
    for e in published:
        events_by_cat[e["category"]].append(e)

    visits = []
    bookings = []

    approved_users = [u for u in users if u["status"] == "approved"]

    for approved_idx, u in enumerate(approved_users):
        tier_mod = approved_idx % 10
        if tier_mod in (0, 1):
            tier = "power"
        elif tier_mod in (2, 3, 4, 5, 6, 7):
            tier = "active"
        else:
            tier = "casual"

        rule = TIER_RULES[tier]
        profile_items = sorted(PROFILES[u["profile"]].items(), key=lambda kv: (-kv[1], kv[0]))
        cats = [c for c, _ in profile_items]

        primary_cats = cats[:2]
        secondary_cats = cats[2:4]
        explore_cats = [c for c in CATEGORY_NAMES if c not in cats]

        primary_pool = _flatten_events(events_by_cat, primary_cats)
        secondary_pool = _flatten_events(events_by_cat, secondary_cats)
        explore_pool = _flatten_events(events_by_cat, explore_cats)

        booked_ids = set()
        seen_ids = set()

        b1, b2, b3 = rule["booking_split"]
        booking_groups = [
            _take_deterministic(primary_pool, b1, f"book:{u['id']}:1"),
            _take_deterministic(secondary_pool, b2, f"book:{u['id']}:2"),
            _take_deterministic(explore_pool, b3, f"book:{u['id']}:3"),
        ]

        for group in booking_groups:
            for ev in group:
                if ev["id"] in booked_ids:
                    continue

                booked_ids.add(ev["id"])
                seen_ids.add(ev["id"])

                ticket = ev["ticket_types"][0]
                qty = 1 + (uuid.uuid5(uuid.NAMESPACE_URL, f"synaxis:qty:{u['id']}:{ev['id']}").int % 3)
                total = qty * float(ticket["price"])

                booked_at = ev["start_datetime"] - timedelta(days=21 + (approved_idx % 14))

                bookings.append((
                    str(uuid.uuid5(uuid.NAMESPACE_URL, f"synaxis:booking:{u['id']}:{ev['id']}")),
                    u["id"],
                    ticket["id"],
                    qty,
                    total,
                    "ACTIVE",
                    booked_at,
                ))

                visits.append((
                    str(uuid.uuid5(uuid.NAMESPACE_URL, f"synaxis:visit:{u['id']}:{ev['id']}")),
                    u["id"],
                    ev["id"],
                    ev["start_datetime"] - timedelta(days=7 + (approved_idx % 10)),
                ))

        c1, c2, c3 = rule["click_split"]
        click_groups = [
            _take_deterministic(primary_pool, c1, f"click:{u['id']}:1", avoid_ids=booked_ids),
            _take_deterministic(secondary_pool, c2, f"click:{u['id']}:2", avoid_ids=booked_ids),
            _take_deterministic(explore_pool, c3, f"click:{u['id']}:3", avoid_ids=booked_ids),
        ]

        for group in click_groups:
            for ev in group:
                if ev["id"] in seen_ids:
                    continue

                seen_ids.add(ev["id"])

                visits.append((
                    str(uuid.uuid5(uuid.NAMESPACE_URL, f"synaxis:visit:{u['id']}:{ev['id']}")),
                    u["id"],
                    ev["id"],
                    ev["start_datetime"] - timedelta(days=3 + (approved_idx % 7)),
                ))

    return visits, bookings


def seed(n_users: int, n_events: int) -> None:
    conn = psycopg2.connect(REC_DB)
    cur = conn.cursor()

    try:
        print("resetting database...")
        reset_db(cur)

        print("seeding categories...")
        category_rows = [(cid, name, None) for name, cid in CATEGORY_IDS.items()]
        bulk_insert(cur, "INSERT INTO category (id, name, parent_id) VALUES %s", category_rows)

        print("seeding venues...")
        bulk_insert(
            cur,
            """
            INSERT INTO venue (id, name, address, city, country, latitude, longitude, capacity)
            VALUES %s
            """,
            VENUES,
        )

        print("seeding users...")
        users = build_users(n_users)
        user_rows = [
            (
                u["id"],
                u["username"],
                PASSWORD_HASH,
                u["first_name"],
                u["last_name"],
                u["email"],
                u["phone"],
                u["address"],
                u["city"],
                u["country"],
                u["tax_id"],
                u["role"],
                u["status"],
            )
            for u in users
        ]
        bulk_insert(
            cur,
            """
            INSERT INTO "user" (
                id, username, password_hash, first_name, last_name, email,
                phone, address, city, country, tax_id, role, status
            )
            VALUES %s
            """,
            user_rows,
        )

        print("building events...")
        events = build_events(n_events)

        event_rows = [
            (
                e["id"],
                e["organizer_id"],
                e["venue_id"],
                e["title"],
                e["event_type"],
                e["status"],
                e["description"],
                e["capacity"],
                e["start_datetime"],
                e["end_datetime"],
            )
            for e in events
        ]
        bulk_insert(
            cur,
            """
            INSERT INTO event (
                id, organizer_id, venue_id, title, event_type, status,
                description, capacity, start_datetime, end_datetime
            )
            VALUES %s
            """,
            event_rows,
        )

        eventcategory_rows = [(e["id"], e["category_id"]) for e in events]
        bulk_insert(cur, "INSERT INTO eventcategory (event_id, category_id) VALUES %s", eventcategory_rows)

        print("building ticket types...")
        ticket_rows = []
        ticket_types_by_event: dict[str, list[dict]] = {}

        for e in events:
            cap = int(e["capacity"])
            general_qty = max(cap * 50, 10000)
            vip_qty = max(cap * 15, 3000)

            general_id = str(uuid.uuid5(uuid.NAMESPACE_URL, f"synaxis:ticket:{e['id']}:general"))
            vip_id = str(uuid.uuid5(uuid.NAMESPACE_URL, f"synaxis:ticket:{e['id']}:vip"))

            base_price = CATEGORY_BASE_PRICE[e["category"]]
            vip_price = round(base_price * 2.5, 2)

            ticket_types_by_event[e["id"]] = [
                {"id": general_id, "price": base_price},
                {"id": vip_id, "price": vip_price},
            ]

            ticket_rows.append((general_id, e["id"], "General", base_price, general_qty, general_qty))
            ticket_rows.append((vip_id, e["id"], "VIP", vip_price, vip_qty, vip_qty))

        bulk_insert(
            cur,
            """
            INSERT INTO tickettype (id, event_id, name, price, quantity, available)
            VALUES %s
            """,
            ticket_rows,
        )

        for e in events:
            e["ticket_types"] = ticket_types_by_event[e["id"]]

        print("building interactions...")
        visits, bookings = build_interactions(users, events)

        bulk_insert(
            cur,
            "INSERT INTO visit (id, user_id, event_id, visited_at) VALUES %s",
            visits,
            page_size=5000,
        )

        bulk_insert(
            cur,
            """
            INSERT INTO booking (
                id, user_id, ticket_type_id, number_of_tickets, total_cost, status, booked_at
            )
            VALUES %s
            """,
            bookings,
            page_size=5000,
        )

        booked_by_ticket = defaultdict(int)
        for row in bookings:
            ticket_type_id = row[2]
            number_of_tickets = row[3]
            booked_by_ticket[ticket_type_id] += number_of_tickets

        for ticket_type_id, booked_count in booked_by_ticket.items():
            cur.execute(
                "UPDATE tickettype SET available = available - %s WHERE id = %s",
                (booked_count, ticket_type_id),
            )

        conn.commit()

        approved_users = sum(1 for u in users if u["status"] == "approved")
        published_events = sum(1 for e in events if e["status"] == "PUBLISHED")

        print(
            f"done: {len(users)} users ({approved_users} approved), "
            f"{len(events)} events ({published_events} published), "
            f"{len(visits)} visits, {len(bookings)} bookings"
        )

    except Exception as e:
        conn.rollback()
        print("seed failed:", e)
        raise
    finally:
        cur.close()
        conn.close()


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--users", type=int, default=DEFAULT_USERS)
    parser.add_argument("--events", type=int, default=DEFAULT_EVENTS)
    args = parser.parse_args()

    seed(args.users, args.events)