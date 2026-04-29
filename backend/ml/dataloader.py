
from __future__ import annotations

import json
from abc import ABC, abstractmethod
from pathlib import Path
from typing import List, Tuple


from abc import ABC, abstractmethod
 
 
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
