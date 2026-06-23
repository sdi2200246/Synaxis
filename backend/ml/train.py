import argparse
import random
from collections import defaultdict, Counter

from dataloader import DataLoader, Database, CsvDataLoader
from model import BiasedMF
from eval import evaluate_model


def train_test_split_userwise(ratings, test_ratio=0.2):
    by_user = defaultdict(list)
    for row in ratings:
        by_user[row[0]].append(row)

    train, test = [], []
    for u, rows in by_user.items():
        random.shuffle(rows)
        cut = max(1, int(len(rows) * test_ratio))
        test.extend(rows[:cut])
        train.extend(rows[cut:])
    return train, test


def run_local_test():
    """
    Pipeline 1: Local validation using CSV files.
    Splits data into train/test sets and evaluates model accuracy.
    """
    print("🚀 Starting Local Test Pipeline (CSV)...")
    random.seed(42)

    db = CsvDataLoader("rel_event_csvs")
    ratings = db.load_ratings()

    print("ratings before filter:", len(ratings))
    uc = Counter(u for u, _, _ in ratings)
    ec = Counter(e for _, e, _ in ratings)
    ratings = [(u, e, r) for u, e, r in ratings if uc[u] >= 7 and ec[e] >= 7]
    print("ratings after filter: ", len(ratings))

    train_ratings, test_ratings = train_test_split_userwise(ratings, test_ratio=0.2)

    model = BiasedMF(k=15, alpha=0.01, lam=0.02, epochs=700, seed=42, shuffle=True)
    model.fit(train_ratings)

    print("\ntest relevant count:", sum(1 for _, _, r in test_ratings if r >= 4.0))

    results = evaluate_model(model, train_ratings, test_ratings, k=20, relevance_threshold=4.0)
    print("── metrics ─────────────────")
    for metric, value in results.items():
        print(f"  {metric}: {value:.4f}")


def run_production():
    """
    Pipeline 2: Production execution using Postgres.
    Trains on 100% of available data, generates recommendations for ALL
    system users (including cold-starts), and writes results back to the database.
    """
    print("🔥 Starting Production Pipeline (Postgres DB)...")
    random.seed(42)

    db = Database()
    ratings = db.load_ratings()
    
    if not ratings:
        print("❌ Error: No ratings loaded from the database. Exiting.")
        return

    print(f"Total production ratings pulled: {len(ratings)}")
    
    uc = Counter(u for u, _, _ in ratings)
    ec = Counter(e for _, e, _ in ratings)
    filtered_ratings = [(u, e, r) for u, e, r in ratings]
    print(f"Ratings remaining after density filters: {len(filtered_ratings)}")

    model = BiasedMF(k=15, alpha=0.01, lam=0.02, epochs=700, seed=42, shuffle=True)
    model.fit(filtered_ratings)

    all_db_users = db.load_users()
    bookings = db.load_bookings()

    user_bookings_map = defaultdict(set)
    for user_id, event_id, _ in bookings:
        user_bookings_map[user_id].add(event_id)


    print(f"Generating recommendations for {len(all_db_users)} users...")
    all_recommendations = []
    
    for user_id in all_db_users:
        previous_bookings = user_bookings_map[user_id]
        

        user_top_n = model.top_n(user_id, n=20, exclude=previous_bookings)
        
        for event_id, score in user_top_n:
            all_recommendations.append((user_id, event_id, score))

    # 6. Push calculations straight into PostgreSQL
    print(f"Saving {len(all_recommendations)} recommendations to the database...")
    db.save_recommendations(all_recommendations)
    print("✅ Production pipeline complete. Recommendations synced successfully.")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Biased Matrix Factorization Pipeline Execution Router"
    )
    parser.add_argument(
        "mode",
        choices=["test", "prod"],
        help="Run 'test' for local evaluation via CSVs, or 'prod' to execute against live Postgres storage."
    )

    args = parser.parse_args()

    if args.mode == "test":
        run_local_test()
    elif args.mode == "prod":
        run_production()