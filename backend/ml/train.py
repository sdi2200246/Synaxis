from dataloader import DataLoader , Database
from model import BiasedMF
from eval import evaluate_model
import matplotlib.pyplot as plt
from collections import defaultdict
import random
 
 
def train(loader : DataLoader ,k =10, alpha=0.00001, lam=0.5, epochs=200):
    
    ratings = loader.load_ratings()
    
    print("── merged ratings (visit + booking) ─────────────────")
    for u, e, r in sorted(ratings)[:20]:
        print(f"  {u}  {e}  {r:.1f}")
 
    model = BiasedMF(k=k, alpha=alpha, lam=lam, epochs=epochs)
    history = model.fit(ratings)
    print(f"\nfinal loss: {history[-1]:.6f}")

    plt.figure(figsize=(8,5))
    plt.plot(history, linewidth=2)
    plt.title("Training Loss")
    plt.xlabel("Epoch")
    plt.ylabel("Loss")
    plt.grid(True, alpha=0.3)
    plt.show()

    return model, ratings
    
def train_test_split_userwise(ratings, test_ratio=0.2):
    by_user = defaultdict(list)

    for row in ratings:
        by_user[row[0]].append(row)

    train = []
    test = []

    for u, rows in by_user.items():
        random.shuffle(rows)

        cut = max(1, int(len(rows) * test_ratio))
        test.extend(rows[:cut])
        train.extend(rows[cut:])

    return train, test


 
if __name__ == "__main__":
    db = Database()
    users = db.load_users()
    bookings = db.load_bookings()
    ratings = db.load_ratings()
    clicks = db.load_visits()

    train_ratings, test_ratings = train_test_split_userwise(ratings)

    model = BiasedMF(
        k=15,
        alpha=0.07,
        lam=0.02,
        epochs=200,
        seed=42,
        shuffle=True
    )
    model.fit(train_ratings)
    all_recs = model.recomendations(20 , bookings)

    first_user = users[1]
    user_id = first_user

    user_recs = [r for r in all_recs if r[0] == user_id]
    user_booked_events = {e for u, e, r in bookings if u == user_id}
    user_visited_events = {e for u , e , r in clicks if u == user_id}

    print("User:", user_id)
    print("Bookings made by this user:", user_booked_events)
    print("clicks made by this user " , user_visited_events)
    print("Recommendations:")
    for _, event_id, score in user_recs:
        mark = "BOOKED" if event_id in user_booked_events else ""
        print(event_id, score, mark)