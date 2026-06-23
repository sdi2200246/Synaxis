from __future__ import annotations
import json
import csv
from abc import ABC, abstractmethod
from pathlib import Path
from typing import List, Tuple
import psycopg2
from psycopg2.extras import execute_values
from uuid import UUID
import os
from dotenv import load_dotenv

 
class DataLoader(ABC):
    """
    Interface for loading visit and booking data.
    Swap MockDataLoader for PostgresDataLoader later without touching the model.
    """
 
    @abstractmethod
    def load_visits(self) -> list[tuple[str, str, float]]:
        pass
 
    @abstractmethod
    def load_bookings(self) -> list[tuple[str, str, float]]:
        pass
 
    @abstractmethod
    def load_users(self) -> list[str]:
        pass
 
    @abstractmethod
    def load_events(self) -> list[str]:
        pass
 
    def load_ratings(self) -> list[tuple[str, str, float]]:

        combined: dict[tuple[str, str], float] = {}
 
        for u, e, r in self.load_visits():
            combined[(u, e)] = combined.get((u, e), 0.0) + r
 
        for u, e, r in self.load_bookings():
            combined[(u, e)] = combined.get((u, e), 0.0) + r
 
        return [(u, e, min(r, 5.0)) for (u, e), r in combined.items()]
 
 
class MockDataLoader(DataLoader):
    """
    Synthetic data with deliberate patterns so we can verify the model learns:
      - Users 0-2: prefer events 0-2 (group A)
      - Users 3-5: prefer events 3-5 (group B)
      - User 6: mixed tastes (control)
 
    Bookings reinforce the strongest preferences within each group.
    After merging:
      - booked + visited event → 5.0 (capped)
      - visited only           → 1.0
      - weak cross-group visit → 1.0
    """
 
    def __init__(self):
        self._users = [f"user_{i}" for i in range(7)]
        self._events = [f"event_{i}" for i in range(7)]
 
        self._visits = [
            # group A
            ("user_0", "event_0", 1.0), ("user_0", "event_1", 1.0), ("user_0", "event_2", 1.0),
            ("user_1", "event_0", 1.0), ("user_1", "event_1", 1.0), ("user_1", "event_2", 1.0),
            ("user_2", "event_0", 1.0), ("user_2", "event_1", 1.0), ("user_2", "event_2", 1.0),
            # weak cross-group visits
            ("user_0", "event_3", 1.0), ("user_1", "event_4", 1.0), ("user_2", "event_5", 1.0),
 
            # group B
            ("user_3", "event_3", 1.0), ("user_3", "event_4", 1.0), ("user_3", "event_5", 1.0),
            ("user_4", "event_3", 1.0), ("user_4", "event_4", 1.0), ("user_4", "event_5", 1.0),
            ("user_5", "event_3", 1.0), ("user_5", "event_4", 1.0), ("user_5", "event_5", 1.0),
            # weak cross-group visits
            ("user_3", "event_0", 1.0), ("user_4", "event_1", 1.0), ("user_5", "event_2", 1.0),
 
            # user 6: mixed
            ("user_6", "event_1", 1.0), ("user_6", "event_3", 1.0), ("user_6", "event_6", 1.0),
        ]
 
        # bookings add 4.0 on top — visited+booked pairs will hit the 5.0 cap
        self._bookings = [
            # group A: each user booked their favourite
            ("user_0", "event_0", 4.0),
            ("user_1", "event_1", 4.0),
            ("user_2", "event_2", 4.0),
 
            # group B
            ("user_3", "event_3", 4.0),
            ("user_4", "event_4", 4.0),
            ("user_5", "event_5", 4.0),
 
            # user 6: confirmed both interests
            ("user_6", "event_1", 4.0),
            ("user_6", "event_3", 4.0),
        ]
 
    def load_visits(self) -> list[tuple[str, str, float]]:
        return self._visits
 
    def load_bookings(self) -> list[tuple[str, str, float]]:
        return self._bookings
 
    def load_users(self) -> list[str]:
        return self._users
 
    def load_events(self) -> list[str]:
        return self._events
 

class FileDataLoader(DataLoader):
    """
    Load synthetic recommender data from a JSON file with keys:
    users, events, visits, bookings.
    """

    def __init__(self, dataset_path: str | Path):
        self.dataset_path = Path(dataset_path)
        if not self.dataset_path.exists():
            raise FileNotFoundError(f"Dataset not found: {self.dataset_path}")

        with self.dataset_path.open("r", encoding="utf-8") as f:
            self._data = json.load(f)

        self._users = list(self._data.get("users", []))
        self._events = list(self._data.get("events", []))
        self._visits = [tuple(item) for item in self._data.get("visits", [])]
        self._bookings = [tuple(item) for item in self._data.get("bookings", [])]

    def load_visits(self) -> List[Tuple[str, str, float]]:
        return self._visits

    def load_bookings(self) -> List[Tuple[str, str, float]]:
        return self._bookings

    def load_users(self) -> List[str]:
        return self._users

    def load_events(self) -> List[str]:
        return self._events

class Database(DataLoader):
    def __init__(self):
        try:
            load_dotenv("../.env")
            self.conn = psycopg2.connect(os.getenv("DATABASE_URL"))
            self.cur = self.conn.cursor()
            self.bookings = None

        except Exception as e:
            print("DB connection failed:", e)


    def test(self):
        cur = self.conn.cursor()
        cur.execute("SELECT version();")
        print(cur.fetchone())

    def load_visits(self) -> List[Tuple[str, str, float]]:
        try:

            self.cur.execute("""
                SELECT DISTINCT user_id , event_id
                FROM visit
            """)

            rows = self.cur.fetchall()
            return [(str(r[0]), str(r[1]) ,  1.0) for r in rows]

        except Exception as e:
            print("DB visits fetch failed:", e)
            return []

    def load_bookings(self) -> List[Tuple[str, str, float]]:
        try:
            self.cur.execute("""
                SELECT DISTINCT b.user_id, e.id AS event_id
                FROM booking b
                JOIN tickettype tt ON tt.id = b.ticket_type_id
                JOIN event e ON e.id = tt.event_id
            """)
            rows = self.cur.fetchall()
            self.bookings =  [(str(r[0]), str(r[1]) ,  6.0) for r in rows]
            return self.bookings
        
        except Exception as e:
            print("DB bookings fetch failed:", e)
            return []
    
    def load_users(self) -> List[str]:
        try:
            self.cur.execute("""
                SELECT id
                FROM "user"
                WHERE status = 'approved'
                AND role = 'user';
            """)

            users = self.cur.fetchall()
            return [str(u[0]) for u in users]

        except Exception as e:
            print("DB users fetch failed:", e)
            return []

    def load_events(self) -> List[str]:
        try:
            self.cur.execute("""
                SELECT id
                FROM event 
                WHERE status != 'DRAFT';
            """)
            events = self.cur.fetchall()
            return [str(e[0]) for e in events]
        except Exception as e:
            print("DB events fetch failed:", e)
            return []


    def save_recommendations(self, all_recommendations):
        cleaned_recs = [
            (str(rec[0]), str(rec[1]), float(rec[2])) 
            for rec in all_recommendations
        ]
        
        try:
            print("🧹 Clearing old recommendations...")
            self.cur.execute("TRUNCATE TABLE recommendation;")
            
            print(f"📥 Saving {len(cleaned_recs)} fresh recommendations to the database...")
            query = """
                INSERT INTO recommendation (user_id, event_id, score)
                VALUES %s;
            """
            execute_values(self.cur, query, cleaned_recs)
            self.conn.commit()
            
        except Exception as e:
            self.conn.rollback()
            raise e

    
class CsvDataLoader(DataLoader):
    def __init__(self, data_dir: str | Path):
        self.data_dir = Path(data_dir)

    # status string -> rating
    ATTEND_SCALE = {"yes": 5.0, "maybe": 3.0, "invited": 2.0, "no": 0.5}

    def load_visits(self) -> list[tuple[str, str, float]]:
        # all positive/negative interest signals EXCEPT yes (that's a booking)
        out = []
        with open(self.data_dir / "event_interest.csv") as f:
            for row in csv.DictReader(f):
                if row["interested"] == "1":
                    out.append((row["user"], row["event"], 3.0))
                elif row["interested"] == "0":
                    out.append((row["user"], row["event"], 1.0))  # not_interested
        with open(self.data_dir / "event_attendees.csv") as f:
            for row in csv.DictReader(f):
                if not row["user_id"].strip():
                    continue
                s = row["status"]
                if s in ("maybe", "invited", "no"):
                    out.append((row["user_id"], row["event"], self.ATTEND_SCALE[s]))
        return out

    def load_bookings(self) -> list[tuple[str, str, float]]:
        out = []
        with open(self.data_dir / "event_attendees.csv") as f:
            for row in csv.DictReader(f):
                if row["status"] == "yes" and row["user_id"].strip():
                    out.append((row["user_id"], row["event"], 5.0))
        return out

    def load_users(self) -> list[str]:
        return list({u for u, _, _ in self.load_visits()} |
                    {u for u, _, _ in self.load_bookings()})

    def load_events(self) -> list[str]:
        return list({e for _, e, _ in self.load_visits()} |
                    {e for _, e, _ in self.load_bookings()})

if __name__ == "__main__":
    db = CsvDataLoader("rel_event_csvs")
    results = len(db.load_users())
    print(results)
